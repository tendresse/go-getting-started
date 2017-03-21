package models

import (
	// "encoding/json"
	_"gopkg.in/pg.v5"
)

type Tag struct {
	ID 		int 		`json:"id,omitempty"`
	Title 		string		`json:"title,omitempty"`
	Banned		bool		`json:"banned,omitempty"`
	Achievements    []*Achievement 	`json:"achievements,omitempty"`
	Gifs		[]Gif 		`pg:",many2many:gifs_tags,joinFK:Gif" json:"gifs,omitempty"`
}

// func (tag *Tag) get_random_tag() {
// 	db.Order("RANDOM()").Find(&tag)
// }