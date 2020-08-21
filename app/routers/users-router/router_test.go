package users_router_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	users_router "github.com/aristat/golang-example-app/app/routers/users-router"

	"golang.org/x/crypto/bcrypt"

	mocket "github.com/selvatico/go-mocket"
	"github.com/stretchr/testify/assert"
)

func TestGetLogin(t *testing.T) {
	provider, _, _ := users_router.BuildTest()
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/login", nil)
	provider.Router.GetLogin(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestPostLogin(t *testing.T) {
	mocket.Catcher.Logging = true

	tests := []struct {
		name         string
		urlValues    url.Values
		mock         func()
		expectedCode int
	}{
		{
			name:      "successful logged",
			urlValues: url.Values{"email": {"test_email"}, "password": {"test_password"}},
			mock: func() {
				ePassword, e := bcrypt.GenerateFromPassword([]byte("test_password"), 8)
				assert.Nil(t, e, "Password is correct")
				reply := []map[string]interface{}{{"id": 1, "email": "test_email", "encrypted_password": ePassword}}
				mocket.Catcher.Reset().NewMock().WithQuery(`WHERE (users.email = $1) LIMIT 1`).WithArgs("test_email").WithReply(reply)
			},
			expectedCode: http.StatusFound,
		},
		{
			name:         "user does not exist",
			urlValues:    url.Values{"email": {"test_email"}, "password": {"test_password"}},
			mock:         func() {},
			expectedCode: http.StatusNotFound,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mocket.Catcher.Reset()
			test.mock()

			provider, _, e := users_router.BuildTest()
			assert.Nil(t, e, "err should be nil")

			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/login", strings.NewReader(test.urlValues.Encode()))
			req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
			provider.Router.PostLogin(rec, req)

			assert.Equal(t, test.expectedCode, rec.Code)
		})
	}
}
