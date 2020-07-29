package entrypoint_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/aristat/golang-example-app/app/entrypoint"

	"github.com/stretchr/testify/assert"
)

func TestInitialize(t *testing.T) {
	wd, _ := filepath.Abs(os.Getenv("APP_WD"))
	entryPoint, err := entrypoint.Initialize(wd, nil)

	assert.Nil(t, err, "err should be nil")
	assert.NotNil(t, entryPoint, "entryPoint should not be nil")
}
