package main

import (
	"errors"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	"github.com/spaceuptech/galaxy/model"
	"github.com/spaceuptech/galaxy/proxy"
	"github.com/spaceuptech/galaxy/runner"
	"github.com/spaceuptech/galaxy/runner/driver"
	"github.com/spaceuptech/galaxy/server"
	"github.com/spaceuptech/galaxy/utils/auth"
)

func actionRunner(c *cli.Context) error {
	// Get runner config flags
	port := c.String("port")
	proxyPort := c.String("proxy-port")
	loglevel := c.String("log-level")

	// Get jwt config
	jwtAlgo := auth.JWTAlgorithm(c.String("jwt-algo"))
	jwtSecret := c.String("jwt-secret")
	jwtProxySecret := c.String("jwt-proxy-secret")

	// Get driver config
	driverType := c.String("driver")
	driverConfig := c.String("driver-config")
	outsideCluster := c.Bool("outside-cluster")

	// Set the log level
	setLogLevel(loglevel)

	// Create a new runner object
	r, err := runner.New(&runner.Config{
		Port:      port,
		ProxyPort: proxyPort,
		Auth: &auth.Config{
			Mode:         auth.Runner,
			JWTAlgorithm: jwtAlgo,
			Secret:       jwtSecret,
			ProxySecret:  jwtProxySecret,
		},
		Driver: &driver.Config{
			DriverType:     model.DriverType(driverType),
			ConfigFilePath: driverConfig,
			IsInCluster:    !outsideCluster,
		},
	})
	if err != nil {
		logrus.Errorf("Failed to start runner - %s", err.Error())
		os.Exit(-1)
	}

	return r.Start()
}

func actionProxy(c *cli.Context) error {
	// Get all flags
	addr := c.String("addr")
	token := c.String("token")
	loglevel := c.String("log-level")

	// Set the log level
	setLogLevel(loglevel)

	// Throw an error if invalid token provided
	if len(strings.Split(token, ".")) != 3 {
		return errors.New("invalid token provided")
	}

	// Start the proxy
	p := proxy.New(addr, token)
	return p.Start()
}

func actionServer(c *cli.Context) error {
	// Get server config flags
	port := c.String("port")
	loglevel := c.String("log-level")

	// Set the log level
	setLogLevel(loglevel)

	s := server.New(&server.Config{Port: port})
	return s.Start()
}

func setLogLevel(loglevel string) {
	switch loglevel {
	case loglevelDebug:
		logrus.SetLevel(logrus.DebugLevel)
	case loglevelInfo:
		logrus.SetLevel(logrus.InfoLevel)
	case logLevelError:
		logrus.SetLevel(logrus.ErrorLevel)
	default:
		logrus.Errorf("Invalid log level (%s) provided", loglevel)
		logrus.Infoln("Defaulting to `info` level")
		logrus.SetLevel(logrus.InfoLevel)
	}
}
