package oauth_test

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-session/session"

	"github.com/aristat/golang-example-app/app/logger"
	"github.com/aristat/golang-example-app/app/oauth"
	"github.com/gavv/httpexpect"
	"github.com/stretchr/testify/assert"
)

var (
	testSrv      *httptest.Server
	clientSrv    *httptest.Server
	clientID     = "123456"
	clientSecret = "12345678"
)

var requestTests = []struct {
	init             func(s *oauth.Router)
	authorizeRequest func(e *httpexpect.Expect)
	tokenRequest     func(e *httpexpect.Expect, code string) *httpexpect.Object
}{
	{
		func(router *oauth.Router) {},
		func(e *httpexpect.Expect) {
			e.GET("/authorize").
				WithQuery("response_type", "code").
				WithQuery("client_id", clientID).
				WithQuery("scope", "all").
				WithQuery("state", "123").
				WithQuery("redirect_uri", clientSrv.URL+"/oauth2").
				Expect().Status(http.StatusOK)
		},
		func(e *httpexpect.Expect, code string) *httpexpect.Object {
			return e.POST("/token").
				WithFormField("redirect_uri", clientSrv.URL+"/oauth2").
				WithFormField("code", code).
				WithFormField("grant_type", "authorization_code").
				WithFormField("client_id", clientID).
				WithBasicAuth(clientID, clientSecret).
				Expect().
				Status(http.StatusOK).
				JSON().Object()
		},
	},
	{
		func(router *oauth.Router) {},
		func(e *httpexpect.Expect) {
			e.GET("/authorize").
				WithQuery("response_type", "code").
				WithQuery("client_id", clientID).
				WithQuery("scope", "all").
				WithQuery("state", "123").
				WithQuery("redirect_uri", clientSrv.URL+"/oauth2").
				Expect().Status(http.StatusOK)
		},
		func(e *httpexpect.Expect, code string) *httpexpect.Object {
			e.POST("/token").
				WithFormField("code", code).
				WithFormField("grant_type", "authorization_code").
				WithFormField("client_id", clientID).
				WithBasicAuth(clientID, clientSecret).
				Expect().
				Status(http.StatusBadRequest).
				JSON().Object()

			return nil
		},
	},
	{
		func(router *oauth.Router) {},
		func(e *httpexpect.Expect) {
			e.GET("/authorize").Expect().Status(http.StatusBadRequest)
		},
		func(e *httpexpect.Expect, code string) *httpexpect.Object { return nil },
	},
	{
		func(router *oauth.Router) {
			memoryStore := &oauth.MemoryStore{}

			memoryStore.CheckFn = func(ctx context.Context, sid string) (bool, error) {
				return false, nil
			}

			memoryStore.CreateFn = func(ctx context.Context, sid string, expired int64) (session.Store, error) {
				return nil, errors.New("don't create store")
			}

			router.SessionManager = session.NewManager(
				session.SetStore(memoryStore),
			)
		},
		func(e *httpexpect.Expect) {
			e.GET("/authorize").
				WithQuery("response_type", "code").
				WithQuery("client_id", clientID).
				WithQuery("scope", "all").
				WithQuery("state", "123").
				WithQuery("redirect_uri", clientSrv.URL+"/oauth2").
				Expect().Status(http.StatusInternalServerError)
		},
		func(e *httpexpect.Expect, code string) *httpexpect.Object { return nil },
	},
	{
		func(router *oauth.Router) {
			store := &oauth.KeyValueStore{}
			store.GetFn = func(key string) (interface{}, bool) {
				return nil, false
			}
			store.SessionIDFn = func() string {
				return "12345"
			}
			store.DeleteFn = func(key string) interface{} { return nil }
			store.SaveFn = func() error { return nil }
			store.SetFn = func(key string, value interface{}) {
				return
			}

			memoryStore := &oauth.MemoryStore{}
			memoryStore.CheckFn = func(ctx context.Context, sid string) (bool, error) {
				return true, nil
			}
			memoryStore.CreateFn = func(ctx context.Context, sid string, expired int64) (session.Store, error) {
				return store, nil
			}
			memoryStore.UpdateFn = func(ctx context.Context, sid string, expired int64) (session.Store, error) {

				return store, nil
			}

			sessionManager := session.NewManager(
				session.SetStore(memoryStore),
			)
			router.SessionManager = sessionManager
		},
		func(e *httpexpect.Expect) {
			e.GET("/authorize").
				WithQuery("response_type", "code").
				WithQuery("client_id", clientID).
				WithQuery("scope", "all").
				WithQuery("state", "123").
				WithQuery("redirect_uri", clientSrv.URL+"/oauth2").
				Expect().Status(http.StatusOK)
		},
		func(e *httpexpect.Expect, code string) *httpexpect.Object {
			return e.POST("/token").
				WithFormField("redirect_uri", clientSrv.URL+"/oauth2").
				WithFormField("code", code).
				WithFormField("grant_type", "authorization_code").
				WithFormField("client_id", clientID).
				WithBasicAuth(clientID, clientSecret).
				Expect().
				Status(http.StatusOK).
				JSON().Object()
		},
	},
}

func TestNew(t *testing.T) {
	manager, _, e := oauth.BuildTest()

	assert.Nil(t, e, "BuildTest is correct")
	assert.NotNil(t, manager, "manager not nil")

	testSrv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/authorize":
			manager.Router.Authorize(w, r)
		case "/token":
			manager.Router.Token(w, r)
		}
	}))
	defer testSrv.Close()

	httpExpect := httpexpect.New(t, testSrv.URL)

	mockLogger := manager.Logger.(*logger.Mock)
	go func() {
		val := <-mockLogger.Catch()
		t.Logf("Log trace, level: %s, format: %#v\n", val.Level, val.Format)
	}()

	for _, testData := range requestTests {
		testData.init(manager.Router)
		clientSrv = httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.URL.Path {
			case "/oauth2":
				err := r.ParseForm()
				if err != nil {
					t.Error("invalid parse form:", err)
					return
				}

				code, state := r.Form.Get("code"), r.Form.Get("state")
				if state != "123" {
					t.Error("unrecognized state:", state)
					return
				}

				resObj := testData.tokenRequest(httpExpect, code)

				if resObj != nil {
					t.Logf("oauth2 response %v\n", resObj.Raw())

					validationAccessToken(t, resObj.Value("access_token").String().Raw(), manager)
				}
			}
		}))

		l, err := net.Listen("tcp", "127.0.0.1:8090")
		if err != nil {
			log.Fatal(err)
		}
		clientSrv.Listener = l
		clientSrv.Start()

		manager.Router.Server.SetUserAuthorizationHandler(func(w http.ResponseWriter, r *http.Request) (userID string, err error) {
			return "000000", nil
		})
		testData.authorizeRequest(httpExpect)

		clientSrv.Close()
	}
}

func validationAccessToken(t *testing.T, accessToken string, manager *oauth.Manager) {
	req := httptest.NewRequest("GET", "http://example.com", nil)

	req.Header.Set("Authorization", "Bearer "+accessToken)

	ti, err := manager.Router.Server.ValidationBearerToken(req)
	if err != nil {
		t.Error(err.Error())
		return
	}

	if ti.GetClientID() != clientID {
		t.Error("invalid access token")
	}
}
