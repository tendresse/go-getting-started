// controllers/gifs_controllers.go

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

type GifsController struct {
}

func (c GifsController) RandomGif() models.Gif {
	gif := models.Gif{}
	if err := config.Global.DB.Preload("Tags").Order("RANDOM()").First(&gif).Error; err != nil {
		log.Error(err)
		return nil
	}
	return gif
}

func (c GifsController) RandomJSONGif() string {
	gif := models.Gif{}
	if err := config.Global.DB.Preload("Tags").Order("RANDOM()").First(&gif).Error; err != nil {
		log.Error(err)
		return `{"success":false, "error":"no gif found"}`
	}
	b, err := json.Marshal(gif)
	if err != nil {
		log.Error(err)
		return `{"success":false, "error":"error while marshaling gif"}`
	}
	return strings.Join([]string{`{"success":true, "gif":`, string(b), "}"}, "")

	// return '{"Id":3,"BlogID":2,"Url":"http://i.giphy.com/xTgeIYpiaiWvWHyWwU.gif","LameScore":4}'
}

//admin
func (c GifsController) GetGif(gif_id int) string {
	gif := models.Gif{}
	if err := config.Global.DB.First(&gif, gif_id).Error; err != nil {
		log.Error(err)
		return `{"success":false, "error":"gif not found"}`
	}
	b, err := json.Marshal(gif)
	if err != nil {
		log.Error(err)
		return `{"success":false, "error":"error while marshaling gif"}`
	}
	return strings.Join([]string{`{"success":true, "gif":`, string(b), "}"}, "")
}

//admin
func (c GifsController) SearchGifsByTags(tags_json string) string {
	// db.Limit(3).Find(&users)
	// db.Preload("Orders", "state = ?", "paid").Preload("Orders.OrderItems").Find(&users)
	type TagsJSON struct {
		Tags []models.Tag `json:"tags,omitempty"`
	}
	tags := TagsJSON{}
	if err := json.Unmarshal([]byte(tags_json), &tags); err != nil {
		log.Error(err)
	}
	gifs := []models.Gif{}
	// loop for finding the gifs matching ALL the tags given
	for i, _ := range tags.Tags {
		if err := config.Global.DB.Where(models.Tag{Name: tags.Tags[i].Name}).First(&tags.Tags[i]).Error; err != nil {
			log.Error(err)
			continue
		}
		if i == 0 {
			if err := config.Global.DB.Debug().Model(&tags.Tags[i]).Related(&gifs, "Gifs").Error; err != nil {
				log.Error(err)
				return `{"success":false, "error":"error while fetching data"}`
			}
		} else {
			if err := config.Global.DB.Debug().Model(&tags.Tags[i]).Where(&gifs).Related(&gifs, "Gifs").Error; err != nil {
				log.Error(err)
				return `{"success":false, "error":"error while fetching data"}`
			}
		}
	}
	b, err := json.Marshal(gifs)
	if err != nil {
		log.Error(err)
		return `{"success":false, "error":"marshal achievement json error"}`
	}
	return strings.Join([]string{`{"success":true, "gifs":`, string(b), "}"}, "")
}

//admin
func (c GifsController) AddGif(gif_json string) string {
	gif := models.Gif{}
	if err := json.Unmarshal([]byte(gif_json), &gif); err != nil {
		log.Error(err)
		return `{"success":false, "error":"error while marshaling gif json"}`
	}
	if config.Global.DB.Where("Url = ?", gif.Url).First(&gif).Error != nil {
		for i, _ := range gif.Tags {
			if err := config.Global.DB.Where(models.Tag{Name: gif.Tags[i].Name}).FirstOrCreate(&gif.Tags[i]).Error; err != nil {
				log.Error(err)
				return `{"success":false, "error":"error while fetching Tags"}`
			}
		}
		if err := config.Global.DB.Create(&gif).Error; err != nil {
			log.Error(err)
			return `{"success":false, "error":"error while creating Gif"}`
		}
		return strings.Join([]string{`{"success":true, "gif":`, string(gif.ID), "}"}, "")
	}
	return `{"success":false, "error":"gif already exists"}`
}

//admin
func (c GifsController) UpdateGif(gif_json string) string {
	gif := models.Gif{}
	updated_gif := models.Gif{}
	if err := json.Unmarshal([]byte(gif_json), &updated_gif); err != nil {
		log.Error(err)
		return `{"success":false, "error":"error while marshaling gif json"}`
	}
	// First is used without Preloading Tags in order to replace them
	if err := config.Global.DB.First(&gif, updated_gif.ID).Error; err != nil {
		log.Error(err)
		return `{"success":false, "error":"gif not found"}`
	}
	// we check if the new Title is not already taken
	if strings.Compare(updated_gif.Url, gif.Url) != 0 {
		if err := config.Global.DB.Where("Title = ?", updated_gif.Url).First(&gif).Error; err != nil {
			log.Error(err)
			return `{"success":false, "error":"gif url already taken"}`
		}
	}
	for i, _ := range gif.Tags {
		if err := config.Global.DB.Where(models.Tag{Name: gif.Tags[i].Name}).FirstOrCreate(&gif.Tags[i]).Error; err != nil {
			log.Error(err)
			return `{"success":false, "error":"erro while fetching Tags"}`
		}
	}
	config.Global.DB.Save(&gif)
	config.Global.DB.Model(&gif).Association("Tags").Replace(&updated_gif.Tags)
	return `{"success":true}`
}

//admin
func (c GifsController) DeleteGif(gif_id int) string {
	gif := models.Achievement{}
	if err := config.Global.DB.Delete(&gif, gif_id).Error; err != nil {
		log.Error(err)
		return `{"success":false, "error":"gif already deleted"}`
	}
	return `{"success":true}`
}
