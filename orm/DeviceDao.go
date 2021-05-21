package orm

import (
	"bark-server/model"
	"bark-server/util"
	"errors"
	"github.com/lithammer/shortuuid/v3"
	"github.com/mritd/logger"
)

func GetDeviceTokenByDeviceKey(deviceKey string) (string, error) {
	var device model.Device
	var err error
	if err = db.Where(model.Device{
		DeviceKey: deviceKey,
		Enabled:   true,
	}, "device_key", "enabled").First(&device).Error; err != nil {
		logger.Error("query device by device key error. ", err)
		if err.Error() == "record not found" {
			err = errors.New("device not registered")
		}
		return "", err
	}
	return device.DeviceToken, nil
}

func GetDeviceByDeviceKey(deviceKey string) (model.Device, error) {
	var device model.Device
	var err error
	if err = db.Where(model.Device{
		DeviceKey: deviceKey,
		Enabled:   true,
	}, "device_key", "enabled").First(&device).Error; err != nil {
		if err.Error() != "record not found" {
			logger.Error("query device by device key error. %+v", err)
			return device, err
		}
	}
	return device, err
}

func GetDeviceByDeviceToken(deviceToken string) (model.Device, error) {
	var device model.Device
	var err error
	if err = db.Where(model.Device{
		DeviceToken: deviceToken,
		Enabled:     true,
	}, "device_token", "enabled").First(&device).Error; err != nil {
		if err.Error() != "record not found" {
			logger.Error("query device by device token error. %+v", err)
			return device, err
		}
	}
	return device, err
}

func SaveOrUpdateDevice(deviceToken string, deviceKey string) (model.Device, error) {
	var device model.Device
	var err error
	if err = db.Where(model.Device{
		DeviceToken: deviceToken,
	}, "device_token", "deleted").First(&device).Error; err != nil {
		if err.Error() != "record not found" {
			logger.Error("query device error. %+v", err)
			return device, err
		}
	}
	if deviceKey != "" {
		device.DeviceKey = deviceKey
	} else {
		device.DeviceKey = shortuuid.New()
	}
	//device.UpdateTime = time.Now()
	if device.ID == 0 {
		device.DeviceToken = deviceToken
		//device.CreateTime = time.Now()
		if err = db.Create(&device).Error; err != nil {
			logger.Error("create device error. %+v", err)
			return device, err
		}
	} else {
		if err = db.Model(&device).Update("device_key", device.DeviceKey).Error; err != nil {
			logger.Error("update device error. %+v", err)
			return device, err
		}
	}
	return device, err
}

func ChangeDeviceStatus(deviceToken string, enabled bool) (bool, error) {
	var device model.Device
	var err error
	if err = db.Where(model.Device{
		DeviceToken: deviceToken,
	}, "device_token", "deleted").First(&device).Error; err != nil {
		logger.Error("query device error. %+v", err)
		return false, err
	}
	if err = db.Model(&device).Update("enabled", enabled).Error; err != nil {
		logger.Error("update device status error. %+v", err)
		return false, err
	}
	return true, nil
}

func CountDevice() (int, error) {
	var count int64
	if err := db.Model(&model.Device{}).Count(&count).Error; err != nil {
		logger.Error("count device error. %+v", err)
		return 0, err
	}
	return util.Int64ToInt(count), nil
}
