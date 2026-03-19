package zentests

import (
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/assert"
)

func TestSuiteHooks(t *testing.T) {
	zt := New(t)
	var order []string

	zt.Describe("Test Suite", func(ctx *SuiteContext) {
		ctx.BeforeEach(func(zt *T) { //nolint:revive
			order = append(order, "before")
		}).AfterEach(func(zt *T) { //nolint:revive
			order = append(order, "after")
		})

		ctx.It("test 1", func(zt *T) { //nolint:revive
			order = append(order, "test1")
		})

		ctx.It("test 2", func(zt *T) { //nolint:revive
			order = append(order, "test2")
		})
	})

	assert.Equal(t, []string{"before", "test1", "after", "before", "test2", "after"}, order)
}

func TestSuiteWithHTTP(t *testing.T) {
	zt := New(t)

	zt.Describe("API Suite", func(ctx *SuiteContext) {
		var app *fiber.App

		ctx.BeforeEach(func(zt *T) { //nolint:revive
			app = fiber.New()
			app.Get("/test", func(c fiber.Ctx) error {
				return c.JSON(fiber.Map{"status": "ok"})
			})
		})

		ctx.It("returns ok", func(zt *T) {
			zt.Get(app, "/test").OK().IsJSON().Has("status", "ok")
		})
	})
}

func TestSuiteNested(t *testing.T) {
	zt := New(t)
	var count int

	zt.Describe("Outer", func(ctx *SuiteContext) {
		ctx.It("outer test", func(zt *T) { //nolint:revive
			count++
		})

		// Note: Nested describes would need more implementation
	})

	assert.Equal(t, 1, count)
}

func TestSuiteNoHooks(t *testing.T) {
	zt := New(t)
	var ran bool

	zt.Describe("No Hooks", func(ctx *SuiteContext) {
		ctx.It("still runs", func(zt *T) { //nolint:revive
			ran = true
		})
	})

	assert.True(t, ran)
}
