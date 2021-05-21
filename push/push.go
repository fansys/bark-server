package push

import (
	"bark-server/model"
	"bark-server/push/provider/apns"
	"bark-server/push/provider/hms"
	"fmt"
	"github.com/sideshow/apns2/payload"
	"strings"
)

func Push(msg *model.PushMessage) error {
	pl := payload.NewPayload().
		AlertTitle(msg.Title).
		AlertBody(msg.Body).
		Sound(msg.Sound).
		Category(msg.Category)

	for k, v := range msg.ExtParams {
		// Change all parameter names to lowercase to prevent inconsistent capitalization
		pl.Custom(strings.ToLower(k), fmt.Sprintf("%v", v))
	}
	var isIos = false
	if strings.Compare(msg.DeviceToken, strings.ToLower(msg.DeviceToken)) == 0 {
		isIos = true
	}
	if isIos {
		return apns.Push(msg)
	} else {
		return hms.Push(msg)
	}
}
