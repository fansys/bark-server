package main

import (
	"fansys/bark-server/v2/orm"
	apns2 "fansys/bark-server/v2/push/apns"
	"fansys/bark-server/v2/push/getui"
	"github.com/mritd/logger"
	"net/url"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func init() {
	registerRoute("push", func(router *fiber.App) {
		router.Post("/push", func(c *fiber.Ctx) error { return routeDoPush(c, false) })
	})

	// compatible with old requests
	registerRouteWithWeight("push_compat", 1, func(router *fiber.App) {
		router.Get("/:device_key", func(c *fiber.Ctx) error { return routeDoPush(c, true) })
		router.Post("/:device_key", func(c *fiber.Ctx) error { return routeDoPush(c, true) })

		router.Get("/:device_key/:body", func(c *fiber.Ctx) error { return routeDoPush(c, true) })
		router.Post("/:device_key/:body", func(c *fiber.Ctx) error { return routeDoPush(c, true) })

		router.Get("/:device_key/:title/:body", func(c *fiber.Ctx) error { return routeDoPush(c, true) })
		router.Post("/:device_key/:title/:body", func(c *fiber.Ctx) error { return routeDoPush(c, true) })

		router.Get("/:device_key/:category/:title/:body", func(c *fiber.Ctx) error { return routeDoPush(c, true) })
		router.Post("/:device_key/:category/:title/:body", func(c *fiber.Ctx) error { return routeDoPush(c, true) })
	})
}

func routeDoPush(c *fiber.Ctx, compat bool) error {
	// default value
	msg := orm.PushMessage{
		Category:  "myNotificationCategory",
		Body:      "NoContent",
		Sound:     "1107",
		ExtParams: make(map[string]interface{}),
	}
	// always parse body(Lowest priority)
	if err := c.BodyParser(&msg); err != nil && err != fiber.ErrUnprocessableEntity {
		return c.Status(400).JSON(failed(400, "request bind failed: %v", err))
	}

	if compat {
		params := make(map[string]string)
		visitor := func(key, value []byte) {
			params[strings.ToLower(string(key))] = string(value)
		}
		// parse query args (medium priority)
		c.Request().URI().QueryArgs().VisitAll(visitor)
		// parse post args
		c.Request().PostArgs().VisitAll(visitor)

		// parse multipartForm values
		form, err := c.Request().MultipartForm()
		if err == nil {
			for key, val := range form.Value {
				if len(val) > 0 {
					params[key] = val[0]
				}
			}
		}

		for key, val := range params {
			switch strings.ToLower(string(key)) {
			case "device_key":
				msg.DeviceKey = val
			case "category":
				msg.Category = val
			case "title":
				msg.Title = val
			case "body":
				msg.Body = val
			case "sound":
				msg.Sound = val + ".caf"
			case "group":
				// 服务端使用 msg.Group 设置 ThreadID, 对通知中心的推送进行分组
				msg.Group = val
				// 客户端从 Custom payload 中拿到 group 参数进行分组
				msg.ExtParams[strings.ToLower(string(key))] = val
			default:
				msg.ExtParams[strings.ToLower(string(key))] = val
			}
		}

		// parse url path (highest priority)
		if pathDeviceKey := c.Params("device_key"); pathDeviceKey != "" {
			msg.DeviceKey = pathDeviceKey
		}
		if category := c.Params("category"); category != "" {
			str, err := url.QueryUnescape(category)
			if err != nil {
				return err
			}
			msg.Category = str
		}
		if title := c.Params("title"); title != "" {
			str, err := url.QueryUnescape(title)
			if err != nil {
				return err
			}
			msg.Title = str
		}
		if body := c.Params("body"); body != "" {
			str, err := url.QueryUnescape(body)
			if err != nil {
				return err
			}
			msg.Body = str
		}
	}

	if msg.DeviceKey == "" {
		return c.Status(400).JSON(failed(400, "device key is empty"))
	}

	device, err := orm.GetDeviceByKey(msg.DeviceKey)
	if err != nil {
		return c.Status(400).JSON(failed(400, "failed to get device token: %v", err))
	}
	msg.DeviceToken = device.DeviceToken
	msg.DeviceType = device.DeviceType
	orm.SaveMessage(&msg)
	logger.Infof("bark push message, device key: %v, token: %v", device.DeviceKey, device.DeviceToken)
	if device.DeviceType == "ios" {
		err = apns2.Push(&msg)
	} else {
		err = getui.Push(&msg)
	}
	if err != nil {
		logger.Errorf("bark push message, device key: %v, error: %v", device.DeviceKey, err)
		return c.Status(500).JSON(failed(500, "push failed: %v", err))
	}
	logger.Infof("bark push message success, device key: %v", device.DeviceKey)
	return c.JSON(success())
}
