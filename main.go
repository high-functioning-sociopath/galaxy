package main

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

const (
	loglevelDebug = "debug"
	loglevelInfo  = "info"
	logLevelError = "error"
)

func main() {

	// Setup logrus
	logrus.SetFormatter(&logrus.TextFormatter{})
	logrus.SetOutput(os.Stdout)

	app := cli.NewApp()
	app.Name = "galaxy"
	app.Version = "0.1.0"

	app.Commands = []cli.Command{
		{
			Name:  "runner",
			Usage: "Starts a galaxy runner instance",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "port",
					EnvVar: "PORT",
					Usage:  "The port the runner will bind too",
					Value:  "4050",
				},
				cli.StringFlag{
					Name:   "proxy-port",
					EnvVar: "PROXY_PORT",
					Usage:  "The port the proxy will bind too",
					Value:  "4055",
				},
				cli.StringFlag{
					Name:   "log-level",
					EnvVar: "LOG_LEVEL",
					Usage:  "Set the log level [debug | info | error]",
					Value:  loglevelInfo,
				},

				// JWT config
				cli.StringFlag{
					Name:   "jwt-algo",
					EnvVar: "JWT_ALGO",
					Usage:  "The jwt algorithm to use for verification and signing [ hs256 | rsa256 ]",
					Value:  "hs256",
				},
				cli.StringFlag{
					Name:   "jwt-secret",
					EnvVar: "JWT_SECRET",
					Usage:  "The jwt secret to use when the algorithm is set to HS256",
					Value:  "some-secret",
				},
				cli.StringFlag{
					Name:   "jwt-proxy-secret",
					EnvVar: "JWT_PROXY_SECRET",
					Usage:  "The jwt secret to use for authenticating the proxy",
					Value:  "some-proxy-secret",
				},

				// Driver config
				cli.StringFlag{
					Name:   "driver",
					EnvVar: "DRIVER",
					Usage:  "The driver to use for deployment",
					Value:  "istio",
				},
				cli.StringFlag{
					Name:   "driver-config",
					EnvVar: "DRIVER_CONFIG",
					Usage:  "Driver config file path",
				},
				cli.BoolFlag{
					Name:   "outside-cluster",
					EnvVar: "OUTSIDE_CLUSTER",
					Usage:  "Indicates whether galaxy in running inside the cluster",
				},
			},
			Action: actionRunner,
		},
		{
			Name:  "proxy",
			Usage: "Starts the proxy to collect metrics directly from envoy",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "addr",
					Usage:  "Address of the galaxy runner instance",
					EnvVar: "ADDR",
					Value:  "runner.galaxy.svc.cluster.local:4050",
				},
				cli.StringFlag{
					Name:   "token",
					Usage:  "The token to be used for authentication",
					EnvVar: "TOKEN",
				},
				cli.StringFlag{
					Name:   "log-level",
					EnvVar: "LOG_LEVEL",
					Usage:  "Set the log level [debug | info | error]",
					Value:  loglevelInfo,
				},
			},
			Action: actionProxy,
		},
		{
			Name:  "server",
			Usage: "Starts the galaxy server instance",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "port",
					Usage:  "The port the server will bind to",
					EnvVar: "PORT",
					Value:  "4050",
				},
				cli.StringFlag{
					Name:   "log-level",
					EnvVar: "LOG_LEVEL",
					Usage:  "Set the log level [debug | info | error]",
					Value:  loglevelInfo,
				},
			},
			Action: actionServer,
		},
	}

	// Start the app
	if err := app.Run(os.Args); err != nil {
		logrus.Fatalln("Failed to start galaxy:", err)
	}
}
