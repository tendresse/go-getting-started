package models

import (
	// "encoding/json"
	_ "gopkg.in/pg.v5"
)

type Role struct {
	ID 	int		`json:"id,omitempty"`
	Title 	string		`json:"title,omitempty"`
	Users 	[]*User		`pg:",many2many:users_roles,joinFK:User" json:"users,omitempty"`
}