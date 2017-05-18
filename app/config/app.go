package config

import (
	"gopkg.in/pg.v5"
)

type GlobalStruct struct {
	DB 		*pg.DB
	SecretKey	string
	DatabaseURI	string
	TumblrAPIKey	string
}

var Global GlobalStruct