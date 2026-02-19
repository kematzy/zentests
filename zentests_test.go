package zentests

import (
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	zt := New(t)
	assert.NotNil(t, zt)
	assert.NotNil(t, zt.T)
}

func TestUse(t *testing.T) {
	app := fiber.New()

	zt := New(t).Use(app)
	assert.NotNil(t, zt)
	assert.NotNil(t, zt.T)
	assert.NotNil(t, zt.app)
	assert.Equal(t, zt.app, app)
}
