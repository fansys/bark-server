package main

import (
	"bark-server/misc"
	"bark-server/orm"
	"bark-server/route"
	"fmt"
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

func main() {
	app := &cli.App{
		Name:    "bark-server",
		Usage:   "Push Server For Bark",
		Version: fmt.Sprintf("%s %s %s", version, commitID, buildDate),
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "addr",
				Usage:   "Server listen address",
				EnvVars: []string{"BARK_SERVER_ADDRESS"},
				Value:   "0.0.0.0:8080",
			},
			&cli.StringFlag{
				Name:    "data",
				Usage:   "Server data storage dir",
				EnvVars: []string{"BARK_SERVER_DATA_DIR"},
				Value:   "/data",
			},
			&cli.StringFlag{
				Name:    "cert",
				Usage:   "Server TLS certificate",
				EnvVars: []string{"BARK_SERVER_CERT"},
				Value:   "",
			},
			&cli.StringFlag{
				Name:    "key",
				Usage:   "Server TLS certificate key",
				EnvVars: []string{"BARK_SERVER_KEY"},
				Value:   "",
			},
			&cli.BoolFlag{
				Name:    "case-sensitive",
				Usage:   "Enable HTTP URL case sensitive",
				EnvVars: []string{"BARK_SERVER_CASE_SENSITIVE"},
				Value:   false,
			},
			&cli.BoolFlag{
				Name:    "strict-routing",
				Usage:   "Enable strict routing distinction",
				EnvVars: []string{"BARK_SERVER_STRICT_ROUTING"},
				Value:   false,
			},
			&cli.BoolFlag{
				Name:    "reduce-memory-usage",
				Usage:   "Aggressively reduces memory usage at the cost of higher CPU usage if set to true",
				EnvVars: []string{"BARK_SERVER_REDUCE_MEMORY_USAGE"},
				Value:   false,
			},
			&cli.StringFlag{
				Name:    "user",
				Usage:   "Basic auth username",
				EnvVars: []string{"BARK_SERVER_BASIC_AUTH_USER"},
				Value:   "",
			},
			&cli.StringFlag{
				Name:    "password",
				Usage:   "Basic auth password",
				EnvVars: []string{"BARK_SERVER_BASIC_AUTH_PASSWORD"},
				Value:   "",
			},
			&cli.StringFlag{
				Name:    "proxy-header",
				Usage:   "The remote IP address used by the bark server http header",
				EnvVars: []string{"BARK_SERVER_PROXY_HEADER"},
				Value:   "",
			},
			&cli.StringFlag{
				Name:    "db-type",
				Usage:   "Server database type, use mysql or sqlite",
				EnvVars: []string{"DB_TYPE"},
				Value:   "sqlite",
			},
			&cli.StringFlag{
				Name:    "sqlite-path",
				Usage:   "Database Sqlite path",
				EnvVars: []string{"SQLITE_PATH"},
				Value:   "bark.db",
			},
			&cli.StringFlag{
				Name:    "mysql-host",
				Usage:   "Database MySQL host",
				EnvVars: []string{"MYSQL_HOST"},
				Value:   "127.0.0.1",
			},
			&cli.IntFlag{
				Name:    "mysql-port",
				Usage:   "Database MySQL port",
				EnvVars: []string{"MYSQL_PORT"},
				Value:   3306,
			},
			&cli.StringFlag{
				Name:    "mysql-user",
				Usage:   "Database MySQL user",
				EnvVars: []string{"MYSQL_USER"},
				Value:   "root",
			},
			&cli.StringFlag{
				Name:    "mysql-pass",
				Usage:   "Database MySQL pass",
				EnvVars: []string{"MYSQL_PASS"},
				Value:   "root",
			},
			&cli.StringFlag{
				Name:    "mysql-database",
				Usage:   "Database MySQL database",
				EnvVars: []string{"MYSQL_DATABASE"},
				Value:   "bark",
			},
			&cli.StringFlag{
				Name:    "mysql-params",
				Usage:   "Database MySQL parameters",
				EnvVars: []string{"MYSQL_PARAMS"},
				Value:   "charset=utf8mb4&parseTime=True&loc=Local",
			},
			&cli.IntFlag{
				Name:    "concurrency",
				Usage:   "Maximum number of concurrent connections",
				EnvVars: []string{"BARK_SERVER_CONCURRENCY"},
				Value:   256 * 1024,
				Hidden:  true,
			},
			&cli.DurationFlag{
				Name:    "read-timeout",
				Usage:   "The amount of time allowed to read the full request, including the body",
				EnvVars: []string{"BARK_SERVER_READ_TIMEOUT"},
				Value:   3 * time.Second,
				Hidden:  true,
			},
			&cli.DurationFlag{
				Name:    "write-timeout",
				Usage:   "The maximum duration before timing out writes of the response",
				EnvVars: []string{"BARK_SERVER_WRITE_TIMEOUT"},
				Value:   3 * time.Second,
				Hidden:  true,
			},
			&cli.DurationFlag{
				Name:    "idle-timeout",
				Usage:   "The maximum amount of time to wait for the next request when keep-alive is enabled",
				EnvVars: []string{"BARK_SERVER_IDLE_TIMEOUT"},
				Value:   10 * time.Second,
				Hidden:  true,
			},
		},
		Authors: []*cli.Author{
			{Name: "mritd", Email: "mritd@linux.com"},
			{Name: "Finb", Email: "to@day.app"},
		},
		Action: func(c *cli.Context) error {
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
					return c.Status(code).JSON(route.CommonResp{
						Code:      code,
						Message:   err.Error(),
						Timestamp: time.Now().Unix(),
					})
				},
			})

			misc.InitVersion(version, buildDate, commitID)
			route.RouterAuth(c.String("user"), c.String("password"), fiberApp)
			route.RouterSetup(fiberApp)
			orm.GormSetup(c)

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
