package main

import (
    "time"
    "strings"

    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/limiter"
    "github.com/gofiber/fiber/v2/middleware/csrf"
    "github.com/gofiber/helmet/v2"
)

func (cl *cutlink) setupRoutes() {
    cl.App.Get("/", cl.HomePage)
    cl.App.Get("/auth/signup", cl.SignupPage)
    cl.App.Post("/auth/signup", cl.SignupUser)
    cl.App.Get("/auth/login", cl.LoginPage)
    cl.App.Post("/auth/login", cl.LoginUser)
    cl.App.Post("/auth/logout", cl.LogoutUser)
    cl.App.Post("/auth/delete", cl.DeleteUser)
    cl.App.Get("/r/:hash", cl.Redirector)
    cl.App.Post("/add", cl.AddUrl)
    cl.App.Delete("/delete/:hash", cl.DeleteUrl)
}


func (cl *cutlink) setupMiddlewares() {
    // setup rate limiter middleware
    cl.App.Use(limiter.New(limiter.Config{
        Next: func (c *fiber.Ctx) bool {
            return (strings.HasPrefix(c.Path(), "/r") || strings.HasPrefix(c.Path(), "/delete"))
        },
        Max: 5,
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
        // CookieName: "csrf_",
        CookieSecure: true,
        CookieHTTPOnly: true,
        CookieSameSite: "Strict",
    }))

    cl.App.Use(helmet.New())
}
