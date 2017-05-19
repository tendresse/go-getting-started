package dao

import (
	"github.com/tendresse/go-getting-started/app/config"
	"github.com/tendresse/go-getting-started/app/models"
)

type Tag struct {}


func (c Tag) CreateTag(tag *models.Tag) error {
	return config.Global.DB.Insert(tag)
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
	_,err := config.Global.DB.Model(&tag).
		Column("id").
		Where("title = ?", tag.Title).
		OnConflict("DO NOTHING"). // OnConflict is optional
		Returning("id").
		SelectOrInsert()
	return err
}


func (c Tag) GetBannedTags(tags []*models.Tag) error {
	count, err := config.Global.DB.Model(&models.Tag{}).Count()
	if err != nil {
		return err
	}
	return config.Global.DB.Model(tags).
		Where("banned = ?",true).
		Limit(count).
		Select()
}


func (c Tag) GetTag(tag *models.Tag) error {
	return config.Global.DB.Select(&tag)
}
func (c Tag) GetTags(tags []*models.Tag) error {
	return config.Global.DB.Model(&tags).Select()
}

func (c Tag) GetAllTags(tags *[]models.Tag) (error) {
	count, err := config.Global.DB.Model(&models.Tag{}).Count()
	if err != nil {
		return err
	}
	return config.Global.DB.Model(&tags).Limit(count).Select()
}


func (c Tag) UpdateTag(tag *models.Tag) error {
	return config.Global.DB.Update(tag)
}


func (c Tag) DeleteTag(tag *models.Tag) error {
	return config.Global.DB.Delete(&tag)
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
	return config.Global.DB.Model(&tag).Column("tag.*", "Gifs").Column("tag.*", "Achievements").Where("tag.id = ?",tag.Id).First()
}
func (c Tag) GetFullTags(tags []*models.Tag) error {
	return config.Global.DB.Model(&tags).Column("tag.*", "Gifs").Column("tag.*", "Achievements").Select()
}


func (c Tag) GetTagByTitle(title string, tag *models.Tag) error {
	return config.Global.DB.Model(&tag).Where("title = ?",title).Select()
}


func (c Tag) AddTagToGif(tag *models.Tag, gif *models.Gif) error {
	return config.Global.DB.Insert(&models.GifsTags{tag.Id, gif.Id})
}
func (c Tag) DeleteTagFromGif(tag *models.Tag, gif *models.Gif) error {
	return config.Global.DB.Delete(&models.GifsTags{tag.Id, gif.Id})
}