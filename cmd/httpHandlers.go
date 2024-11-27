package main

import (
	"github.com/thehxdev/cutlink/models"
	"regexp"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	urlMatcher *regexp.Regexp = regexp.MustCompile(
		`^((http|https)://)[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-z]{2,6}\b([-a-zA-Z0-9@:%_\+.~#?&//=]*)$`)

	userIdMatcher *regexp.Regexp = regexp.MustCompile(
		`[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}`)
)

func (cl *cutlink) ErrorPage(c *fiber.Ctx, errMsg string) error {
	return c.Render("error", fiber.Map{
		"title": "Error",
		"msg":   errMsg,
	}, "layouts/main")
}

func (cl *cutlink) HomePage(c *fiber.Ctx) error {
	var urls []*models.Url

	sess, err := cl.Store.Get(c)
	if err != nil {
		cl.ErrorLog.Println(err.Error())
		return err
	}

	id := sess.Get("authenticatedUserID")
	if id == nil {
		id = 0
	}

	if id != 0 {
		urls, err = cl.Conn.GetAllUrls(id.(int))
		if err != nil {
			cl.ErrorLog.Println(err.Error())
			return fiber.ErrInternalServerError
		}
	}

	err = c.Render("index", fiber.Map{
		"title":         "Home",
		"Urls":          urls,
		"authenticated": id,
	}, "layouts/main")

	if err != nil {
		cl.ErrorLog.Println(err.Error())
	}
	sess.Delete("errMsg")
	sess.Save()

	return err
}

func (cl *cutlink) SignupPage(c *fiber.Ctx) error {
	sess, err := cl.Store.Get(c)
	if err != nil {
		cl.ErrorLog.Println(err.Error())
		return fiber.ErrInternalServerError
	}

	if !sess.Fresh() {
		err = sess.Regenerate()
		if err != nil {
			cl.ErrorLog.Println(err.Error())
			return fiber.ErrInternalServerError
		}
	}

	errMsg := sess.Get("errMsg")
	if errMsg != nil {
		errMsg = errMsg.(string)
	}

	err = c.Render("signup", fiber.Map{
		"title":    "Signup",
		"errMsg":   errMsg,
		"disabled": cl.Cfg.Management.NoSignup,
	}, "layouts/main")

	if errMsg != nil {
		sess.Delete("errMsg")
	}

	if err != nil {
		cl.ErrorLog.Println(err.Error())
	}

	return err
}

func (cl *cutlink) SignupUser(c *fiber.Ctx) error {
	if cl.Cfg.Management.NoSignup {
		return c.SendString("Signup is disabled.")
	}

	sess, err := cl.Store.Get(c)
	if err != nil {
		cl.ErrorLog.Println(err.Error())
		return fiber.ErrInternalServerError
	}

	err = sess.Regenerate()
	if err != nil {
		cl.ErrorLog.Println(err.Error())
		return fiber.ErrInternalServerError
	}

	password := c.FormValue("password", "")
	if password == "" || len(password) <= 8 {
		sess.Set("errMsg", "Password must be more than 8 characters.")
		sess.Save()
		return c.Redirect("/auth/signup", fiber.StatusSeeOther)
	}

	userID, err := uuid.NewRandom()
	if err != nil {
		cl.ErrorLog.Println(err.Error())
		return fiber.ErrInternalServerError
	}

	err = cl.Conn.CreateUser(userID.String(), password)
	if err != nil {
		cl.ErrorLog.Println(err.Error())
		return fiber.ErrInternalServerError
	}

	sess.Set("userid", userID.String())
	err = sess.Save()
	if err != nil {
		cl.ErrorLog.Println(err.Error())
		return fiber.ErrInternalServerError
	}

	return c.Redirect("/auth/login", fiber.StatusSeeOther)
}

func (cl *cutlink) LoginPage(c *fiber.Ctx) error {
	sess, err := cl.Store.Get(c)
	if err != nil {
		return fiber.ErrInternalServerError
	}

	errMsg := sess.Get("errMsg")
	if errMsg != nil {
		errMsg = errMsg.(string)
	}

	userID := sess.Get("userid")
	if userID != nil {
		userID = userID.(string)
	}

	err = c.Render("login", fiber.Map{
		"errMsg": errMsg,
		"userid": userID,
		"title":  "Login",
	}, "layouts/main")

	sess.Delete("errMsg")
	sess.Delete("userid")
	sess.Save()
	if err != nil {
		cl.ErrorLog.Println(err.Error())
	}

	return err
}

func (cl *cutlink) LoginUser(c *fiber.Ctx) error {
	sess, err := cl.Store.Get(c)
	if err != nil {
		return fiber.ErrInternalServerError
	}

	password := c.FormValue("password", "")
	userID := strings.TrimSpace(c.FormValue("uuid", ""))

	if !userIdMatcher.Match([]byte(userID)) {
		sess.Set("errMsg", "Invalid UserID or Password.")
		sess.Save()
		return c.Redirect("/auth/login", fiber.StatusSeeOther)
	}

	id, isAdmin, err := cl.Conn.AuthenticatUser(userID, password)
	if err != nil {
		sess.Set("errMsg", "Invalid UserID or Password.")
		sess.Save()
		return c.Redirect("/auth/login", fiber.StatusSeeOther)
	}

	sess.Regenerate()
	sess.Set("authenticatedUserID", id)

	if isAdmin {
		sess.Set("admin", true)
	} else {
		sess.Set("admin", false)
	}

	err = sess.Save()
	if err != nil {
		cl.ErrorLog.Println(err.Error())
		return fiber.ErrInternalServerError
	}

	return c.Redirect("/", fiber.StatusSeeOther)
}

