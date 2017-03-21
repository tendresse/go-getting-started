package dao

import (
	"github.com/tendresse/go-getting-started/app/models"
	"gopkg.in/pg.v5"
)

type Blog struct{
	DB 	*pg.DB
}

func (c Blog) CreateBlog(blog *models.Blog) error {
	return c.DB.Insert(blog)
}

func (c Blog) CreateBlogs(blogs *[]models.Blog) error {
	return c.DB.Insert(blogs)
}

func (c Blog) UpdateBlog(blog *models.Blog) error {
	return c.DB.Update(blog)
}

func (c Blog) GetBlog(blog *models.Blog) error {
	return c.DB.Select(&blog)
}

func (c Blog) DeleteBlog(blog *models.Blog) error {
	// TODO : delete cascade on blog delete
	return c.DB.Delete(&blog)
}

func (c Blog) GetFullBlog(blog *models.Blog) error {
	return c.DB.Model(&blog).Column("blog.*", "Gifs").First()
}

func (c Blog) GetBlogByTitle(title string, blog *models.Blog) error {
	return c.DB.Model(&blog).Where("title = ?",title).Select()
}

func (c Blog) GetBlogs(blogs *[]models.Blog) error {
	return c.DB.Model(&blogs).Select()
}

func (c Blog) GetFullBlogs(blogs *[]models.Blog) error {
	return c.DB.Model(&blogs).Column("blog.*", "Gifs").Select()
}

func (c Blog) DeleteBlogs(blogs *[]models.Blog) error {
	// TODO : delete cascade on blogs delete
	return c.DB.Delete(&blogs)
}

