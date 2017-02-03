package tests

import (
	_ "github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	_ "github.com/lib/pq"
)

type Tag struct {
	ID 		int 	`json:"id,omitempty"`
	Name 		string 	`json:"name,omitempty"`
	Banned		bool	`json:"banned,omitempty"`
}
