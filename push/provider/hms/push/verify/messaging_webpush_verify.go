package verify

import (
	"errors"

	"bark-server/push/provider/hms/push/constant"
	"bark-server/push/provider/hms/push/model"
)

func validateWebPushConfig(webPushConfig *model.WebPushConfig) error {
	if webPushConfig == nil {
		return nil
	}

	if err := validateWebPushHeaders(webPushConfig.Headers); err != nil {
		return err
	}

	return validateWebPushNotification(webPushConfig.Notification)
}

func validateWebPushHeaders(headers *model.WebPushHeaders) error {
	if headers == nil {
		return nil
	}

	if headers.TTL != "" && !ttlPattern.MatchString(headers.TTL) {
		return errors.New("malformed ttl")
	}

	if headers.Urgency != "" &&
		headers.Urgency != constant.UrgencyHigh &&
		headers.Urgency != constant.UrgencyNormal &&
		headers.Urgency != constant.UrgencyLow &&
		headers.Urgency != constant.UrgencyVeryLow {
		return errors.New("priority must be 'high', 'normal', 'low' or 'very-low'")
	}
	return nil
}

func validateWebPushNotification(notification *model.WebPushNotification) error {
	if notification == nil {
		return nil
	}

	if err := validateWebPushAction(notification.Actions); err != nil {
		return err
	}

	if err := validateWebPushDirection(notification.Dir); err != nil {
		return err
	}
	return nil
}

func validateWebPushAction(actions []*model.WebPushAction) error {
	if actions == nil {
		return nil
	}

	for _, action := range actions {
		if action.Action == "" {
			return errors.New("web common action can't be empty")
		}
	}
	return nil
}

func validateWebPushDirection(dir string) error {
	if dir != constant.DirAuto && dir != constant.DirLtr && dir != constant.DirRtl {
		return errors.New("web common dir must be 'auto', 'ltr', 'rtl'")
	}
	return nil
}
