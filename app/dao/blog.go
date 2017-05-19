package dao

import (
	"github.com/tendresse/go-getting-started/app/config"
	"github.com/tendresse/go-getting-started/app/models"
)

type Blog struct{}


func (c Blog) CreateBlog(blog *models.Blog) error {
	return config.Global.DB.Insert(blog)
}
func (c Blog) CreateBlogs(blogs []*models.Blog) error {
	for _,blog := range blogs {
		if err := c.CreateBlog(blog); err != nil {
			return err
		}
	}
	return nil
}


func (c Blog) UpdateBlog(blog *models.Blog) error {
	return config.Global.DB.Update(blog)
}

func (c Blog) GetBlog(blog *models.Blog) error {
	return config.Global.DB.Select(&blog)
}
func (c Blog) GetBlogs(blogs []*models.Blog) error {
	return config.Global.DB.Model(&blogs).Select()
}


func (c Blog) GetAllBlogs(blogs *[]models.Blog) error {
	count, err := config.Global.DB.Model(&models.Blog{}).Count()
	if err != nil {
		return err
	}
	return config.Global.DB.Model(&blogs).Limit(count).Select()
}


func (c Blog) GetFullBlog(blog *models.Blog) error {
	return config.Global.DB.Model(&blog).Column("blog.*", "Gifs").Where("id = ?",blog.Id).First()
}
func (c Blog) GetFullBlogs(blogs []*models.Blog) error {
	return config.Global.DB.Model(&blogs).Column("blog.*", "Gifs").Select()
}


func (c Blog) GetBlogByTitle(title string, blog *models.Blog) error {
	return config.Global.DB.Model(&blog).Where("title = ?",title).First()
}

func (c Blog) GetBlogByUrl(url string, blog *models.Blog) error {
	return config.Global.DB.Model(&blog).Where("url = ?",url).First()
}


func (c Blog) DeleteBlog(blog *models.Blog) error {
	// TODO : delete cascade on blog delete
	return config.Global.DB.Delete(&blog)
}
func (c Blog) DeleteBlogs(blogs []*models.Blog) error {
	// TODO : delete cascade on blogs delete
	for _,blog := range blogs {
		if err := c.DeleteBlog(blog); err != nil {
			return err
		}
	}
	return nil
}