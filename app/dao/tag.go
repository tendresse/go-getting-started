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
func (c Tag) CreateTags(tags []*models.Tag) error {
	for _,tag := range tags {
		if err := c.CreateTag(tag); err != nil {
			return err
		}
	}
	return nil
}


func (c Tag) GetOrCreateTag(tag *models.Tag) error {
	_,err := c.DB.Model(&tag).
		Column("id").
		Where("title = ?", tag.Title).
		OnConflict("DO NOTHING"). // OnConflict is optional
		Returning("id").
		SelectOrInsert()
	return err
}


func (c Tag) GetBannedTags(tags []*models.Tag) error {
	count, err := c.DB.Model(&models.Tag{}).Count()
	if err != nil {
		return err
	}
	return c.DB.Model(tags).
		Where("banned = ?",true).
		Limit(count).
		Select()
}


func (c Tag) GetTag(tag *models.Tag) error {
	return c.DB.Select(&tag)
}
func (c Tag) GetTags(tags []*models.Tag) error {
	return c.DB.Model(&tags).Select()
}

func (c Tag) GetAllTags(tags *[]models.Tag) (error) {
	count, err := c.DB.Model(&models.Tag{}).Count()
	if err != nil {
		return err
	}
	return c.DB.Model(&tags).Limit(count).Select()
}


func (c Tag) UpdateTag(tag *models.Tag) error {
	return c.DB.Update(tag)
}


func (c Tag) DeleteTag(tag *models.Tag) error {
	return c.DB.Delete(&tag)
}
func (c Tag) DeleteTags(tags []*models.Tag) error {
	for _,tag := range tags {
		if err := c.DeleteTag(tag); err != nil {
			return err
		}
	}
	return nil
}


func (c Tag) GetFullTag(tag *models.Tag) error {
	return c.DB.Model(&tag).Column("tag.*", "Gifs").Column("tag.*", "Achievements").First()
}
func (c Tag) GetFullTags(tags []*models.Tag) error {
	return c.DB.Model(&tags).Column("tag.*", "Gifs").Column("tag.*", "Achievements").Select()
}


func (c Tag) GetTagByTitle(title string, tag *models.Tag) error {
	return c.DB.Model(&tag).Where("title = ?",title).Select()
}


func (c Tag) AddTagToGif(tag *models.Tag, gif *models.Gif) error {
	return c.DB.Insert(&models.GifsTags{tag.ID,gif.ID})
}
func (c Tag) DeleteTagFromGif(tag *models.Tag, gif *models.Gif) error {
	return c.DB.Delete(&models.GifsTags{tag.ID,gif.ID})
}