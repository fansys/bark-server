package orm

import (
	"bark-server/model"
	"github.com/lithammer/shortuuid/v3"
	"github.com/mritd/logger"
)

func SaveUser(user model.User) (model.User, error) {
	newUser := model.User{
		Username:  user.Username,
		Password:  user.Password,
		Nickname:  user.Nickname,
		OpenId:    user.OpenId,
		UniqueKey: shortuuid.New(),
	}
	if err := db.Save(&newUser).Error; err != nil {
		logger.Error("save user error.%+v", err)
		return model.User{}, err
	}
	return newUser, nil
}

func BindUser(bindType string, key string, uniqueKey string) (bool, error) {
	userBind := model.UserBind{
		UniqueKey: uniqueKey,
		BindType:  bindType,
		Key:       key,
	}
	if err := db.Save(&userBind).Error; err != nil {
		logger.Error("bind user error.%+v", err)
		return false, err
	}
	return true, nil
}
