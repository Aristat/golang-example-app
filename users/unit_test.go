package users

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/go-session/session"
	"golang.org/x/crypto/bcrypt"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"fmt"

	"github.com/aristat/golang-gin-oauth2-example-app/common"
)

var (
	databaseUrl = "postgresql://localhost:5432/oauth2_test?sslmode=disable"
	db          = common.InitDB(databaseUrl)
)

var requestTests = []struct {
	init            func(*http.Request)
	url             string
	method          string
	requestBodyData string
	expectedCode    int
	msg             string
}{
	{
		func(req *http.Request) {
			resetDBWithMock()
		},
		"/login",
		"GET",
		"",
		http.StatusOK,
		"",
	},
	{
		func(req *http.Request) {
			resetDBWithMock()
		},
		"/login",
		"POST",
		`email=user1@linkedin.com&password=password123`,
		http.StatusFound,
		"",
	},
	{
		func(req *http.Request) {
			resetDBWithMock()
		},
		"/login",
		"POST",
		`email=user3@linkedin.com&password=password123`,
		http.StatusNotFound,
		"",
	},
	{
		func(req *http.Request) {
		},
		"/auth",
		"GET",
		"",
		http.StatusFound,
		"",
	},
}

func userModelMocker(n int) []UserModel {
	var ret []UserModel

	for i := 1; i <= n; i++ {
		userModel := UserModel{
			Email: fmt.Sprintf("user%v@linkedin.com", i),
		}
		userModel.setPassword("password123", bcrypt.MinCost)

		sqlStatement := `
INSERT INTO users (email, encrypted_password)
VALUES ($1, $2)`
		_ = db.QueryRow(sqlStatement, userModel.Email, userModel.EncryptedPassword)
		ret = append(ret, userModel)
	}

	return ret
}

func resetDBWithMock() {
	common.ClearDataTestDB(db)
	userModelMocker(2)
}

func TestRequest(t *testing.T) {
	t.Parallel()

	asserts := assert.New(t)

	r := gin.New()
	r.LoadHTMLGlob("../resources/templates/**/*")

	service := &Service{DB: db}
	Run(r, service)

	for _, testData := range requestTests {
		service.SessionManager = session.NewManager()

		requestBodyData := testData.requestBodyData
		req, err := http.NewRequest(testData.method, testData.url, bytes.NewBufferString(requestBodyData))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		asserts.NoError(err)
		testData.init(req)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		asserts.Equal(testData.expectedCode, w.Code)
	}
}

func TestMain(m *testing.M) {
	exitVal := m.Run()

	common.ClearDataTestDB(db)
	os.Exit(exitVal)
}
