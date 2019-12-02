package resolver_test

import (
	"context"
	"testing"

	"github.com/spf13/cast"

	"golang.org/x/crypto/bcrypt"

	"github.com/stretchr/testify/assert"

	"github.com/aristat/golang-example-app/app/resolver"
	graphql1 "github.com/aristat/golang-example-app/generated/graphql"
	mocket "github.com/selvatico/go-mocket"
)

type UserParams struct {
	testName string
	id       int
	email    string
}

func TestOne(t *testing.T) {
	ctx := context.Background()

	cfg, _, err := resolver.BuildTest()
	if err != nil {
		assert.Failf(t, "resolver instance failed, err: %v", err.Error())
		return
	}

	ePassword, e := bcrypt.GenerateFromPassword([]byte("12345"), 8)
	assert.Nil(t, e, "Password is correct")

	obj := graphql1.UsersQuery{}
	defaultEmail := "test@gmail.com"
	id := 1
	mockDefaultResult := func() {
		reply := []map[string]interface{}{{"id": id, "email": defaultEmail, "encrypted_password": ePassword}}
		mocket.Catcher.Reset().NewMock().WithQuery(`WHERE (users.email = $1) LIMIT 1`).WithArgs(defaultEmail).WithReply(reply)
	}

	tests := []struct {
		userParams UserParams
		mock       func()
		asserts    func(userParams UserParams, out *graphql1.UsersOneOut, err error)
	}{
		{
			userParams: UserParams{id: 2, email: "test1@gmail.com", testName: "USER NOT EXIST"},
			mock:       mockDefaultResult,
			asserts: func(userParams UserParams, out *graphql1.UsersOneOut, err error) {
				assert.NotNil(t, err, "err should not be nil")
				assert.Nil(t, out, "out should be nil")
			},
		},
		{
			userParams: UserParams{id: id, email: defaultEmail, testName: "USER EXIST"},
			mock:       mockDefaultResult,
			asserts: func(userParams UserParams, out *graphql1.UsersOneOut, err error) {
				assert.Nil(t, err, "err should be nil")
				assert.Equal(t, cast.ToString(userParams.id), out.ID, "id should be equals")
				assert.Equal(t, userParams.email, out.Email, "email should be equals")
			},
		},
	}

	for _, test := range tests {
		t.Run(test.userParams.testName, func(t *testing.T) {
			mocket.Catcher.Reset()
			test.mock()

			out, err := cfg.Resolvers.UsersQuery().One(ctx, &obj, test.userParams.email)
			test.asserts(test.userParams, out, err)
		})
	}
}
