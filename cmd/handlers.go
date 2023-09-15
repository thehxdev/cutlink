package main

import (
    "fmt"
    "strings"
    "net/http"
    "encoding/json"

    "github.com/julienschmidt/httprouter"
)


func (a *app) Root(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
    fmt.Fprintf(w, "Welcome to CutLink!\n")
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
    token := r.Header.Get("Token")
    if token != a.AdminToken || token == "" {
        a.ErrorLog.Println("Not authorized request:", token)
        http.Error(w, "Not Authorized", http.StatusInternalServerError)
        return
    }

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

    jsonData, err := json.Marshal(shortUrl)
    if err != nil {
        a.ErrorLog.Println(err.Error())
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.Write(jsonData)
}


func (a *app) ViewAll(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
    token := r.Header.Get("Token")
    if token != a.AdminToken || token == "" {
        a.ErrorLog.Println("Not authorized request:", token)
        http.Error(w, "Not Authorized", http.StatusInternalServerError)
        return
    }

    urls, err := a.Urls.GetAll()
    if (err != nil) {
        a.ErrorLog.Println(err.Error())
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }

    jsonData, err := json.Marshal(urls)
    if err != nil {
        a.ErrorLog.Println(err.Error())
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.Write(jsonData)
}


func (a *app) AddUrl(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
    token := r.Header.Get("Token")
    if token != a.AdminToken || token == "" {
        a.ErrorLog.Println("Not authorized request:", token)
        http.Error(w, "Not Authorized", http.StatusInternalServerError)
        return
    }

    // if r.Method != http.MethodPost {
    //     w.Header().Set("Allow", "POST")
    //     http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
    //     return
    // }

    target := r.URL.Query().Get("target")
    // customName := r.URL.Query().Get("name")
    if target == "" {
        http.Error(w, "Empty query", http.StatusInternalServerError)
        return
    }

    _, hash, err := a.Urls.Create(target)
    if (err != nil) {
        a.ErrorLog.Println(err.Error())
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }

    a.InfoLog.Println("New URL added:", target)
    http.Redirect(w, r, fmt.Sprintf("/view?hash=%s", hash), http.StatusSeeOther)
}


func (a *app) Delete(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    token := r.Header.Get("Token")
    if token != a.AdminToken || token == "" {
        a.ErrorLog.Println("Not authorized request:", token)
        http.Error(w, "Not Authorized", http.StatusInternalServerError)
        return
    }

    err := a.Urls.Delete(ps.ByName("hash"))
    if err != nil {
        a.ErrorLog.Println(err.Error())
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }

    w.Write([]byte("Deleted!"))
}
