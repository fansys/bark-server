package hms

import (
	"bark-server/model"
	"bark-server/push/provider/hms/push/config"
	"bark-server/push/provider/hms/push/constant"
	"bark-server/push/provider/hms/push/core"
	hms_model "bark-server/push/provider/hms/push/model"
	"context"
	"encoding/json"
	"fmt"
	"github.com/mritd/logger"
	"github.com/sirupsen/logrus"
	"os"
	"strings"
	"time"
)

const (
	// get token address
	hmsAuthUrl = "https://login.cloud.huawei.com/oauth2/v2/token"
	// send push msg address
	hmsPushUrl = "https://api.push.hicloud.com"
)

var cli *core.HttpPushClient

func init() {
	conf := &config.Config{
		AppId:     os.Getenv("HMS_APP_ID"),
		AppSecret: os.Getenv("HMS_APP_SECRET"),
		AuthUrl:   hmsAuthUrl,
		PushUrl:   hmsPushUrl,
	}
	client, err := core.NewHttpClient(conf)
	if err != nil {
		logger.Fatalf("failed to init hms client: %+v", err.Error())
	} else {
		cli = client
		logger.Info("init hms client success...")
	}
}

func Push(msg *model.PushMessage) error {
	msgRequest, err := buildHmsMsgRequest(msg)
	if err != nil {
		logger.Error("failed to build message request: %+v", err.Error())
		return err
	}
	msgRequest.Message.Token = []string{msg.DeviceToken}
	resp, err := cli.SendMessage(context.Background(), msgRequest)
	if err != nil {
		logger.Error("failed to send message: %+v", err.Error())
		return nil
	}
	if resp.Code != constant.Success {
		logger.Error("failed to send message: %+v", resp)
		return fmt.Errorf("HMS push failed: %s", resp.Msg)
	}
	return nil
}

func buildHmsMsgRequest(msg *model.PushMessage) (*hms_model.MessageRequest, error) {
	if msg.Title == "" {
		msg.Title = "Bark"
	}
	params := make(map[string]interface{})
	for k, v := range msg.ExtParams {
		// Change all parameter names to lowercase to prevent inconsistent capitalization
		params[strings.ToLower(k)] = fmt.Sprintf("%v", v)
	}
	// 解析额外参数
	params["category"] = msg.Category
	params["title"] = msg.Title
	params["body"] = msg.Body
	params["time"] = time.Now().UnixNano() / 1e6
	paramsBytes, err := json.Marshal(params)
	if err != nil {
		logrus.Error("Failed to marshal the params! Error is %s\n", err.Error())
		return nil, err
	}
	paramsStr := string(paramsBytes)

	msgRequest := hms_model.NewNotificationMsgRequest()
	msgRequest.Message.Data = paramsStr
	// 通用通知
	notification := &hms_model.Notification{
		Title: msg.Title,
		Body:  msg.Body,
	}
	msgRequest.Message.Notification = notification
	// Android 通知
	msgRequest.Message.Android = hms_model.GetDefaultAndroid()
	msgRequest.Message.Android.FastAppTarget = constant.FastAppTargetDevelop

	// 解析额外参数
	msgRequest.Message.Android.Data = paramsStr

	// 安卓通知内容
	androidNotification := hms_model.GetDefaultAndroidNotification()
	androidNotification.Title = msg.Title
	androidNotification.Body = msg.Body
	androidNotification.ForegroundShow = false
	androidNotification.Importance = constant.NotificationPriorityDefault
	androidNotification.ChannelId = "RingRing"
	androidNotification.ClickAction = &hms_model.ClickAction{
		Type:   constant.TypeIntentOrAction,
		Intent: "intent://com.fansy.cloud.bark/push?#Intent;scheme=push;launchFlags=0x4000000;end",
	}
	msgRequest.Message.Android.Notification = androidNotification

	return msgRequest, nil
}
