// controllers/tags_controllers.go

package controllers

import (
	"encoding/json"
	"strings"

	"github.com/tendresse/go-getting-started/app/config"
	"github.com/tendresse/go-getting-started/app/dao"
	"github.com/tendresse/go-getting-started/app/models"

	log "github.com/Sirupsen/logrus"
)

type Tag struct {
}

// wrap with admin rights
func (c Tag) GetTags() string {
	// TODO : CHOOSE WHAT JSON TO RETURN
	tags_dao := dao.Tag{DB:config.Global.DB}
	tags := []models.Tag{}
	if err := tags_dao.GetAllTags(&tags); err != nil {
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
func (c Tag) GetTag(tag_id int) string {
	tags_dao := dao.Tag{DB:config.Global.DB}
	tag := models.Tag{ID:tag_id}
	if err := tags_dao.GetFullTag(&tag); err != nil {
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
func (c Tag) AddTag(tag_json string) string {
	tags_dao := dao.Tag{DB:config.Global.DB}
	tag := models.Tag{}
	if err := json.Unmarshal([]byte(tag_json), &tag); err != nil {
		log.Error(err)
		return `{"success":false, "error":"error while marshaling tag json"}`
	}
	if err := tags_dao.GetTagByTitle(tag.Title, &tag); err == nil {
		log.Error(err)
		return `{"success":false, "error":"tag already exists"}`
	}
	if err := tags_dao.CreateTag(&tag); err != nil {
		log.Error(err)
		return `{"success":false, "error":"error while creating Gif"}`
	}
	return strings.Join([]string{`{"success":true, "tag":`, string(tag.ID), "}"}, "")
}

//admin
func (c Tag) UpdateTag(tag_json string) string {
	tags_dao := dao.Tag{DB:config.Global.DB}
	updated_tag := models.Tag{}
	if err := json.Unmarshal([]byte(tag_json), &updated_tag); err != nil {
		log.Error(err)
		return `{"success":false, "error":"error while marshaling tag json"}`
	}
	tag := models.Tag{ID:updated_tag.ID}
	if err := tags_dao.GetTag(&tag); err != nil {
		log.Error(err)
		return `{"success":false, "error":"tag not found"}`
	}
	// if Title has changed we check if not already taken
	if strings.Compare(updated_tag.Title, tag.Title) != 0 {
		if err := tags_dao.GetTagByTitle(updated_tag.Title,&updated_tag); err == nil {
			log.Error(err)
			return `{"success":false, "error":"tag name already taken"}`
		}
	}
	tag.Title = updated_tag.Title
	if err := tags_dao.UpdateTag(&updated_tag); err == nil {
		log.Error(err)
		return `{"success":false, "error":"error while updating tag"}`
	}
	return `{"success":true}`
}

//admin
func (c Tag) DeleteTag(tag_id int) string {
	tags_dao := dao.Tag{DB:config.Global.DB}
	tag := models.Tag{ID:tag_id}
	if err := tags_dao.DeleteTag(&tag).Error; err != nil {
		log.Error(err)
		return `{"success":false, "error":"tag already deleted"}`
	}
	return `{"success":true}`
}
