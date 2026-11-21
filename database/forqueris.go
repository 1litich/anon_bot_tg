package database

func GetTarget(userID int64) (int64, error) {
	var targetID int64
	querySelectTarget := "SELECT target_id FROM targetid WHERE user_id = $1"

	err := db.QueryRow(querySelectTarget, userID).Scan(&targetID)

	return targetID, err
}

func GetSender(messageID int) (int64, error) {
	var senderID int64
	querySelectSender := "SELECT sender_id FROM senderid WHERE message_id = $1"

	err := db.QueryRow(querySelectSender, messageID).Scan(&senderID)

	return senderID, err
}

func SaveSenderID(messageID int, senderID int64) error {
	queryInsertSentMsg := "INSERT INTO senderid (message_id, sender_id) VALUES ($1, $2)"

	_, err := db.Exec(queryInsertSentMsg, messageID, senderID)

	return err
}

func SaveTargetID(targetID, userID int64) error {
	queryInsertTarget := "INSERT INTO targetid (target_id, user_id) VALUES ($1, $2) ON CONFLICT (user_id) DO UPDATE SET target_id = EXCLUDED.target_id"

	_, err := db.Exec(queryInsertTarget, targetID, userID)
	return err
}
