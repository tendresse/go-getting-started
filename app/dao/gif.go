package dao

import (
	"github.com/tendresse/go-getting-started/app/models"
	"gopkg.in/pg.v5"
	//log "github.com/Sirupsen/logrus"
)

type Gif struct {
	DB 	*pg.DB
}

func (c Gif) CreateGif(gif *models.Gif) error {
	return c.DB.Insert(gif)
}

func (c Gif) CreateGifs(gifs *[]models.Gif) error {
	return c.DB.Insert(gifs)
}

func (c Gif) UpdateGif(gif *models.Gif) error {
	return c.DB.Update(gif)
}

func (c Gif) GetGif(gif *models.Gif) error {
	return c.DB.Select(&gif)
}

func (c Gif) DeleteGif(gif *models.Gif) error {
	return c.DB.Delete(&gif)
}

func (c Gif) GetFullGif(gif *models.Gif) error {
	return c.DB.Model(&gif).Column("gif.*", "Blog").Column("gif.*", "Tags").First()
}

func (c Gif) GetGifs(gifs *[]models.Gif) error {
	return c.DB.Model(&gifs).Select()
}

func (c Gif) GetFullGifs(gifs *[]models.Gif) error {
	return c.DB.Model(&gifs).Column("gif.*", "Blog").Column("gif.*", "Tags").Select()
}

func (c Gif) DeleteGifs(gifs *[]models.Gif) error {
	return c.DB.Delete(&gifs)
}

func (c Gif) AddTagToGif(tag *models.Tag, gif *models.Gif) error {
	return c.DB.Insert(&models.GifsTags{tag.ID,gif.ID})
}

func (c Gif) AddTagsToGif(tags *[]models.Tag, gif *models.Gif) error {
	for _,tag := range *tags {
		if err := c.AddTagToGif(&tag, gif); err != nil {
			return err
		}
	}
	return nil
}