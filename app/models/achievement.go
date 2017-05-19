package models

import (
    	// "encoding/json"
	_"github.com/go-pg/pg"
)

type Achievement struct{
	Id        int 		`json:"id,omitempty"`
	Condition int 		`json:"condition,omitempty"`
	Icon      string 	`json:"icon,omitempty"`
	TagId     int		`json:"tag_id,omitempty"`
	Tag       *Tag    	`json:"tag,omitempty"`
	Title     string 	`json:"title,omitempty"`
	TypeOf    string 	`json:"type_of,omitempty"`
	XP        int 		`json:"xp,omitempty"`
}
