package main

import (
    "fmt"
    "strings"
    "net/http"
    "html/template"

    "github.com/julienschmidt/httprouter"
)


func (a *app) Root(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
    urls, err := a.Urls.GetAll()
    if err != nil {
        a.ErrorLog.Println(err.Error())
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }

    files := []string{
        "./ui/html/base.tmpl.html",
        "./ui/html/partials/nav.tmpl.html",
        "./ui/html/partials/footer.tmpl.html",
        "./ui/html/partials/urltable.tmpl.html",
        "./ui/html/partials/urlrow.tmpl.html",
        "./ui/html/pages/home.tmpl.html",
    }

    ts, err := template.ParseFiles(files...)
    if err != nil {
        a.ErrorLog.Println(err.Error())
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }

    data := &templateData{
        Urls: urls,
    }

    err = ts.ExecuteTemplate(w, "base", data)
    if err != nil {
        a.ErrorLog.Println(err.Error())
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
    }
}


func (a *app) Redirector(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    hash := ps.ByName("hash")
    if hash == "" {
        a.ErrorLog.Printf("Cannot find hash %s", hash)
        http.NotFound(w, r)
        return
    }

    target, err := a.Urls.Get(hash)
    if err != nil {
        a.ErrorLog.Println(err.Error())
        http.NotFound(w, r)
        return
    }

    tUrl := target.Target
    if !strings.HasPrefix(tUrl, "https://") {
        tUrl = fmt.Sprintf("https://%s", tUrl)
    }

    err = a.Urls.IncrementClicked(hash)
    if err != nil {
        a.ErrorLog.Println("Incrementing for", target.Hash, "Failed")
    }

    http.Redirect(w, r, tUrl, http.StatusSeeOther)
}


func (a *app) ViewUrl(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
    hash := r.URL.Query().Get("hash")
    if hash == "" {
        http.NotFound(w, r)
        return
    }

    shortUrl, err := a.Urls.Get(hash)
    if err != nil {
        a.ErrorLog.Println(err)
        http.NotFound(w, r)
        return
    }

    w.Header().Set("HX-Trigger", "newUrl")
    ts, _ := template.ParseFiles("./ui/html/partials/urlrow.tmpl.html")
    _ = ts.ExecuteTemplate(w, "urlrow", shortUrl)
}


func (a *app) ViewAll(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
    urls, err := a.Urls.GetAll()
    if (err != nil) {
        a.ErrorLog.Println(err.Error())
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }

    files := []string{
        "./ui/html/partials/urltable.tmpl.html",
        "./ui/html/partials/urlrow.tmpl.html",
    }

    data := &templateData{
        Urls: urls,
    }

    ts, _ := template.ParseFiles(files...)
    _ = ts.ExecuteTemplate(w, "urltable", data)
}


func (a *app) AddUrl(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
    err := r.ParseForm()
    if err != nil {
        a.ErrorLog.Println(err.Error())
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }
    target := r.Form["target"][0]

    _, hash, err := a.Urls.Create(target)
    if (err != nil) {
        a.ErrorLog.Println(err.Error())
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }

    a.InfoLog.Println("New URL added:", target)
    http.Redirect(w, r, fmt.Sprintf("/api/view?hash=%s", hash), http.StatusSeeOther)
}


func (a *app) Delete(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    hash := ps.ByName("hash")
    err := a.Urls.Delete(hash)
    if err != nil {
        a.ErrorLog.Println(err.Error())
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }

    a.InfoLog.Printf("url with hash %s has been deleted", hash)
}
