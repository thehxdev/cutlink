package main

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// TODO: Make arrors handling better
func (cl *cutlink) AdminHome(c *fiber.Ctx) error {
	sess, err := cl.Store.Get(c)
	if err != nil {
		cl.ErrorLog.Println(err.Error())
		return fiber.ErrInternalServerError
	}

	// sess.Regenerate()
	id := sess.Get("authenticatedUserID")
	if id == nil {
		return fiber.NewError(fiber.StatusInternalServerError, "not authenticated")
	}

	var isAdmin bool
	a := sess.Get("admin")
	if a == nil {
		return fiber.ErrInternalServerError
	}

	isAdmin = a.(bool)
	if !isAdmin {
		return fiber.NewError(fiber.StatusInternalServerError, "not authenticated as admin")
	}

	return c.Render("admin", fiber.Map{
		"title":         "Admin",
		"authenticated": id,
		"route":         cl.Cfg.Admin.Route,
		"noSignup":      cl.Cfg.Management.NoSignup,
		"rateLimit":     cl.Cfg.Management.RateLimitMax,
	}, "layouts/main")
}

func (cl *cutlink) AdminToggleSignup(c *fiber.Ctx) error {
	sess, err := cl.Store.Get(c)
	if err != nil {
		cl.ErrorLog.Println(err.Error())
		return fiber.ErrInternalServerError
	}

	var isAdmin bool
	a := sess.Get("admin")
	if a == nil {
		return fiber.ErrInternalServerError
	}

	isAdmin = a.(bool)
	if !isAdmin {
		return fiber.NewError(fiber.StatusInternalServerError, "not authenticated as admin")
	}

	cl.Cfg.Management.NoSignup = !cl.Cfg.Management.NoSignup
	return c.Redirect(cl.Cfg.Admin.Route, fiber.StatusSeeOther)
}

func (cl *cutlink) AdminSetRateLimitMax(c *fiber.Ctx) error {
	sess, err := cl.Store.Get(c)
	if err != nil {
		cl.ErrorLog.Println(err.Error())
		return fiber.ErrInternalServerError
	}

	var isAdmin bool
	a := sess.Get("admin")
	if a == nil {
		return fiber.ErrInternalServerError
	}

	isAdmin = a.(bool)
	if !isAdmin {
		return fiber.NewError(fiber.StatusInternalServerError, "not authenticated as admin")
	}

	tmp := c.FormValue("ratelimit", "20")
	newRateLimitMax, err := strconv.Atoi(tmp)
	if err != nil {
		return fiber.ErrInternalServerError
	}

	cl.Cfg.Management.RateLimitMax = newRateLimitMax
	return c.Redirect(cl.Cfg.Admin.Route, fiber.StatusSeeOther)
}
