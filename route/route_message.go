package route

import (
	"bark-server/model"
	"bark-server/orm"
	"github.com/gofiber/fiber/v2"
)

func init() {
	registerRoute("message", func(router *fiber.App) {
		router.Get("/message", func(c *fiber.Ctx) error { return getMessageList(c) })
	})
}

func getMessageList(c *fiber.Ctx) error {
	page := model.Page{}
	c.QueryParser(&page)
	page, err := orm.GetMessageList(page)
	if err != nil {
		return c.Status(400).JSON(failed(400, "get message list failed: %v", err))
	}
	return c.JSON(dataPage(page))
}
