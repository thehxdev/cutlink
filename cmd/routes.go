package main

import (
    "time"
    "strings"

    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/limiter"
    "github.com/gofiber/fiber/v2/middleware/csrf"
    // "github.com/gofiber/helmet/v2"
)

func (cl *cutlink) setupRoutes() {
    cl.App.Get("/", cl.HomePage)
    cl.App.Get("/r/:hash", cl.Redirector)
    cl.App.Post("/r/:hash", cl.RedirectorPassword)
    cl.App.Post("/add", cl.AddUrl)
    cl.App.Delete("/delete/:hash", cl.DeleteUrl)

    authRoute := cl.App.Group("/auth")
    authRoute.Get("/signup", cl.SignupPage)
    authRoute.Post("/signup", cl.SignupUser)
    authRoute.Get("/login", cl.LoginPage)
    authRoute.Post("/login", cl.LoginUser)
    authRoute.Post("/logout", cl.LogoutUser)
    authRoute.Post("/delete", cl.DeleteUser)

    // adminRoute := cl.App.Group("/admin")
}


func (cl *cutlink) setupMiddlewares() {
    // setup rate limiter middleware
    cl.App.Use(limiter.New(limiter.Config{
        Next: func (c *fiber.Ctx) bool {
            return (c.IP() == "127.0.0.1" && (strings.HasPrefix(c.Path(), "/r") || strings.HasPrefix(c.Path(), "/delete")))
        },
        Max: cl.Cfg.Management.RateLimitMax,
        Expiration: 30 * time.Second,
        LimitReached: func (c *fiber.Ctx) error {
            return c.Render("rateLimit", fiber.Map{
                "title": "Rate Limit",
            }, "layouts/main")
        },
    }))


    // setup CSRF protection
    cl.App.Use(csrf.New(csrf.Config{
        KeyLookup: "cookie:csrf_",
        CookieSecure: true,
        CookieHTTPOnly: true,
        CookieSameSite: "Strict",
    }))

    // cl.App.Use(helmet.New())
}
