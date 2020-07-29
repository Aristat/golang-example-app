package repo_test

import (
	"errors"
	"testing"

	"golang.org/x/crypto/bcrypt"

	"github.com/aristat/golang-example-app/app/db"
	"github.com/aristat/golang-example-app/app/db/domain"

	mocket "github.com/selvatico/go-mocket"

	"github.com/aristat/golang-example-app/app/db/repo"
	"github.com/stretchr/testify/assert"
)

var customInsertError = errors.New("sql: custom insert error")

func TestCreateUser(t *testing.T) {
	mocket.Catcher.Logging = true

	dbManager, _, e := db.BuildTest()
	assert.Nil(t, e, "DB manager error should be nil")

	userRepo, _, e := repo.NewUsersRepo(dbManager.DB)
	assert.Nil(t, e, "Repo error should be nil")

	defaultEmail := "test@gmail.com"
	defaultPassword := "123456789"
	primaryKey := 333

	mockCreateWithError := func() {
		mocket.Catcher.
			Reset().
			NewMock().
			WithQuery("INSERT  INTO \"users\" (\"email\",\"encrypted_password\") VALUES ($1,$2)").
			WithArgs(defaultEmail, defaultPassword).
			WithError(customInsertError)
	}
	mockCreate := func() {
		reply := []map[string]interface{}{{"id": primaryKey}}
		mocket.Catcher.
			Reset().
			NewMock().
			WithQuery("INSERT  INTO \"users\" (\"email\",\"encrypted_password\") VALUES ($1,$2)").
			WithArgs(defaultEmail, defaultPassword).
			WithReply(reply)
	}

	tests := []struct {
		name    string
		mock    func()
		asserts func(user *domain.User, err error)
	}{
		{
			name: "USER NOT CREATED",
			mock: mockCreateWithError,
			asserts: func(user *domain.User, err error) {
				assert.Nil(t, user, "user should be nil")
				assert.NotNil(t, err, "err should not be nil")
				assert.Equal(t, customInsertError.Error(), err.Error())
			},
		},
		{
			name: "USER CREATED",
			mock: mockCreate,
			asserts: func(user *domain.User, err error) {
				assert.Equal(t, primaryKey, user.ID)
				assert.Nil(t, err, "err should be nil")
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mock()

			user, err := userRepo.CreateUser(defaultEmail, defaultPassword)
			test.asserts(user, err)
		})
	}
}

func TestFindByEmail(t *testing.T) {
	mocket.Catcher.Logging = true

	dbManager, _, e := db.BuildTest()
	assert.Nil(t, e, "DB manager error should be nil")

	userRepo, _, e := repo.NewUsersRepo(dbManager.DB)
	assert.Nil(t, e, "Repo error should be nil")

	ePassword, e := bcrypt.GenerateFromPassword([]byte("12345"), 8)
	assert.Nil(t, e, "Password is correct")

	defaultEmail := "test@gmail.com"
	mockSelect := func() {
		reply := []map[string]interface{}{{"id": 1, "email": defaultEmail, "encrypted_password": ePassword}}
		mocket.Catcher.Reset().NewMock().WithQuery(`WHERE (users.email = $1) LIMIT 1`).WithArgs(defaultEmail).WithReply(reply)
	}

	tests := []struct {
		name    string
		mock    func()
		asserts func(user *domain.User, err error)
	}{
		{
			name: "USER IS EMPTY",
			mock: func() {},
			asserts: func(user *domain.User, err error) {
				assert.Equal(t, "", user.Email, "email should be equals")
				assert.NotNil(t, err, "err should not be nil")
				assert.Equal(t, err.Error(), "record not found", "error should be corrects")
			},
		},
		{
			name: "USER EXIST",
			mock: mockSelect,
			asserts: func(user *domain.User, err error) {
				assert.Equal(t, defaultEmail, user.Email, "email should be equals")
				assert.Nil(t, e, "err should be nil")
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mocket.Catcher.Reset()
			test.mock()

			user, err := userRepo.FindByEmail(defaultEmail)
			test.asserts(user, err)
		})
	}
}
