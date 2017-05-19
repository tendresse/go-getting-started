package config

import (
	"github.com/go-pg/pg"
)

type GlobalStruct struct {
	DB 		*pg.DB
	SecretKey	string
	DatabaseURI	string
	TumblrAPIKey	string
}

var Global GlobalStruct