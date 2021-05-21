package orm

import (
	"bark-server/model"
	"bark-server/util"
	"encoding/json"
	"github.com/mritd/logger"
)

func SaveMessage(msg model.PushMessage, pushErr error) (bool, error) {
	var pushed = true
	var errorMsg = ""
	var extParams = ""
	if pushErr != nil {
		pushed = false
		errorMsg = pushErr.Error()
	}
	if msg.Source == "" {
		msg.Source = "NORMAL"
	}
	if msg.ExtParams != nil {
		b, err := json.Marshal(msg.ExtParams)
		if err != nil {
			logger.Error("convert extParams error", err)
		} else {
			extParams = string(b)
		}
	}
	message := model.Message{
		Source:      msg.Source,
		DeviceKey:   msg.DeviceKey,
		DeviceToken: msg.DeviceToken,
		Title:       msg.Title,
		Content:     msg.Body,
		Sound:       msg.Sound,
		ExtParams:   extParams,
		Pushed:      pushed,
		ErrorMsg:    errorMsg,
	}
	if err := db.Save(&message).Error; err != nil {
		logger.Error("save message error.%+v", err)
		return false, err
	}
	return true, nil
}

func GetMessageCount() int {
	var count int64
	if err := db.Model(&model.Message{}).Count(&count).Error; err != nil {
		logger.Error("count message error.", err)
		return 0
	}
	return util.Int64ToInt(count)
}

func GetMessageList(page model.Page) (model.Page, error) {
	if page.PageNo <= 0 {
		page.PageNo = 1
	}
	if page.PageSize <= 0 {
		page.PageSize = 10
	}
	var messages []model.Message
	offset := (page.PageNo - 1) * page.PageSize
	limit := page.PageSize
	if err := db.Limit(limit).Offset(offset).Find(&messages).Error; err != nil {
		logger.Error("query message list error.", err)
		return page, err
	}
	page.Data = messages
	page.Total = GetMessageCount()
	page.TotalPage = page.Total / page.PageSize
	if page.Total%page.PageSize > 0 {
		page.TotalPage += 1
	}
	return page, nil
}
