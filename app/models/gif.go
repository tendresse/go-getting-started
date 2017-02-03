package models

import (
	// "encoding/json"
	_"gopkg.in/pg.v5"
)
type Gif struct {
	ID 		int 		`json:"id,omitempty"`
	BlogID		int		`json:"blog_id,omitempty"`
	Blog     	Blog		`json:"blog,omitempty"`
	Url		string 		`json:"url,omitempty"`
	LameScore	int		`json:"lame_score,omitempty"`
	Tags		[]Tag 		`pg:",many2many:item_to_items,joinFK:Tag" json:"tags,omitempty"`
	Tendresses 	[]*Tendresse	`json:"tendresses,omitempty"`
}

