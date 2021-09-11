package main

import (
	"fansys/bark-server/v2/orm"
	"fansys/bark-server/v2/push/getui"
	"fmt"
	"github.com/urfave/cli/v2/altsrc"
	"os"
	"os/signal"
	"syscall"
	"time"

	jsoniter "github.com/json-iterator/go"

	"github.com/gofiber/fiber/v2"

	"github.com/mritd/logger"
	"github.com/urfave/cli/v2"
)

var (
	version   string
	buildDate string
	commitID  string
)

const (
	ConfigName    = "config"
	DefaultConfig = "config.yml"
)

func main() {
	app := &cli.App{
		Name:    "bark-server",
		Usage:   "Push Server For Bark",
		Version: fmt.Sprintf("%s %s %s", version, commitID, buildDate),
		Before:  altsrc.InitInputSourceWithContext(flags, altsrc.NewYamlSourceFromFlagFunc(ConfigName)),
		Flags:   flags,
		Authors: []*cli.Author{
			{Name: "mritd", Email: "mritd@linux.com"},
			{Name: "Finb", Email: "to@day.app"},
		},
		Action: func(c *cli.Context) error {
			getui.New(getui.Config{
				AppId:        c.String("getui.app-id"),
				AppKey:       c.String("getui.app-key"),
				MasterSecret: c.String("getui.master-secret"),
			})
			orm.New(c.String("data"))
			fiberApp := fiber.New(fiber.Config{
				ServerHeader:      "Bark",
				CaseSensitive:     c.Bool("case-sensitive"),
				StrictRouting:     c.Bool("strict-routing"),
				Concurrency:       c.Int("concurrency"),
				ReadTimeout:       c.Duration("read-timeout"),
				WriteTimeout:      c.Duration("write-timeout"),
				IdleTimeout:       c.Duration("idle-timeout"),
				ProxyHeader:       c.String("proxy-header"),
				ReduceMemoryUsage: c.Bool("reduce-memory-usage"),
				JSONEncoder:       jsoniter.Marshal,
				ErrorHandler: func(c *fiber.Ctx, err error) error {
					code := fiber.StatusInternalServerError
					if e, ok := err.(*fiber.Error); ok {
						code = e.Code
					}
					return c.Status(code).JSON(CommonResp{
						Code:      code,
						Message:   err.Error(),
						Timestamp: time.Now().Unix(),
					})
				},
			})

			routerAuth(c.String("user"), c.String("password"), fiberApp)
			routerSetup(fiberApp)

			go func() {
				sigs := make(chan os.Signal)
				signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
				for range sigs {
					logger.Warn("Received a termination signal, bark server shutdown...")
					if err := fiberApp.Shutdown(); err != nil {
						logger.Errorf("Server forced to shutdown error: %v", err)
					}
				}
			}()

			logger.Infof("Bark Server Listen at: %s", c.String("addr"))
			if cert, key := c.String("cert"), c.String("key"); cert != "" && key != "" {
				return fiberApp.ListenTLS(c.String("addr"), cert, key)
			}
			return fiberApp.Listen(c.String("addr"))
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		logger.Error(err)
		os.Exit(1)
	}
}
