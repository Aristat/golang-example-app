package client

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/aristat/golang-example-app/app/config"

	"github.com/aristat/golang-example-app/app/logger"

	"github.com/aristat/golang-example-app/common"

	"github.com/opentracing-contrib/go-stdlib/nethttp"

	"github.com/go-chi/chi"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
)

// Config
type Config struct {
	Url          string
	Port         string
	AuthUrl      string
	TokenUrl     string
	ClientID     string
	ClientSecret string
}

var token *oauth2.Token

var (
	//bind string
	Cmd = &cobra.Command{
		Use:           "oauth-client",
		Short:         "Test oauth client",
		SilenceUsage:  true,
		SilenceErrors: true,
		Run: func(_ *cobra.Command, _ []string) {
			conf, c, e := config.Build()
			if e != nil {
				panic(e)
			}
			defer c()

			clientConfig := Config{}
			e = conf.UnmarshalKey("services.oauthClient", &clientConfig)
			if e != nil {
				log.Fatal("Config initialize error")
			}

			oauth2Config := oauth2.Config{
				ClientID:     clientConfig.ClientID,
				ClientSecret: clientConfig.ClientSecret,
				Scopes:       []string{"all"},
				RedirectURL:  clientConfig.Url + "/oauth2",
				Endpoint: oauth2.Endpoint{
					AuthURL:  clientConfig.AuthUrl,
					TokenURL: clientConfig.TokenUrl,
				},
			}

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

			tracer := common.GenerateTracerForTestClient("golang-example-app-client", conf)
			client := &http.Client{Transport: &nethttp.Transport{}}

			r := chi.NewRouter()
			r.Get("/login", func(w http.ResponseWriter, r *http.Request) {
				u := oauth2Config.AuthCodeURL("xyz")
				http.Redirect(w, r, u, http.StatusFound)
			})

			r.Get("/oauth2", func(w http.ResponseWriter, r *http.Request) {
				r.ParseForm()
				state := r.Form.Get("state")
				if state != "xyz" {
					http.Error(w, "State invalid", http.StatusBadRequest)
					return
				}
				code := r.Form.Get("code")
				if code == "" {
					http.Error(w, "Code not found", http.StatusBadRequest)
					return
				}
				tokenC, err := oauth2Config.Exchange(context.Background(), code)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				token = tokenC
				e := json.NewEncoder(w)
				e.SetIndent("", "  ")
				e.Encode(*token)

			})

			r.Get("/user", func(w http.ResponseWriter, r *http.Request) {
				req, err := http.NewRequest("GET", "http://localhost:9096/user", nil)
				if err != nil {
					log.Printf("[ERROR] %s", err.Error())
					http.Redirect(w, r, "/login", http.StatusFound)
					return
				}

				req.Header.Set("Authorization", token.Type()+" "+token.AccessToken)

				req = req.WithContext(r.Context())
				req, ht := nethttp.TraceRequest(tracer, req, nethttp.OperationName("HTTP GET"))
				defer ht.Finish()

				res, err := client.Do(req)
				if err != nil {
					log.Printf("[ERROR] %s", err.Error())
					http.Redirect(w, r, "/login", http.StatusFound)
					return
				}
				defer res.Body.Close()

				w.Header().Set("Content-Type", "application/json;charset=UTF-8")
				w.Header().Set("Cache-Control", "no-store")
				w.Header().Set("Pragma", "no-cache")

				d, err := ioutil.ReadAll(res.Body)
				w.WriteHeader(http.StatusOK)
				w.Write(d)
			})

			log.Printf("[INFO] Client is running at %s port.", clientConfig.Port)

			server := &http.Server{
				Addr:    ":" + clientConfig.Port,
				Handler: r,
			}

			err := server.ListenAndServe()
			if err != nil {
				panic(err)
			}
		},
	}
)

func init() {
}
