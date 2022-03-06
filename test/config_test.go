package test

import (
	"testing"

	"github.com/mrlightwood/golang-products-api/config"

	"github.com/stretchr/testify/assert"
)

func TestConfig_NewConfig(t *testing.T) {
	c, err := config.NewConfig("../config/config_test.yaml")
	assert.Nil(t, err)
	assert.NotNil(t, c)
}
