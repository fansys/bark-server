package main

import (
	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
	"time"
)

var flags []cli.Flag

func init() {
	flags = []cli.Flag{
		&cli.StringFlag{
			Name:       ConfigName,
			Usage:      "Server config dir",
			EnvVars:    []string{"BARK_CONFIG_DIR"},
			Value:      DefaultConfig,
			HasBeenSet: true,
		},
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:    "addr",
			Usage:   "Server listen address",
			EnvVars: []string{"BARK_SERVER_ADDRESS"},
			Value:   "0.0.0.0:8080",
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:    "data",
			Usage:   "Server data storage dir",
			EnvVars: []string{"BARK_SERVER_DATA_DIR"},
			Value:   "data",
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:    "cert",
			Usage:   "Server TLS certificate",
			EnvVars: []string{"BARK_SERVER_CERT"},
			Value:   "",
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:    "key",
			Usage:   "Server TLS certificate key",
			EnvVars: []string{"BARK_SERVER_KEY"},
			Value:   "",
		}),
		altsrc.NewBoolFlag(&cli.BoolFlag{
			Name:    "case-sensitive",
			Usage:   "Enable HTTP URL case sensitive",
			EnvVars: []string{"BARK_SERVER_CASE_SENSITIVE"},
			Value:   false,
		}),
		altsrc.NewBoolFlag(&cli.BoolFlag{
			Name:    "strict-routing",
			Usage:   "Enable strict routing distinction",
			EnvVars: []string{"BARK_SERVER_STRICT_ROUTING"},
			Value:   false,
		}),
		altsrc.NewBoolFlag(&cli.BoolFlag{
			Name:    "reduce-memory-usage",
			Usage:   "Aggressively reduces memory usage at the cost of higher CPU usage if set to true",
			EnvVars: []string{"BARK_SERVER_REDUCE_MEMORY_USAGE"},
			Value:   false,
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:    "user",
			Usage:   "Basic auth username",
			EnvVars: []string{"BARK_SERVER_BASIC_AUTH_USER"},
			Value:   "",
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:    "password",
			Usage:   "Basic auth password",
			EnvVars: []string{"BARK_SERVER_BASIC_AUTH_PASSWORD"},
			Value:   "",
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:    "proxy-header",
			Usage:   "The remote IP address used by the bark server http header",
			EnvVars: []string{"BARK_SERVER_PROXY_HEADER"},
			Value:   "",
		}),
		altsrc.NewIntFlag(&cli.IntFlag{
			Name:    "concurrency",
			Usage:   "Maximum number of concurrent connections",
			EnvVars: []string{"BARK_SERVER_CONCURRENCY"},
			Value:   256 * 1024,
			Hidden:  true,
		}),
		altsrc.NewDurationFlag(&cli.DurationFlag{
			Name:    "read-timeout",
			Usage:   "The amount of time allowed to read the full request, including the body",
			EnvVars: []string{"BARK_SERVER_READ_TIMEOUT"},
			Value:   3 * time.Second,
			Hidden:  true,
		}),
		altsrc.NewDurationFlag(&cli.DurationFlag{
			Name:    "write-timeout",
			Usage:   "The maximum duration before timing out writes of the response",
			EnvVars: []string{"BARK_SERVER_WRITE_TIMEOUT"},
			Value:   3 * time.Second,
			Hidden:  true,
		}),
		altsrc.NewDurationFlag(&cli.DurationFlag{
			Name:    "idle-timeout",
			Usage:   "The maximum amount of time to wait for the next request when keep-alive is enabled",
			EnvVars: []string{"BARK_SERVER_IDLE_TIMEOUT"},
			Value:   10 * time.Second,
			Hidden:  true,
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:    "getui.app-id",
			Usage:   "GeTui appId",
			EnvVars: []string{"GETUI_APP_ID"},
			Value:   "",
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:    "getui.app-key",
			Usage:   "GeTui appKey",
			EnvVars: []string{"GETUI_APP_KEY"},
			Value:   "",
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:    "getui.master-secret",
			Usage:   "GeTui masterSecret",
			EnvVars: []string{"GETUI_MASTER_SECRET"},
			Value:   "",
		}),
	}
}
