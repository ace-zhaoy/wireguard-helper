package main

import (
	"context"
	"fmt"
	"github.com/ace-zhaoy/errors"
	"github.com/ace-zhaoy/glog/log"
	"github.com/ace-zhaoy/go-utils/usignal"
	"github.com/ace-zhaoy/gviper"
	"github.com/ace-zhaoy/wireguard-helper/internal/config"
	"github.com/ace-zhaoy/wireguard-helper/internal/detection"
	config2 "github.com/ace-zhaoy/wireguard-helper/internal/detection/config"
	"github.com/ace-zhaoy/wireguard-helper/internal/logs"
	"github.com/ace-zhaoy/wireguard-helper/internal/tunnelmanager"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var (
	logLevel      string
	detectionObj  = detection.NewDetection()
	tunnelManager *tunnelmanager.TunnelManager
	tunnelNames   []string
)

func main() {
	ctx, cancel := usignal.WithSignalContext(context.Background())
	defer errors.Recover(func(e error) {
		log.ErrorContext(ctx, fmt.Sprintf("panic: %+v", e))
		cancel()
		os.Exit(2)
	})

	var configPath string
	rootCmd := &cobra.Command{}
	rootCmd.Flags().StringVarP(&configPath, "config", "c", "./config", "config path")
	rootCmd.Flags().StringVarP(&logLevel, "log-level", "l", "", "log level, options: debug, info, warn, error")
	rootCmd.Flags().StringSliceVarP(&tunnelNames, "tunnel-name", "t", nil, "tunnel name, if not set, all tunnels will be used")
	rootCmd.RunE = func(cmd *cobra.Command, args []string) (err error) {
		defer errors.Recover(func(e error) { err = errors.Wrap(e, "param: %v", args) })
		errors.Check(InitWithConfigPath(ctx, configPath))
		log.Info("connect start")
		err = tunnelManager.Connect(ctx)
		errors.Check(err)
		log.Info("connect end")
		return
	}

	err := rootCmd.Execute()
	errors.Check(errors.Wrap(err, "cmd execute error"))
}

func InitWithConfigPath(ctx context.Context, configPath string) (err error) {
	defer errors.Recover(func(e error) { err = errors.Wrap(e, "param: %v", configPath) })
	conf := gviper.Default(configPath)

	conf.RegisterNotification(config.NewLogHook())

	InitLogWithConfig(ctx, conf, logLevel)

	InitDetectionWithConfig(ctx, conf)

	InitTunnelManagerWithConfig(ctx, conf)

	err = conf.Load()
	errors.Check(err)

	conf.Watch()
	return
}

func InitLogWithConfig(_ context.Context, config *gviper.Config, logLevel string) {
	var lc logs.Config
	config.BindAndListen("log", &lc, func(_ *viper.Viper) error {
		if logLevel != "" {
			lc.Level = logLevel
			logLevel = ""
		}
		return logs.Init(lc)
	})
}

func InitDetectionWithConfig(_ context.Context, config *gviper.Config) {
	var dc config2.Config
	config.BindAndListen("detection", &dc, func(_ *viper.Viper) error {
		return detectionObj.LoadConfig(dc)
	})
}

func InitTunnelManagerWithConfig(_ context.Context, config *gviper.Config) {
	var c tunnelmanager.Config
	config.BindAndListen("tunnel_manager", &c, func(_ *viper.Viper) (err error) {
		defer errors.Recover(func(e error) { err = e })
		if len(tunnelNames) > 0 {
			c.ConnectTunnelNames = tunnelNames
		}
		if tunnelManager == nil {
			tunnelManager, err = tunnelmanager.NewTunnelManager(c, detectionObj)
			errors.Check(err)
			return
		}
		errors.Check(tunnelManager.LoadConfig(c))
		return
	})
}
