package main

import (
	"context"
	"fmt"
	"html"

	"anon_bottg/database"

	"github.com/mymmrac/telego"
)

func text(bot *telego.Bot, update telego.Update) {
	userID := update.Message.From.ID
	userText := update.Message.Text

	if update.Message.ReplyToMessage != nil { // ответ владельца ссылки
		messageID := update.Message.ReplyToMessage.MessageID

		senderID, err := database.GetSender(messageID)
		if err != nil {
			fmt.Println("Ошибка! ID получателя не найден!", err)
			return
		}

		replyText := fmt.Sprintf(
			"💬 <b>Пришёл ответ на твоё сообщение:</b>\n\n"+
				"<blockquote>%s</blockquote>",
			html.EscapeString(userText),
		)

		_, _ = bot.SendMessage(context.Background(), &telego.SendMessageParams{
			ChatID:    telego.ChatID{ID: senderID},
			Text:      replyText,
			ParseMode: telego.ModeHTML,
		})

		_, _ = bot.SendMessage(context.Background(), &telego.SendMessageParams{
			ChatID:    telego.ChatID{ID: userID},
			Text:      "✅ <b>Ответ отправлен!</b>",
			ParseMode: telego.ModeHTML,
		})
	} else { // гость, перешедший по ссылке
		targetID, err := database.GetTarget(userID)
		if err != nil {
			_, _ = bot.SendMessage(context.Background(), &telego.SendMessageParams{
				ChatID:    telego.ChatID{ID: userID},
				Text:      "⚠️ <b>Кому отправить?</b>\nСначала перейди по чьей-либо персональной ссылке или свайпни входящее сообщение для ответа.",
				ParseMode: telego.ModeHTML,
			})
			return
		}

		textForAdmin := fmt.Sprintf(
			"📩 <b>Новое анонимное сообщение!</b>\n\n"+
				"<blockquote>%s</blockquote>\n\n"+
				"<i>↩️ Чтобы ответить, свайпни это сообщение</i>",
			html.EscapeString(userText),
		)

		sentMsg, err := bot.SendMessage(context.Background(), &telego.SendMessageParams{
			ChatID:    telego.ChatID{ID: targetID},
			Text:      textForAdmin,
			ParseMode: telego.ModeHTML,
		})
		if err != nil {
			fmt.Println("Ошибка отправки сообщения(1)", err)
			return
		}

		err = database.SaveSenderID(sentMsg.MessageID, userID)
		if err != nil {
			fmt.Println("Ошибка записи в бд senderid!", err)
			return
		}

		_, err = bot.SendMessage(context.Background(), &telego.SendMessageParams{
			ChatID:    telego.ChatID{ID: userID},
			Text:      "📨 <b>Доставлено!</b> Ожидай ответа в этом чате.",
			ParseMode: telego.ModeHTML,
		})
		if err != nil {
			fmt.Println("Ошибка отправки ответа пользователю!", err)
			return
		}
	}
}
