package models

import (
	// "encoding/json"
	_ "gopkg.in/pg.v5"
)

type Role struct {
	ID 	int		`json:"id,omitempty"`
	Name 	string		`json:"name,omitempty"`
	Users 	[]User		`pg:",many2many:item_to_items,joinFK:User" json:"users,omitempty"`
}