package models

import (
	// "encoding/json"
	_"github.com/go-pg/pg"
)

type Gif struct {
	Id         int 		`json:"id,omitempty"`
	BlogId     int		`json:"blog_id,omitempty"`
	Blog       Blog		`json:"blog,omitempty"`
	Url        string 	`json:"url,omitempty"`
	LameScore  int		`json:"lame_score,omitempty"`
	Tags       []Tag 	`pg:",many2many:gifs_tags,joinFK:Tag" json:"tags,omitempty"`
	Tendresses []*Tendresse	`json:"tendresses,omitempty"`
}