func (cl *cutlink) LogoutUser(c *fiber.Ctx) error {
	sess, err := cl.Store.Get(c)
	if err != nil {
		cl.ErrorLog.Println(err.Error())
		return fiber.ErrInternalServerError
	}

	err = sess.Regenerate()
	if err != nil {
		cl.ErrorLog.Println(err.Error())
		return fiber.ErrInternalServerError
	}

	sess.Destroy()
	return c.Redirect("/", fiber.StatusSeeOther)
}

func (cl *cutlink) DeleteUser(c *fiber.Ctx) error {
	sess, err := cl.Store.Get(c)
	if err != nil {
		cl.ErrorLog.Println(err.Error())
		return fiber.ErrInternalServerError
	}

	err = sess.Regenerate()
	if err != nil {
		cl.ErrorLog.Println(err.Error())
		return fiber.ErrInternalServerError
	}
	id := sess.Get("authenticatedUserID")

	err = cl.Conn.DeleteUser(id.(int))
	if err != nil {
		return fiber.ErrInternalServerError
	}

	if id != nil {
		sess.Destroy()
	}

	return c.Redirect("/", fiber.StatusSeeOther)
}

func (cl *cutlink) Redirector(c *fiber.Ctx) error {
	hash := c.Params("hash")
	if hash == "" {
		return cl.ErrorPage(c, "URL hash is not valid.")
	}
	target, err := cl.Conn.GetUrl(hash)
	if err != nil {
		return fiber.ErrNotFound
	}

	if target.PassHash != "" {
		return c.Render("redirect", fiber.Map{
			"title": "Protected",
			"hash":  target.Hash,
		}, "layouts/main")
	}

	err = cl.Conn.IncrementClicked(hash)
	if err != nil {
		cl.ErrorLog.Println("Incrementing for", target.Hash, "Failed")
	}

	direct := c.QueryBool("direct", false)
	if !direct {
		return c.Render("safe", fiber.Map{
			"title":  "Safe Mode",
			"target": target.Target,
		}, "layouts/main")
	}

	return c.Redirect(target.Target, fiber.StatusSeeOther)
}

func (cl *cutlink) RedirectorPassword(c *fiber.Ctx) error {
	hash := c.Params("hash")
	if hash == "" {
		return fiber.ErrNotFound
	}

	target, err := cl.Conn.GetUrl(hash)
	if err != nil {
		return fiber.ErrNotFound
	}
	password := c.FormValue("password", "")

	err = bcrypt.CompareHashAndPassword([]byte(target.PassHash), []byte(password))
	if err != nil {
		return cl.ErrorPage(c, "Invalid password.")
	}

	err = cl.Conn.IncrementClicked(hash)
	if err != nil {
		cl.ErrorLog.Println("Incrementing for", target.Hash, "Failed")
	}

	direct := c.QueryBool("direct", false)
	if !direct {
		return c.Render("safe", fiber.Map{
			"title":  "Safe Mode",
			"target": target.Target,
		}, "layouts/main")
	}

	return c.Redirect(target.Target, fiber.StatusSeeOther)
}

func (cl *cutlink) AddUrl(c *fiber.Ctx) error {
	sess, err := cl.Store.Get(c)
	if err != nil {
		cl.ErrorLog.Println(err.Error())
		return fiber.ErrInternalServerError
	}
	id := sess.Get("authenticatedUserID")

	if id == nil {
		return c.SendString("Not authenticated request.")
	}

	target := strings.TrimSpace(c.FormValue("target", ""))
	if target == "" || !urlMatcher.Match([]byte(target)) {
		c.Set("HX-Retarget", "#add-url-form")
		c.Set("HX-Reswap", "outerHTML")
		return c.Render("partials/addUrlForm", fiber.Map{
			"errMsg": "Target URL is not valid.",
		})
	} else if len(target) > 1024 {
		c.Set("HX-Refresh", "true")
		c.Set("HX-Reswap", "outerHTML")
		return c.Render("partials/addUrlForm", fiber.Map{
			"errMsg": "Target URL must be less than 1024 characters.",
		})
	}
	password := strings.TrimSpace(c.FormValue("password", ""))

	_, hash, err := cl.Conn.CreateUrl(id.(int), target, password)
	if err != nil {
		cl.ErrorLog.Println(err.Error())
		return fiber.ErrInternalServerError
	}

	newUrl, err := cl.Conn.GetUrl(hash)
	if err != nil {
		cl.ErrorLog.Println(err.Error())
		return fiber.ErrInternalServerError
	}

	return c.Render("partials/urlrow", newUrl)
}

func (cl *cutlink) DeleteUrl(c *fiber.Ctx) error {
	sess, err := cl.Store.Get(c)
	if err != nil {
		cl.ErrorLog.Println(err.Error())
		return fiber.NewError(fiber.StatusInternalServerError, "Internal Server Error")
	}
	id := sess.Get("authenticatedUserID")

	if id == nil {
		return c.SendString("Not authenticated request.")
	}

	hash := c.Params("hash")
	err = cl.Conn.DeleteUrl(id.(int), hash)
	if err != nil {
		cl.ErrorLog.Println(err.Error())
		return fiber.NewError(fiber.StatusInternalServerError, "Internal Server Error")
	}

	return nil
}
