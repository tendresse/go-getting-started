package config

import (
	"github.com/tendresse/go-getting-started/app/models"

	"github.com/googollee/go-socket.io"
	"gopkg.in/pg.v5"
)

type GlobalStruct struct {
	DB 		*pg.DB
	CurrentUser 	models.User
	Server 		*socketio.Server
	Socket		socketio.Socket
	SecretKey	string
	DatabaseURI	string
	TumblrAPIKey	string
}

var Global GlobalStruct