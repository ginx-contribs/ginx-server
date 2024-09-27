package main

import (
	"context"
	"fmt"
	"github.com/ginx-contribs/ginx"
	"github.com/ginx-contribs/ginx-server/internal/conf"
	"github.com/ginx-contribs/ginx-server/internal/server"
	"github.com/ginx-contribs/ginx-server/pkg/logh"
	"github.com/spf13/cobra"
	"log/slog"
)

var (
	Author     = "ginx-contribs"
	Version    = "unknown"
	BuildTime  = "1970.01.01"
	ConfigFile = "conf.toml"
)

var rootCmd = &cobra.Command{
	Use:          "ginx-server [commands] [-flags]",
	Short:        "ginx-server is quickstart for a http api server",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		httpserver, err := NewServer(ctx, Author, Version, BuildTime, ConfigFile)
		if err != nil {
			return err
		}
		return httpserver.Spin()
	},
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&ConfigFile, "config", "f", "conf.toml", "server configuration file")
	rootCmd.AddCommand(versionCmd)
}

func main() {
	rootCmd.Execute()
}

func NewServer(ctx context.Context, author, version, buildTime, configFile string) (*ginx.Server, error) {
	// read config file
	appConf, err := conf.ReadFrom(configFile)
	if err != nil {
		return nil, err
	}

	appConf.Author = author
	appConf.Version = version
	appConf.BuildTime = buildTime

	// revise configuration
	appConf, err = conf.Revise(appConf)
	if err != nil {
		return nil, err
	}

	// initialize app logger
	logger, err := logh.NewLogger(appConf.Log)
	if err != nil {
		return nil, err
	}
	defer logger.Close()

	// setup default logger
	slog.SetDefault(logger.Slog())
	slog.Info(fmt.Sprintf("logging in level: %s", appConf.Log.Level.String()))

	// initialize app
	app, err := server.NewHTTPServer(ctx, &appConf, logger)
	if err != nil {
		return nil, err
	}
	return app, nil
}
