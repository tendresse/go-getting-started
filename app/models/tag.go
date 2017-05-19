package models

import (
	// "encoding/json"
	_"github.com/go-pg/pg"
)

type Tag struct {
	Id           int 		`json:"id,omitempty"`
	Title        string		`json:"title,omitempty"`
	Banned       bool		`json:"banned,omitempty"`
	Achievements []*Achievement 	`json:"achievements,omitempty"`
	Gifs         []Gif 		`pg:",many2many:gifs_tags,joinFK:Gif" json:"gifs,omitempty"`
}

// func (tag *Tag) get_random_tag() {
// 	db.Order("RANDOM()").Find(&tag)
// }