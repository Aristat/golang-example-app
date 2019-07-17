package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/aristat/golang-gin-oauth2-example-app/cmd/client"

	"github.com/aristat/golang-gin-oauth2-example-app/app/entrypoint"
	"github.com/aristat/golang-gin-oauth2-example-app/app/logger"

	"go.uber.org/automaxprocs/maxprocs"

	"github.com/aristat/golang-gin-oauth2-example-app/cmd/daemon"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	configPath    string
	debug         bool
	v             *viper.Viper
	gracefulDelay time.Duration
	log           logger.Logger
)

const prefix = "cmd.root"

// Root command
var rootCmd = &cobra.Command{
	Use:           "bin [command]",
	Long:          "",
	SilenceUsage:  true,
	SilenceErrors: true,
	PersistentPreRun: func(cmd *cobra.Command, _ []string) {
		l, c, e := logger.Build()
		defer c()
		if e != nil {
			panic(e)
		}

		log = l.WithFields(logger.Fields{"service": prefix})

		v.SetConfigFile(configPath)

		if configPath != "" {
			e := v.ReadInConfig()
			if e != nil {
				log.Error("can't read config, %v", logger.Args(errors.WithMessage(e, prefix)))
				os.Exit(1)
			}
		}

		if debug {
			b, _ := json.Marshal(v.AllSettings())
			var out bytes.Buffer
			e := json.Indent(&out, b, "", "  ")
			if e != nil {
				log.Error("can't prettify config")
				os.Exit(1)
			}
			fmt.Println(string(out.Bytes()))
		}

		_, _ = maxprocs.Set(maxprocs.Logger(log.Printf))
	},
}

func init() {
	v = viper.New()
	v.SetConfigType("yaml")
	v.SetEnvPrefix("APP")
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_", ".", "_"))
	v.AutomaticEnv()

	// pflags
	rootCmd.PersistentFlags().StringVarP(&configPath, "config", "c", "", "config file")
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "debug mode")
	rootCmd.PersistentFlags().DurationVar(&gracefulDelay, "graceful.delay", 50*time.Millisecond, "graceful delay")

	// initializing
	wd := os.Getenv("APP_WD")
	if len(wd) == 0 {
		wd, _ = filepath.Abs(filepath.Dir(os.Args[0]))
	}
	wd, _ = filepath.Abs(wd)
	ep, _ := entrypoint.Initialize(wd, v)

	// bin pflags to viper
	_ = v.BindPFlags(rootCmd.PersistentFlags())

	go func() {
		reloadSignal := make(chan os.Signal)
		signal.Notify(reloadSignal, syscall.SIGHUP)
		for {
			sig := <-reloadSignal
			ep.Reload()
			fmt.Printf("OS signaled `%v`, reload", sig.String())
		}
	}()

	go func() {
		shutdownSignal := make(chan os.Signal)
		signal.Notify(shutdownSignal, syscall.SIGTERM, syscall.SIGINT)
		sig := <-shutdownSignal
		fmt.Printf("OS signaled `%v`, graceful shutdown in %s", sig.String(), gracefulDelay)
		ctx, _ := context.WithTimeout(context.Background(), gracefulDelay)
		ep.Shutdown(ctx, 0)
	}()
}

func Execute() {
	rootCmd.AddCommand(daemon.Cmd, client.Cmd)
	if e := rootCmd.Execute(); e != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v\n", e.Error())
		os.Exit(1)
	}
}
