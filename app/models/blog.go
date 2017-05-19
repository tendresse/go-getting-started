package models

import (
	// "encoding/json"
	_"github.com/go-pg/pg"
)

type Blog struct{
	Id    int 	`json:"id,omitempty"`
	Title string	`json:"title,omitempty"`
	Url   string 	`json:"url,omitempty"`
	Gifs  []*Gif	`json:"gifs,omitempty"`
}

