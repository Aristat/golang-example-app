package client

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/opentracing-contrib/go-stdlib/nethttp"

	"github.com/go-chi/chi"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"

	"github.com/uber/jaeger-client-go"
	jaegerConfig "github.com/uber/jaeger-client-go/config"
)

var (
	config = oauth2.Config{
		ClientID:     "123456",
		ClientSecret: "12345678",
		Scopes:       []string{"all"},
		RedirectURL:  "http://localhost:9094/oauth2",
		Endpoint: oauth2.Endpoint{
			AuthURL:  "http://localhost:9096/authorize",
			TokenURL: "http://localhost:9096/token",
		},
	}
)

var token *oauth2.Token

var (
	//bind string
	Cmd = &cobra.Command{
		Use:           "client",
		Short:         "Test client",
		SilenceUsage:  true,
		SilenceErrors: true,
		Run: func(_ *cobra.Command, _ []string) {
			log.SetFlags(log.Lshortfile | log.LstdFlags)

			jaegerCfg := jaegerConfig.Configuration{
				ServiceName: "golang-example-app-client",
				Sampler: &jaegerConfig.SamplerConfig{
					Type:  "const",
					Param: 1,
				},
				Reporter: &jaegerConfig.ReporterConfig{
					LogSpans: false,
				},
			}

			tracer, _, e := jaegerCfg.NewTracer(jaegerConfig.Logger(jaeger.StdLogger))
			if e != nil {
				log.Fatal("Jaeger initialize error")
			}

			client := &http.Client{Transport: &nethttp.Transport{}}

			r := chi.NewRouter()
			r.Get("/login", func(w http.ResponseWriter, r *http.Request) {
				u := config.AuthCodeURL("xyz")
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
				tokenC, err := config.Exchange(context.Background(), code)
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

			log.Println("[INFO] Client is running at 9094 port.")

			server := &http.Server{
				Addr:    ":9094",
				Handler: r,
			}
			log.Fatal(server.ListenAndServe())
		},
	}
)

func init() {
}
