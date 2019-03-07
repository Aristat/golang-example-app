package cmd

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"golang.org/x/oauth2"
)

type ClientCommand struct {
}

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

func (cmd *ClientCommand) Execute(args []string) error {
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		u := config.AuthCodeURL("xyz")
		http.Redirect(w, r, u, http.StatusFound)
	})

	http.HandleFunc("/oauth2", func(w http.ResponseWriter, r *http.Request) {
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

	http.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
		res, err := config.Client(context.Background(), token).Get("http://localhost:9096/user")
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
	log.Fatal(http.ListenAndServe(":9094", nil))

	return nil
}
