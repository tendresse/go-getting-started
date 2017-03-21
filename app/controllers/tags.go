// controllers/tags_controllers.go

package controllers

import (
	"encoding/json"
	"strings"

	"github.com/tendresse/go-getting-started/app/config"
	"github.com/tendresse/go-getting-started/app/models"

	log "github.com/Sirupsen/logrus"
	_ "github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	_ "github.com/lib/pq"
)

type TagsController struct {
}

// wrap with admin rights
func (c TagsController) GetTags() string {
	// TODO : CHOOSE WHAT JSON TO RETURN
	tags := []models.Tag{}
	if err := config.Global.DB.Find(&tags).Error; err != nil {
		log.Error(err)
		return `{"success":false, "error":"no tag found"}`
	}
	b, err := json.Marshal(tags)
	if err != nil {
		log.Error(err)
		return `{"success":false, "error":"marshal tags json error"}`
	}
	return strings.Join([]string{`{"success":true, "tags":`, string(b), "}"}, "")
}

//admin
func (c TagsController) GetTag(tag_id int) string {
	tag := models.Tag{}
	if err := config.Global.DB.First(&tag, tag_id).Error; err != nil {
		log.Error(err)
		return `{"success":false, "error":"tag not found"}`
	}
	b, err := json.Marshal(tag)
	if err != nil {
		log.Error(err)
		return `{"success":false, "error":"error while marshaling tag"}`
	}
	return strings.Join([]string{`{"success":true, "tag":`, string(b), "}"}, "")
}

//admin
func (c TagsController) AddTag(tag_json string) string {
	tag := models.Tag{}
	if err := json.Unmarshal([]byte(tag_json), &tag); err != nil {
		log.Error(err)
		return `{"success":false, "error":"error while marshaling tag json"}`
	}
	if config.Global.DB.Where("Name = ?", tag.Name).First(&tag).Error != nil {
		if err := config.Global.DB.Create(&tag).Error; err != nil {
			log.Error(err)
			return `{"success":false, "error":"error while creating Gif"}`
		}
		return strings.Join([]string{`{"success":true, "tag":`, string(tag.ID), "}"}, "")
	}
	return `{"success":false, "error":"tag already exists"}`
}

//admin
func (c TagsController) UpdateTag(tag_json string) string {
	tag := models.Tag{}
	updated_tag := models.Tag{}
	if err := json.Unmarshal([]byte(tag_json), &updated_tag); err != nil {
		log.Error(err)
		return `{"success":false, "error":"error while marshaling tag json"}`
	}
	// First is used without Preloading Tags in order to replace them
	if err := config.Global.DB.First(&tag, updated_tag.ID).Error; err != nil {
		log.Error(err)
		return `{"success":false, "error":"tag not found"}`
	}
	// we check if the new Title is not already taken
	if strings.Compare(updated_tag.Name, tag.Name) != 0 {
		if err := config.Global.DB.Where("Title = ?", updated_tag.Name).First(&tag).Error; err != nil {
			log.Error(err)
			return `{"success":false, "error":"tag name already taken"}`
		}
	}
	config.Global.DB.Save(&tag)
	return `{"success":true}`
}

//admin
func (c TagsController) DeleteTag(tag_id int) string {
	tag := models.Tag{}
	if err := config.Global.DB.Delete(&tag, tag_id).Error; err != nil {
		log.Error(err)
		return `{"success":false, "error":"tag already deleted"}`
	}
	return `{"success":true}`
}
