package daemon

import (
	"fmt"

	"github.com/aristat/golang-gin-oauth2-example-app/app/http"

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
			s, c, e = http.Build()
			if e != nil {
				fmt.Println(e.Error())
				return
			}
			defer c()
			if err := s.ListenAndServe(bind); err != nil {
				fmt.Println(err.Error())
				return
			}
			fmt.Println("daemon stopped successfully")
		},
	}
)

func init() {
	Cmd.PersistentFlags().StringVarP(&bind, "bind", "b", ":9096", "bind address")
}
