package dao

import (
	"github.com/tendresse/go-getting-started/app/models"
	"gopkg.in/pg.v5"
	//log "github.com/Sirupsen/logrus"
	"fmt"
)

type Gif struct {
	DB 	*pg.DB
}


func (c Gif) CreateGif(gif *models.Gif) error {
	return c.DB.Insert(gif)
}
func (c Gif) CreateGifs(gifs []*models.Gif) error {
	for _,gif := range gifs {
		if err := c.CreateGif(gif); err != nil {
			return err
		}
	}
	return nil
}

func (c Gif) GetOrCreateGif(gif *models.Gif) error {
	_,err := c.DB.Model(&gif).
		Column("id").
		Where("url = ?", gif.Url).
		OnConflict("DO NOTHING"). // OnConflict is optional
		Returning("id").
		SelectOrInsert()
	return err
}

func (c Gif) UpdateGif(gif *models.Gif) error {
	return c.DB.Update(gif)
}

func (c Gif) GetRandomGif(gif *models.Gif) error {
	// get random tag
	var tag_id int
	var gif_id int
	for{
		r, _ := c.DB.Query(&tag_id, "SELECT t.id from tags t where NOT t.banned offset random() * (select count(*) from tags) limit 1 ;")
		if r.RowsReturned() == 0 {
			// Tag has no Gif
			continue
		}
		// get random gif from random tag
		r, _ = c.DB.Query(&gif_id, "SELECT gif_id from gifs_tags where tag_id = ? offset random() * (select count(*) from gifs_tags where tag_id = ?) limit 1 ;",tag_id,tag_id)
		fmt.Println(tag_id,gif_id)
		if r.RowsReturned() == 0 {
			// Tag has no Gif
			continue
		}
		break
	}
	gif.ID = gif_id
	return c.GetFullGif(gif)
}

func (c Gif) GetGif(gif *models.Gif) error {
	return c.DB.Select(&gif)
}
func (c Gif) GetGifs(gifs []*models.Gif) error {
	return c.DB.Model(&gifs).Select()
}
func (c Gif) GetAllGifs(gifs *[]models.Gif) error {
	count, err := c.DB.Model(&models.Gif{}).Count()
	if err != nil {
		return err
	}
	return c.DB.Model(&gifs).Limit(count).Select()
}


func (c Gif) DeleteGif(gif *models.Gif) error {
	return c.DB.Delete(&gif)
}
func (c Gif) DeleteGifs(gifs []*models.Gif) error {
	for _,gif := range gifs {
		if err := c.DeleteGif(gif); err != nil {
			return err
		}
	}
	return nil
}


func (c Gif) GetGifByUrl(gif *models.Gif) error {
	return c.DB.Model(&gif).Where("url = ?",gif.Url).First()
}


func (c Gif) GetFullGif(gif *models.Gif) error {
	return c.DB.Model(&gif).Column("gif.*", "Blog").Column("gif.*", "Tags").First()
}
func (c Gif) GetFullGifs(gifs []*models.Gif) error {
	count, err := c.DB.Model(&models.Gif{}).Count()
	if err != nil {
		return err
	}
	return c.DB.Model(&gifs).Column("gif.*", "Blog").Column("gif.*", "Tags").Limit(count).Select()
}



func (c Gif) AddTagToGif(tag *models.Tag, gif *models.Gif) error {
	return c.DB.Insert(&models.GifsTags{tag.ID,gif.ID})
}
func (c Gif) DeleteTagFromGif(tag *models.Tag, gif *models.Gif) error {
	return c.DB.Delete(&models.GifsTags{tag.ID,gif.ID})
}