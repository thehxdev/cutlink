package main

import (
    "fmt"
    "regexp"
    "cutlink/models"

    "github.com/google/uuid"
    "golang.org/x/crypto/bcrypt"
    "github.com/gofiber/fiber/v2"
)


var (
    urlMatcher *regexp.Regexp = regexp.MustCompile(
        `^((http|https)://)[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-z]{2,6}\b([-a-zA-Z0-9@:%_\+.~#?&//=]*)$`)

    userIdMatcher *regexp.Regexp = regexp.MustCompile(
        `[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}`)
)


func (cl *cutlink) getUserID(c *fiber.Ctx) (int, error) {
    sess, err := cl.Store.Get(c)
    if err != nil {
        return 0, err
    }

    id := sess.Get("authenticatedUserID")
    if id == nil {
        return 0, nil
    }

    return id.(int), nil
}


func (cl *cutlink) ErrorPage(c *fiber.Ctx, errMsg string) error {
    return c.Render("error", fiber.Map{
        "title": "Error",
        "msg": errMsg,
    }, "layouts/main")
}


func (cl *cutlink) HomePage(c *fiber.Ctx) error {
    var urls []*models.Url

    id, err := cl.getUserID(c)
    if err != nil {
        cl.ErrorLog.Println(err.Error())
    }

    if id != 0 {
        urls, err = cl.Conn.GetAllUrls(id)
        if err != nil {
            cl.ErrorLog.Println(err.Error())
            return fiber.ErrInternalServerError
        }
    }

    err = c.Render("index", fiber.Map{
        "title": "Home",
        "Urls": urls,
        "authenticated": id,
    }, "layouts/main")

    if err != nil {
        cl.ErrorLog.Println(err.Error())
    }

    return err
}


func (cl *cutlink) SignupPage(c *fiber.Ctx) error {
    err := c.Render("signup", fiber.Map{
        "title": "Signup",
        "disabled": cl.DisableSignup,
    }, "layouts/main")

    if err != nil {
        cl.ErrorLog.Println(err.Error())
    }

    return err
}


func (cl *cutlink) SignupUser(c *fiber.Ctx) error {
    if cl.DisableSignup {
        return c.SendString("Signup is disabled.")
    }

    password := c.FormValue("password", "")
    if password == "" || len(password) <= 8 {
        retval := `<div class="container alert alert-danger" role="alert">
        <h4>Password is not valid.</h4>
        <p style="font-size: 20px;">Password must be more than 8 characters.</p>
        </div>`
        return c.SendString(retval)
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

    retval := fmt.Sprintf(
        `<div class="container alert alert-success" role="alert"><h3>UUID</h3><code style="font-size: 20px">%s</code></div>`,
        userID.String())
    err = c.SendString(retval)

    if err != nil {
        cl.ErrorLog.Println(err.Error())
    }

    return err
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

    err = c.Render("login", fiber.Map{
        "errMsg": errMsg,
        "title": "Login",
    }, "layouts/main")

    sess.Delete("errMsg")
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
    userID   := c.FormValue("uuid", "")

    if !userIdMatcher.Match([]byte(userID)) {
        sess.Set("errMsg", "Invalid UserID or Password.")
        sess.Save()
        return cl.LoginPage(c)
    }

    id, err := cl.Conn.AuthenticatUser(userID, password)
    if err != nil {
        sess.Set("errMsg", "Invalid UserID or Password.")
        sess.Save()
        return cl.LoginPage(c)
    }

    sess.Regenerate()
    sess.Set("authenticatedUserID", id)
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
    sess.Save()

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
        sess.Save()
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
            "hash": target.Hash,
        }, "layouts/main")
    }

    err = cl.Conn.IncrementClicked(hash)
    if err != nil {
        cl.ErrorLog.Println("Incrementing for", target.Hash, "Failed")
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

    return c.Redirect(target.Target, fiber.StatusSeeOther)
}


func (cl *cutlink) AddUrl(c *fiber.Ctx) error {
    sess, err := cl.Store.Get(c)
    if err != nil {
        cl.ErrorLog.Println(err.Error())
        return fiber.ErrInternalServerError
    }
    id := sess.Get("authenticatedUserID")

    target := c.FormValue("target", "")
    if target == "" || !urlMatcher.Match([]byte(target)) {
        return fiber.ErrInternalServerError
    }
    password := c.FormValue("password", "")

    _, _, err = cl.Conn.CreateUrl(id.(int), target, password)
    if (err != nil) {
        cl.ErrorLog.Println(err.Error())
        return fiber.ErrInternalServerError
    }

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

    return nil
}
