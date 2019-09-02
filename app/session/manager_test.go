package session_test

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/aristat/golang-example-app/app/session"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	manager, _, e := session.BuildTest()
	assert.Nil(t, e, "BuildTest error should be nil")
	assert.NotNil(t, manager, "Manager should not be nil")

	w := httptest.NewRecorder()

	buffer := new(bytes.Buffer)
	params := url.Values{}
	buffer.WriteString(params.Encode())
	req, err := http.NewRequest("POST", "http://example.com", buffer)

	if err != nil {
		assert.Failf(t, "build request instance failed, err: %v", err.Error())
		return
	}

	store, err := manager.Start(context.Background(), w, req)
	if err != nil {
		assert.Failf(t, "build store instance failed, err: %v", err.Error())
		return
	}

	store.Set("UserID", "1")
	err = store.Save()
	if err != nil {
		assert.Failf(t, "save key failed, err: %v", err.Error())
		return
	}
}
