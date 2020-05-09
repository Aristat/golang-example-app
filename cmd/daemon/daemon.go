package daemon

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/aristat/golang-example-app/app/entrypoint"
	"github.com/aristat/golang-example-app/app/http"
	"github.com/aristat/golang-example-app/app/logger"

	"github.com/spf13/cobra"
)

var (
	bind          string
	gracefulDelay time.Duration
	Cmd           = &cobra.Command{
		Use:           "daemon",
		Short:         "Gateway API daemon",
		SilenceUsage:  true,
		SilenceErrors: true,
		Run: func(_ *cobra.Command, _ []string) {
			var (
				e error
				s *http.Http
				c func()
			)

			log, c, e := logger.Build()
			if e != nil {
				panic(e)
			}
			defer c()
			defer func() {
				if r := recover(); r != nil {
					if re, _ := r.(error); re != nil {
						log.Error(re.Error())
					} else {
						log.Alert("unhandled panic, err: %v", logger.Args(r))
					}
				}
			}()

			s, c, e = http.Build()
			if e != nil {
				log.Error(e.Error())
				return
			}
			defer c()

			wg := &sync.WaitGroup{}
			wg.Add(1)

			server := s.ListenAndServe(wg, bind)

			shutdownSignal := make(chan os.Signal)
			signal.Notify(shutdownSignal, syscall.SIGTERM, syscall.SIGINT)
			sig := <-shutdownSignal
			log.Printf("OS signaled `%v`\n", sig.String())

			log.Info("Server shutdown is raised")
			if e := server.Shutdown(context.Background()); e != nil {
				log.Emergency("Graceful shutdown error, %v", logger.Args(e))
			}

			wg.Wait()

			log.Printf("Graceful shutdown in %s\n", gracefulDelay)
			ctx, _ := context.WithTimeout(context.Background(), gracefulDelay)
			entrypoint.Shutdown(ctx, 0)
		},
	}
)

func init() {
	Cmd.PersistentFlags().StringVarP(&bind, "bind", "b", ":9096", "bind address")
	Cmd.PersistentFlags().DurationVar(&gracefulDelay, "graceful.delay", 50*time.Millisecond, "graceful delay")
}
