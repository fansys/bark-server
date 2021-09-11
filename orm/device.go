package orm

import (
	"errors"
	"fansys/bark-server/v2/util"
	"github.com/lithammer/shortuuid/v3"
	"github.com/mritd/logger"
	"strings"
)

func GetDeviceByKey(deviceKey string) (*Device, error) {
	if deviceKey == "" {
		return nil, errors.New("device key is empty")
	}
	var device *Device
	res := db.Where("device_key = ?", deviceKey).First(&device)
	if res.Error != nil {
		logger.Errorf("get device by key error. token: %v, error: %v", deviceKey, res.Error)
		return nil, res.Error
	}
	return device, nil
}

func GetDeviceByToken(deviceToken string) (*Device, error) {
	if deviceToken == "" {
		return nil, errors.New("device token is empty")
	}
	var device *Device
	res := db.Where("device_token = ?", deviceToken).First(&device)
	if res.Error != nil {
		logger.Errorf("get device by token error. token: %v, error: %v", deviceToken, res.Error)
		return nil, res.Error
	}
	return device, nil
}

func SaveDevice(deviceKey, deviceToken, deviceType string) (*Device, error) {
	if deviceToken == "" {
		return nil, errors.New("device token is empty")
	}
	if deviceType == "" {
		deviceType = "ios"
	}
	device, _ := GetDeviceByToken(deviceToken)
	if device != nil {
		if device.DeviceKey == deviceKey {
			return device, nil
		}
		logger.Infof("update device key. token: %v, key: %v", device.DeviceToken, device.DeviceKey)
		// update deviceKey
		device.DeviceKey = genKey()
		db.Model(Device{}).Where("device_token = ?", deviceToken).Update("device_key", device.DeviceKey)
		return device, nil
	} else {
		device := Device{
			DeviceKey:   genKey(),
			DeviceToken: deviceToken,
			DeviceType:  deviceType,
		}
		res := db.Create(&device)
		if res.Error != nil {
			logger.Errorf("save device error. token: %v, error: %v", device.DeviceToken, res.Error)
			return nil, res.Error
		}
		logger.Infof("save device success. token: %v, key: %v", device.DeviceToken, device.DeviceKey)
		return &device, nil
	}
}

func CountDevice() int {
	var count int64
	db.Model(&Device{}).Count(&count)
	return util.Int64ToInt(count)
}

func genKey() string {
	return strings.ToLower(shortuuid.New())
}
