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

func TestSuiteNested_HooksRunOuterBeforeInnerAndInnerAfterOuter(t *testing.T) {
	zt := New(t)
	var order []string

	zt.Describe("Outer", func(outer *SuiteContext) {
		outer.BeforeEach(func(zt *T) { order = append(order, "outer-before") }).
			AfterEach(func(zt *T) { order = append(order, "outer-after") })

		outer.Describe("Inner", func(inner *SuiteContext) {
			inner.BeforeEach(func(zt *T) { order = append(order, "inner-before") }).
				AfterEach(func(zt *T) { order = append(order, "inner-after") })

			inner.It("test", func(zt *T) { order = append(order, "test") })
		})
	})

	assert.Equal(t, []string{"outer-before", "inner-before", "test", "inner-after", "outer-after"}, order)
}

func TestSuiteNested_InnerBlockWithoutHooks_InheritsOuterHooks(t *testing.T) {
	zt := New(t)
	var order []string

	zt.Describe("Outer", func(outer *SuiteContext) {
		outer.BeforeEach(func(zt *T) { order = append(order, "before") })

		outer.Describe("Inner", func(inner *SuiteContext) {
			inner.It("test", func(zt *T) { order = append(order, "test") })
		})
	})

	assert.Equal(t, []string{"before", "test"}, order)
}

func TestSuiteHooks_DeclaredAfterIt_NotApplied(t *testing.T) {
	zt := New(t)
	var order []string

	zt.Describe("Late hooks", func(ctx *SuiteContext) {
		ctx.It("test", func(zt *T) { order = append(order, "test") })
		ctx.BeforeEach(func(zt *T) { order = append(order, "before") })
	})

	assert.Equal(t, []string{"test"}, order)
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

func TestSuiteAfterEachRunsWhenBeforeEachSkips(t *testing.T) {
	var afterEachRan bool

	New(t).Describe("Setup skip", func(ctx *SuiteContext) {
		ctx.BeforeEach(func(zt *T) {
			zt.SkipNow()
		}).AfterEach(func(zt *T) {
			afterEachRan = true
		})

		ctx.It("skips during setup", func(zt *T) {
			t.Fatal("test body must not run after setup skips")
		})
	})

	assert.True(t, afterEachRan, "AfterEach should run when BeforeEach skips")
}
