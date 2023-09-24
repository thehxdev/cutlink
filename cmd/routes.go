package main

import (
    // "strings"
    "time"

    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/limiter"
    // "github.com/gofiber/fiber/v2/middleware/csrf"
)

func (cl *cutlink) setupRoutes() {
    cl.App.Static("/static", "./ui/static", fiber.Static{
        Browse: false,
        CacheDuration: 10 * time.Second,
    })
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
            return !(c.Path() == "/auth/signup")
            // return (!strings.HasPrefix(c.Path(), "/auth/signup"))
        },
        Max: 20,
        Expiration: 60 * time.Second,
        LimitReached: func (c *fiber.Ctx) error {
            return c.SendString("Rate Limit Reached. Wait for 60 seconds.")
        },
    }))
    // cl.App.Use(csrf.New())
}
