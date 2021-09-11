package main

import (
	"fansys/bark-server/v2/orm"

	"github.com/gofiber/fiber/v2"
	"github.com/mritd/logger"
)

type DeviceInfo struct {
	DeviceKey   string `form:"device_key,omitempty" json:"device_key,omitempty" xml:"device_key,omitempty" query:"device_key,omitempty"`
	DeviceToken string `form:"device_token,omitempty" json:"device_token,omitempty" xml:"device_token,omitempty" query:"device_token,omitempty"`
	DeviceType  string `form:"device_type,omitempty" json:"device_type,omitempty" xml:"device_type,omitempty" query:"device_type,omitempty"`

	// compatible with old req
	OldDeviceKey   string `form:"key,omitempty" json:"key,omitempty" xml:"key,omitempty" query:"key,omitempty"`
	OldDeviceToken string `form:"devicetoken,omitempty" json:"devicetoken,omitempty" xml:"devicetoken,omitempty" query:"devicetoken,omitempty"`
}

const (
	bucketName = "device"
)

func init() {
	registerRoute("register", func(router *fiber.App) {
		router.Post("/register", func(c *fiber.Ctx) error { return doRegister(c, false) })
		router.Get("/register/:device_key", doRegisterCheck)
	})

	// compatible with old requests
	registerRouteWithWeight("register_compat", 100, func(router *fiber.App) {
		router.Get("/register", func(c *fiber.Ctx) error { return doRegister(c, true) })
	})
}

func doRegister(c *fiber.Ctx, compat bool) error {
	var deviceInfo DeviceInfo
	if compat {
		if err := c.QueryParser(&deviceInfo); err != nil {
			return c.Status(400).JSON(failed(400, "request bind failed: %v", err))
		}
	} else {
		if err := c.BodyParser(&deviceInfo); err != nil {
			return c.Status(400).JSON(failed(400, "request bind failed: %v", err))
		}
	}

	if deviceInfo.DeviceKey == "" && deviceInfo.OldDeviceKey != "" {
		deviceInfo.DeviceKey = deviceInfo.OldDeviceKey
	}

	if deviceInfo.DeviceToken == "" {
		if deviceInfo.OldDeviceToken != "" {
			deviceInfo.DeviceToken = deviceInfo.OldDeviceToken
		} else {
			return c.Status(400).JSON(failed(400, "device token is empty"))
		}
	}

	device, err := orm.SaveDevice(deviceInfo.DeviceKey, deviceInfo.DeviceToken, deviceInfo.DeviceType)

	if err != nil {
		logger.Errorf("device registration failed: %v", err)
		return c.Status(500).JSON(failed(500, "device registration failed: %v", err))
	}

	return c.Status(200).JSON(data(map[string]string{
		// compatible with old resp
		"key":          device.DeviceKey,
		"device_key":   device.DeviceKey,
		"device_token": device.DeviceToken,
		"device_type":  device.DeviceType,
	}))
}

func doRegisterCheck(c *fiber.Ctx) error {
	deviceKey := c.Params("device_key")

	if deviceKey == "" {
		return c.Status(400).JSON(failed(400, "device key is empty"))
	}

	_, err := orm.GetDeviceByKey(deviceKey)
	if err != nil {
		return c.Status(400).JSON(failed(400, err.Error()))
	}
	return c.Status(200).JSON(success())
}
