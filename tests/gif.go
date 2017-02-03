package tests

import (
	_ "github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	_ "github.com/lib/pq"
)


type Gif struct {
	ID 		int 	`json:"id,omitempty"`
	Url		string 	`json:"url,omitempty"`
	BlogID		int 	`json:"blog_id,omitempty"`
	Lamescore	int	`json:"lamescore,omitempty"`
	Tags 		[]Tag 	`pg:",many2many:gifs_tags,joinFK:Gif"`
}
