// controllers/gifs_controllers.go

package controllers

import (
	"encoding/json"
	"strings"

	"github.com/tendresse/go-getting-started/app/config"
	"github.com/tendresse/go-getting-started/app/dao"
	"github.com/tendresse/go-getting-started/app/models"

	log "github.com/Sirupsen/logrus"
)

type GifsController struct {
}

func (c GifsController) RandomGif() string {
	gifs_dao := dao.Gif{DB:config.Global.DB}
	gif := models.Gif{}
	if err := gifs_dao.GetRandomGif(&gif); err != nil {
		log.Error(err)
		return `{"success":false, "error":"error while getting gif"}`
	}
	b, err := json.Marshal(gif)
	if err != nil {
		log.Error(err)
		return `{"success":false, "error":"error while marshaling gif"}`
	}
	return strings.Join([]string{`{"success":true, "gif":`, string(b), "}"}, "")
}

//admin
func (c GifsController) GetGif(gif_id int) string {
	gifs_dao := dao.Gif{DB:config.Global.DB}
	gif := models.Gif{ID:gif_id}
	if err := gifs_dao.GetFullGif(&gif); err != nil {
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
func (c GifsController) GetGifs() string {
	gifs_dao := dao.Gif{DB:config.Global.DB}
	gifs, err := gifs_dao.GetAllGifs()
	if err != nil {
		log.Error(err)
		return `{"success":false, "error":"gif not found"}`
	}
	b, err := json.Marshal(gifs)
	if err != nil {
		log.Error(err)
		return `{"success":false, "error":"error while marshaling gifs"}`
	}
	return strings.Join([]string{`{"success":true, "gifs":`, string(b), "}"}, "")
}

//admin
func (c GifsController) SearchGifsByTags(tags_json string) string {
	// TODO : test SearchGifsByTags
	type TagsJSON struct {
		Tags []models.Tag `json:"tags,omitempty"`
	}
	tags := TagsJSON{}
	if err := json.Unmarshal([]byte(tags_json), &tags); err != nil {
		log.Error(err)
	}
	tags_dao := dao.Tag{DB:config.Global.DB}
	gifs     := []models.Gif{}
	final_gifs := []models.Gif{}
	first := true
	// for each tag, get the gifs corresponding
	for i,_ := range tags.Tags {
		if err := tags_dao.GetTagByTitle(tags.Tags[i].Title,&tags.Tags[i]); err != nil {
			log.Error(err)
			continue
		}
		if err := tags_dao.GetFullTag(&tags.Tags[i]); err != nil {
			log.Error(err)
			return `{"success":false, "error":"error while fetching data"}`
		}
		if first {
			gifs = tags.Tags[i].Gifs
			first = false
		} else {
			CommonTags:
			for j,_ := range gifs {
				for k,_ := range tags.Tags[i].Gifs {
					if gifs[j].ID == tags.Tags[i].Gifs[k].ID{
						final_gifs = append(final_gifs, gifs[j])
						continue CommonTags
					}
				}
			}
			gifs = final_gifs
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
	gifs_dao := dao.Gif{DB:config.Global.DB}

	gif := models.Gif{}
	if err := json.Unmarshal([]byte(gif_json), &gif); err != nil {
		log.Error(err)
		return `{"success":false, "error":"error while marshaling gif json"}`
	}
	if err := gifs_dao.GetGifByUrl(&gif); err == nil {
		log.Error(err)
		return `{"success":false, "error":"gif already exists"}`
	}
	if err := gifs_dao.CreateGif(&gif); err != nil {
		log.Error(err)
		return `{"success":false, "error":"error while creating gif"}`
	}
	tags_dao := dao.Tag{DB:config.Global.DB}
	for i, _ := range gif.Tags {
		if err := tags_dao.GetOrCreateTag(&gif.Tags[i]); err != nil {
			log.Error(err)
			return `{"success":false, "error":"error while creating Tag"}`
		}
		if err := gifs_dao.AddTagToGif(&models.Tag{ID:gif.Tags[i].ID}, &gif); err != nil {
			log.Error(err)
			return `{"success":false, "error":"error while creating Tag"}`
		}
	}
	return strings.Join([]string{`{"success":true, "gif":`, string(gif.ID), "}"}, "")
}

//admin
func (c GifsController) UpdateGif(gif_json string) string {
	gifs_dao := dao.Gif{DB:config.Global.DB}
	tags_dao := dao.Tag{DB:config.Global.DB}

	updated_gif := models.Gif{}
	if err := json.Unmarshal([]byte(gif_json), &updated_gif); err != nil{
		log.Error(err)
		return `{"success":false, "error":"marshal achievement json error"}`
	}
	gif := models.Gif{ID:updated_gif.ID}
	if err := gifs_dao.GetGif(&gif); err != nil {
		log.Error(err)
		return `{"success":false, "error":"achievement not found"}`
	}

	FetchTags:
	for i, _ := range updated_gif.Tags {
		if err := tags_dao.GetOrCreateTag(&updated_gif.Tags[i]); err != nil {
			log.Error(err)
		}
		for j,_ := range gif.Tags {
			if updated_gif.Tags[i].ID == gif.Tags[j].ID {
				continue FetchTags
			}
		}
		// new tag so should be linked
		gifs_dao.AddTagToGif(&updated_gif.Tags[i], &gif)
	}
	CurrentTags:
	for i,_ := range gif.Tags {
		for j, _ := range updated_gif.Tags {
			if updated_gif.Tags[j].ID == gif.Tags[i].ID {
				continue CurrentTags
			}
		}
		// new tag so should be unlinked
		gifs_dao.DeleteTagFromGif(&gif.Tags[i],&gif)
	}
	return strings.Join([]string{`{"success":true, "gif":` , string(gif.ID) , "}"} , "")
}

//admin
func (c GifsController) DeleteGif(gif_id int) string {
	gifs_dao := dao.Gif{DB : config.Global.DB}
	if err := gifs_dao.DeleteGif(&models.Gif{ID:gif_id}); err != nil {
		log.Error(err)
		return `{"success":false, "error":"error while deleting gif"}`
	}
	return `{"success":true}`
}
