package main

import (
    "fmt"
    "regexp"
    "cutlink/models"

    "github.com/gofiber/fiber/v2"
    "github.com/google/uuid"
)


var (

urlMatcher *regexp.Regexp = regexp.MustCompile(
    `^((http|https)://)[-a-zA-Z0-9@:%._\+~#=]{2,256}\.[a-z]{2,6}\b([-a-zA-Z0-9@:%_\+.~#?&//=]*)$`)

userIdMatcher *regexp.Regexp = regexp.MustCompile(
        `[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}`)
)


func (cl *cutlink) HomePage(c *fiber.Ctx) error {
    var urls []*models.Url
    sess, err := cl.Store.Get(c)
    if err != nil {
        cl.ErrorLog.Println(err.Error())
        return fiber.ErrInternalServerError
    }
    id := sess.Get("authenticatedUserID")

    if id != nil {
        urls, err = cl.Conn.GetAllUrls(id.(int))
        if err != nil {
            cl.ErrorLog.Println(err.Error())
            return fiber.NewError(fiber.StatusInternalServerError, "Internal Server Error")
        }
    }

    err = c.Render("index", fiber.Map{
        "title": "Home",
        "Urls": urls,
        "authenticated": id,
    }, "layouts/main")

    if err != nil {
        cl.ErrorLog.Println(err.Error())
        return err
    }

    return nil
}


func (cl *cutlink) SignupPage(c *fiber.Ctx) error {
    err := c.Render("signup", fiber.Map{
        "title": "Signup",
        "disabled": cl.DisableSignup,
    }, "layouts/main")

    if err != nil {
        cl.ErrorLog.Println(err.Error())
        return err
    }

    return nil
}


func (cl *cutlink) SignupUser(c *fiber.Ctx) error {
    if cl.DisableSignup {
        return c.SendString("Signup is disabled.")
    }

    password := c.FormValue("password", "")
    if password == "" || len(password) <= 8 {
        cl.ErrorLog.Println("provided password is not valid")
        retval := `<div class="container alert alert-danger" role="alert"><h4>Password Is NOT Valid</h4></div>`
        return c.SendString(retval)
    }

    userID, err := uuid.NewRandom()
    if err != nil {
        cl.ErrorLog.Println(err.Error())
        return fiber.NewError(fiber.StatusInternalServerError, "Internal Server Error")
    }

    err = cl.Conn.CreateUser(userID.String(), password)
    if err != nil {
        cl.ErrorLog.Println(err.Error())
        return fiber.NewError(fiber.StatusInternalServerError, "Internal Server Error")
    }

    cl.InfoLog.Println("new user created with uuid", userID.String())

    retval := fmt.Sprintf(
        `<div class="container alert alert-success" role="alert"><h2>UUID</h2><code style="font-size: 20px">%s</code></div>`,
        userID.String())
    err = c.SendString(retval)

    if err != nil {
        cl.ErrorLog.Println(err.Error())
        return err
    }

    return nil
}


func (cl *cutlink) LoginPage(c *fiber.Ctx) error {
    err := c.Render("login", fiber.Map{
        "title": "Login",
    }, "layouts/main")

    if err != nil {
        cl.ErrorLog.Println(err.Error())
        return err
    }

    return nil
}


func (cl *cutlink) LoginUser(c *fiber.Ctx) error {
    password := c.FormValue("password", "")
    userID   := c.FormValue("uuid", "")

    if !userIdMatcher.Match([]byte(userID)) {
        return fiber.ErrInternalServerError
    }

    id, err := cl.Conn.AuthenticatUser(userID, password)
    if err != nil {
        cl.ErrorLog.Println(err.Error())
        return fiber.ErrInternalServerError
    }

    sess, err := cl.Store.Get(c)
    if err != nil {
        cl.ErrorLog.Println(err.Error())
        return fiber.ErrInternalServerError
    }

    sess.Set("authenticatedUserID", id)
    sess.Save()
    return c.Redirect("/", fiber.StatusSeeOther)
}


func (cl *cutlink) LogoutUser(c *fiber.Ctx) error {
    sess, err := cl.Store.Get(c)
    if err != nil {
        cl.ErrorLog.Println(err.Error())
        return fiber.ErrInternalServerError
    }
    id := sess.Get("authenticatedUserID")

    if id != nil {
        sess.Regenerate()
        sess.Destroy()
        sess.Save()
    }

    return c.Redirect("/", fiber.StatusSeeOther)
}


func (cl *cutlink) DeleteUser(c *fiber.Ctx) error {
    sess, err := cl.Store.Get(c)
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
        sess.Regenerate()
        sess.Destroy()
        sess.Save()
    }

    return c.Redirect("/", fiber.StatusSeeOther)
}


func (cl *cutlink) Redirector(c *fiber.Ctx) error {
    hash := c.Params("hash")
    if hash == "" {
        return fiber.NewError(fiber.StatusInternalServerError, "Internal Server Error")
    }

    target, err := cl.Conn.GetUrl(hash)
    if err != nil {
        return fiber.NewError(fiber.StatusNotFound, "Not Found")
    }

    err = cl.Conn.IncrementClicked(hash)
    if err != nil {
        cl.ErrorLog.Println("Incrementing for", target.Hash, "Failed")
    }

    return c.Redirect(target.Target, fiber.StatusSeeOther)
}


func (cl *cutlink) AddUrl(c *fiber.Ctx) error {
    sess, err := cl.Store.Get(c)
    if err != nil {
        cl.ErrorLog.Println(err.Error())
        return fiber.NewError(fiber.StatusInternalServerError, "Internal Server Error")
    }
    id := sess.Get("authenticatedUserID")

    target := c.FormValue("target", "")
    if target == "" {
        return fiber.NewError(fiber.StatusInternalServerError, "Internal Server Error")
    }

    if !urlMatcher.Match([]byte(target)) {
        cl.ErrorLog.Println("target url does not match the pattern")
        return c.Redirect("/", fiber.StatusSeeOther)
    }

    _, _, err = cl.Conn.CreateUrl(id.(int), target)
    if (err != nil) {
        cl.ErrorLog.Println(err.Error())
        return fiber.NewError(fiber.StatusInternalServerError, "Internal Server Error")
    }

    cl.InfoLog.Println("new url added:", target)
    return c.Redirect("/", fiber.StatusSeeOther)
}


func (cl *cutlink) DeleteUrl(c *fiber.Ctx) error {
    sess, err := cl.Store.Get(c)
    if err != nil {
        cl.ErrorLog.Println(err.Error())
        return fiber.NewError(fiber.StatusInternalServerError, "Internal Server Error")
    }
    id := sess.Get("authenticatedUserID")

    hash := c.Params("hash")
    err = cl.Conn.DeleteUrl(id.(int), hash)
    if err != nil {
        cl.ErrorLog.Println(err.Error())
        return fiber.NewError(fiber.StatusInternalServerError, "Internal Server Error")
    }

    cl.InfoLog.Printf("url with hash %s has been deleted", hash)
    return nil
}
