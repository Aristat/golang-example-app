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
		name    string
		id      int
		email   string
		mock    func()
		asserts func(err error)
	}{
		{
			name:  "USER EXIST",
			id:    id,
			email: defaultEmail,
			mock:  mockDefaultResult,
			asserts: func(err error) {
				assert.Nil(t, err, "err should be nil")
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mocket.Catcher.Reset()
			test.mock()

			out, _ := cfg.Resolvers.UsersQuery().One(ctx, &obj, "test@gmail.com")
			t.Log(out.ID, out.Email)

			assert.Equal(t, cast.ToString(test.id), out.ID, "id should be equals")
			assert.Equal(t, test.email, out.Email, "email should be equals")
			test.asserts(err)
		})
	}
}
