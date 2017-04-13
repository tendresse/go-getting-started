package config

import (
	"github.com/googollee/go-socket.io"
	"gopkg.in/pg.v5"
)

type GlobalStruct struct {
	DB 		*pg.DB
	Server 		*socketio.Server
	SecretKey	string
	DatabaseURI	string
	TumblrAPIKey	string
}

var Global GlobalStruct