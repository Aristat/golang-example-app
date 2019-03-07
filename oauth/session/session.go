package session

import (
	"context"
	"log"

	"github.com/go-session/session"
)

func Init(manager session.ManagerStore) *session.Manager {
	if _, err := manager.Check(context.Background(), ""); err != nil {
		log.Fatal(err)
	}

	sessionManager := session.NewManager(
		session.SetStore(manager),
	)

	return sessionManager
}
