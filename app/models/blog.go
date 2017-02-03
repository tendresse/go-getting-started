package models

import (
	// "encoding/json"
	_"gopkg.in/pg.v5"
)

type Blog struct{
	ID 		int 		`json:"id,omitempty"`
	Name        	string		`json:"name,omitempty"`
	Url         	string 		`json:"url,omitempty"`
	Gifs 		[]*Gif		`json:"gifs,omitempty"`
}

