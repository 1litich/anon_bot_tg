package main

import (
	"anon_bottg/database"
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/mymmrac/telego"
)

func main() {
	database.Database()
	defer database.Close()

	botToken := os.Getenv("TOKEN")

	if botToken == "" {
		fmt.Println("Ошибка: переменная окружения TOKEN не задана")
		os.Exit(1)
	}

	go func() {
		port := os.Getenv("PORT")
		if port == "" {
			port = "8080"
		}
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, "Бот работает!")
		})
		_ = http.ListenAndServe(":"+port, nil)
	}()

	bot, err := telego.NewBot(botToken, telego.WithDefaultDebugLogger(), telego.WithHTTPClient(&http.Client{}))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	botUser, err := bot.GetMe(context.Background())
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	updates, err := bot.UpdatesViaLongPolling(context.Background(), &telego.GetUpdatesParams{})
	if err != nil {
		fmt.Println("Ошибка запуска Long Polling:", err)
		return
	}
	fmt.Println("БОТ ЖЕСТКО СЛУШАЕТ КАНАЛ!!!")

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if update.Message.Text != "" {
			userID := update.Message.From.ID
			userText := update.Message.Text

			fmt.Printf("Получено сообщение от ID %d: %s\n", userID, userText)

			if userText == "/start" {
				welcomeText := fmt.Sprintf(
					"🕵️‍♂️ <b>Бот анонимных сообщений</b>\n\n"+
						"Хочешь задать вопрос, признаться в чём-то или просто написать пару слов? Скинь свою ссылку друзьям или в профиль — и получай сообщения полностью анонимно.\n\n"+
						"📎 <b>Поддерживается:</b> текст, фото, видео и стикеры.\n\n"+
						"🔗 <b>Твоя личная ссылка</b> <i>(нажми, чтобы скопировать)</i>:\n"+
						"<code>https://t.me/%s?start=%d</code>",
					botUser.Username, userID,
				)

				shareText := fmt.Sprintf("\nНапиши мне анонимное сообщение по ссылке:\nhttps://t.me/%s?start=%d", botUser.Username, userID)

				keyboard := &telego.InlineKeyboardMarkup{
					InlineKeyboard: [][]telego.InlineKeyboardButton{
						{
							{
								Text:              "📤 Поделиться ссылкой",
								SwitchInlineQuery: &shareText,
							},
						},
					},
				}

				_, err := bot.SendMessage(context.Background(), &telego.SendMessageParams{
					ChatID:      telego.ChatID{ID: userID},
					Text:        welcomeText,
					ParseMode:   telego.ModeHTML,
					ReplyMarkup: keyboard,
				})
				if err != nil {
					fmt.Println("Ошибка старта:", err)
				}
				continue
			}

			if strings.HasPrefix(userText, "/start ") {
				grID := strings.TrimPrefix(userText, "/start ")
				targetID, err := strconv.ParseInt(grID, 10, 64)

				if err != nil {
					_, _ = bot.SendMessage(context.Background(), &telego.SendMessageParams{
						ChatID:    telego.ChatID{ID: userID},
						Text:      "⚠️ <b>Некорректная ссылка.</b> Проверь, что перешёл по полному адресу.",
						ParseMode: telego.ModeHTML,
					})
					continue
				}

				if targetID == userID {
					_, _ = bot.SendMessage(context.Background(), &telego.SendMessageParams{
						ChatID:    telego.ChatID{ID: userID},
						Text:      "🙃 <b>Самому себе написать нельзя!</b> Отправь эту ссылку друзьям или в свой канал.",
						ParseMode: telego.ModeHTML,
					})
					continue
				}

				err = database.SaveTargetID(targetID, userID)
				if err != nil {
					fmt.Println("Ошибка добавления в targetid:", err)
					continue
				}

				_, _ = bot.SendMessage(context.Background(), &telego.SendMessageParams{
					ChatID:    telego.ChatID{ID: userID},
					Text:      "✅ <b>Чат открыт!</b>\n\nОтправь сюда текст, фото, видео или стикер — я анонимно передам всё владельцу ссылки.",
					ParseMode: telego.ModeHTML,
				})
				continue
			}

			text(bot, update)
		} else {
			picture(update, bot)
			sticker(update, bot)
			video(bot, update)
		}
	}
}
