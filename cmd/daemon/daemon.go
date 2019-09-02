package daemon

import (
	"github.com/aristat/golang-example-app/app/http"
	"github.com/aristat/golang-example-app/app/logger"

	"github.com/spf13/cobra"
)

var (
	bind string
	Cmd  = &cobra.Command{
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
			if err := s.ListenAndServe(bind); err != nil {
				log.Error(err.Error())
				return
			}
			log.Info("daemon stopped successfully")
		},
	}
)

func init() {
	Cmd.PersistentFlags().StringVarP(&bind, "bind", "b", ":9096", "bind address")
}
