package orm

import "github.com/mritd/logger"

func SaveMessage(message *PushMessage) (bool, error) {
	res := db.Create(&message)
	if res.Error != nil {
		logger.Errorf("save message error. token: %v, error: %v", message.DeviceToken, res.Error)
		return false, res.Error
	}
	logger.Infof("save message success. token: %v, key: %v", message.DeviceToken, message.DeviceKey)
	return true, nil
}
