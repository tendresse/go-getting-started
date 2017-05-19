package models

import (
	// "encoding/json"
	_ "github.com/go-pg/pg"
)

type Role struct {
	Id    int	`json:"id,omitempty"`
	Title string	`json:"title,omitempty"`
	Users []*User	`pg:",many2many:users_roles,joinFK:User" json:"users,omitempty"`
}