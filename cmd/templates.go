package main

import (
    "cutlink/models"
)

type templateData struct {
    Url  *models.Url
    Urls []*models.Url
    // Token string
}
