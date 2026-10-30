# Migration from Fiber v2

If you're upgrading from Fiber v2 to v3, here are the key changes affecting `zentests`.

For the zentests basics, see the [README](../README.md).

## Import Path Change

```go
// Fiber v2
import "github.com/gofiber/fiber/v2"

// Fiber v3
import "github.com/gofiber/fiber/v3"
```

## Handler Signature Change

Fiber v3 uses `fiber.Ctx` as an interface (not a pointer):

```go
// Fiber v2
app.Get("/route", func(c *fiber.Ctx) error {
    return c.SendString("hello")
})

// Fiber v3
app.Get("/route", func(c fiber.Ctx) error {
    return c.SendString("hello")
})
```

## Body Parser Change

The `BodyParser` method has been replaced with the new `Bind()` API:

```go
// Fiber v2
app.Post("/users", func(c *fiber.Ctx) error {
    var data map[string]any
    c.BodyParser(&data)
    return c.JSON(data)
})

// Fiber v3
app.Post("/users", func(c fiber.Ctx) error {
    var data map[string]any
    if err := c.Bind().Body(&data); err != nil {
        return err
    }
    return c.JSON(data)
})
```

## New TestConfig Support

Fiber v3's `app.Test()` now accepts `TestConfig` instead of timeout duration:

```go
// Fiber v2
resp, err := app.Test(req, 5*time.Second)

// Fiber v3
resp, err := app.Test(req) // Uses default 1s timeout
resp, err := app.Test(req, fiber.TestConfig{
    Timeout: 5 * time.Second,
})
```

`zentests` provides `*WithConfig` methods to expose this functionality:

```go
// Fiber v2 style (still works)
zt.Get(app, "/route")

// Fiber v3 with custom timeout
zt.GetWithConfig(app, "/slow-route", fiber.TestConfig{
    Timeout: 10 * time.Second,
})
```
