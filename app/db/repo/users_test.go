package repo

import (
	"testing"

	"golang.org/x/crypto/bcrypt"

	"github.com/aristat/golang-oauth2-example-app/app/db"
	mocket "github.com/selvatico/go-mocket"
	"github.com/stretchr/testify/assert"
)

func TestFindByEmail(t *testing.T) {
	mocket.Catcher.Logging = true

	dbManager, _, e := db.BuildTest()
	assert.Nil(t, e, "DB manager error should be nil")

	repo, _, e := NewUsersRepo(dbManager.DB)
	assert.Nil(t, e, "Repo error should be nil")

	ePassword, e := bcrypt.GenerateFromPassword([]byte("12345"), 8)
	assert.Nil(t, e, "BuildTest is correct")

	defaultEmail := "test@gmail.com"
	mockDefaultResult := func() {
		reply := []map[string]interface{}{{"id": 1, "email": defaultEmail, "encrypted_password": ePassword}}
		mocket.Catcher.Reset().NewMock().WithQuery(`WHERE (users.email = $1) LIMIT 1`).WithArgs(defaultEmail).WithReply(reply)
	}

	tests := []struct {
		name    string
		email   string
		mock    func()
		asserts func(err error)
	}{
		{
			name:  "USER IS EMPTY",
			email: "",
			mock:  func() {},
			asserts: func(err error) {
				assert.NotNil(t, err, "err should not be nil")
				assert.Equal(t, err.Error(), "record not found", "error should be corrects")
			},
		},
		{
			name:  "USER EXIST",
			email: defaultEmail,
			mock:  mockDefaultResult,
			asserts: func(err error) {
				assert.Nil(t, e, "err should be nil")
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mocket.Catcher.Reset()
			test.mock()

			user, err := repo.FindByEmail(defaultEmail)
			assert.Equal(t, test.email, user.Email, "id should be equals")
			test.asserts(err)
		})
	}
}
