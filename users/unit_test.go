package users

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"fmt"

	"github.com/aristat/golang-gin-oauth2-example-app/common"
)

var env *common.Env
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
		userModel.setPassword("password123")

		sqlStatement := `
INSERT INTO users (email, encrypted_password)
VALUES ($1, $2)`
		_ = env.DB.QueryRow(sqlStatement, userModel.Email, userModel.EncryptedPassword)
		ret = append(ret, userModel)
	}

	return ret
}

func resetDBWithMock() {
	common.ClearDataTestDB(env.DB)
	userModelMocker(2)
}

func TestRequest(t *testing.T) {
	t.Parallel()

	asserts := assert.New(t)

	r := gin.New()
	r.LoadHTMLGlob("../templates/**/*")

	usersRouters := &UserRouters{Env: env}
	InitRouters(r, usersRouters)

	for _, testData := range requestTests {
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
	common.InitConfig()
	env = common.InitTestEnv()

	exitVal := m.Run()

	common.ClearDataTestDB(env.DB)
	os.Exit(exitVal)
}
