package route

import "github.com/gofiber/fiber/v2"

func init() {
	registerRoute("user", func(router *fiber.App) {
		router.Get("/user/auth", func(c *fiber.Ctx) error { return authUser(c) })
	})
}

func authUser(c *fiber.Ctx) error {
	return c.JSON(success())
}
