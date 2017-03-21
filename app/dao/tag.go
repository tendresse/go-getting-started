package dao

import (
    "github.com/tendresse/go-getting-started/app/models"

    "gopkg.in/pg.v5"
)

type Tag struct {
	DB 	*pg.DB
}

func (c Tag) CreateTag(tag *models.Tag) error {
	return c.DB.Insert(tag)
}

func (c Tag) CreateTags(tags *[]models.Tag) error {
	return c.DB.Insert(tags)
}

func (c Tag) UpdateTag(tag *models.Tag) error {
	return c.DB.Update(tag)
}

func (c Tag) GetTag(tag *models.Tag) error {
	return c.DB.Select(&tag)
}

func (c Tag) DeleteTag(tag *models.Tag) error {
	return c.DB.Delete(&tag)
}

func (c Tag) GetFullTag(tag *models.Tag) error {
	return c.DB.Model(&tag).Column("tag.*", "Gifs").Column("tag.*", "Achievements").First()
}

func (c Tag) GetTagByTitle(title string, tag *models.Tag) error {
	return c.DB.Model(&tag).Where("title = ?",title).Select()
}

func (c Tag) GetTags(tags *[]models.Tag) error {
	return c.DB.Model(&tags).Select()
}

func (c Tag) GetFullTags(tags *[]models.Tag) error {
	return c.DB.Model(&tags).Column("tag.*", "Gifs").Column("tag.*", "Achievements").Select()
}

func (c Tag) DeleteTags(tags *[]models.Tag) error {
	return c.DB.Delete(&tags)
}

func (c Tag) AddTagToGif(tag *models.Tag, gif *models.Gif) error {
	return c.DB.Insert(&models.GifsTags{tag.ID,gif.ID})
}

func (c Tag) AddTagsToGif(tags *[]models.Tag, gif *models.Gif) error {
	for _,tag := range *tags {
		if err := c.AddTagToGif(&tag, gif); err != nil {
			return err
		}
	}
	return nil
}