// controllers/blogs_controllers.go

package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/tendresse/go-getting-started/app/dao"
	"github.com/tendresse/go-getting-started/app/config"
	"github.com/tendresse/go-getting-started/app/models"

	log "github.com/Sirupsen/logrus"
)


type Blog struct {
}


type errorString struct {
	s string
}
func (e *errorString) Error() string {
	return e.s
}

var http_client = &http.Client{Timeout: 10 * time.Second}

func fetchTumblr(url string, tumblr *models.Tumblr) (error) {
	resp, err := http_client.Get(url)
	defer resp.Body.Close()
	if err != nil {
		log.Error(err)
		return err
	}

	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
		return err
	}

	err = json.Unmarshal(contents, &tumblr)
	if err != nil {
		log.Error(err)
		return err
	}

	if tumblr.Meta.Status != 200 {
		log.Error(err)
		return &errorString{`{"success":false, "error":"API request is invalid"}`}
	}
	return nil
}

// wrap with admin rights
func (c Blog) GetBlogs() string {
	blogs_dao := dao.Blog{}
	blogs := []models.Blog{}
	if err := blogs_dao.GetAllBlogs(&blogs); err != nil {
		log.Error(err)
		return `{"success":false, "error":"no blog found"}`
	}
	b, err := json.Marshal(blogs)
	if err != nil {
		log.Error(err)
		return `{"success":false, "error":"marshal blogs json error"}`
	}
	return strings.Join([]string{`{"success":true, "blogs":`,string(b),"}"} ,"")
}

// wrap with admin rights
func (c Blog) GetBlog(blog_id int) string {
	blogs_dao := dao.Blog{}
	blog := models.Blog{Id: blog_id}
	if err := blogs_dao.GetFullBlog(&blog); err != nil {
		log.Error(err)
		return `{"success":false, "error":"blog not found"}`
	}
	b, err := json.Marshal(blog)
	if err != nil {
		log.Error(err)
		return `{"success":false, "error":"marshal blog json error"}`
	}
	return strings.Join([]string{`{"success":true, "blog":`,string(b),"}"} ,"")
}

// wrap with admin rights
func (c Blog) AddBlog(blog_url string) string {
	blogs_dao := dao.Blog{}
	tags_dao  := dao.Tag{}
	gifs_dao  := dao.Gif{}
	blog      := models.Blog{}
	if err := blogs_dao.GetBlogByUrl(blog_url, &blog); err != nil {
		log.Error(err)
		return `{"success":false, "error":"blog already exists"}`
	}
	blog = models.Blog{Url: blog_url}
	if err := blogs_dao.CreateBlog(&blog); err != nil {
		log.Error(err)
		return `{"success":false, "error":"error while creating the blog"}`
	}
	blog_url = strings.Join([]string{"https://api.tumblr.com/v2/blog/",blog_url,"/posts/photo?api_key=",config.Global.TumblrAPIKey},"")
	tumblr := models.Tumblr{}
	if err := fetchTumblr(blog_url, &tumblr); err != nil{
		log.Error(err)
		return `{"success":false, "error":"error while fetching the Tumblr"}`
	}

	blog_url = strings.Join([]string{blog_url,"&offset="},"")

	for i := 0; i < tumblr.Response.Blog.TotalPosts; i+=20 {
		posts_url := strings.Join([]string{blog_url,strconv.Itoa(i)},"")
		tumblr = models.Tumblr{}
		if err := fetchTumblr(posts_url, &tumblr); err != nil {
			log.Error(err)
			return `{"success":false, "error":"error while getting contents"}`
		}

		for _,post := range tumblr.Response.Posts {
			if strings.Compare(post.Type,"photo") != 0 {
				continue
			}
			gif_url := post.Photos[0].OriginalSize.Url
			if strings.Compare(gif_url[len(gif_url)-3:],"gif") != 0 {
				continue
			}
			gif := models.Gif{Url:gif_url}
			if err := gifs_dao.GetOrCreateGif(&gif); err != nil {
				log.Error(err)
			}
			gif.BlogId = blog.Id
			if err := gifs_dao.UpdateGif(&gif); err != nil {
				log.Error(err)
			}
			Tags:
			for _,tag_title := range post.Tags {
				tag := models.Tag{Title: tag_title}
				// check if tag is not too short
				if len(tag_title) > 2 {
					if err := tags_dao.GetOrCreateTag(&tag); err != nil {
						log.Error(err)
						continue Tags
					}
					if tag.Banned {
						continue Tags
					}
					if err := gifs_dao.AddTagToGif(&tag, &gif); err != nil {
						log.Error(err)
					}
				}
			}
		}
	}
	return `{"success":true}`
}

//admin
func (c Blog) DeleteBlog(blog_id int) string {
	blogs_dao := dao.Blog{}
	if err := blogs_dao.DeleteBlog(&models.Blog{Id: blog_id}); err != nil {
		log.Error(err)
		return `{"success":false, "error":"error while deleting blog"}`
	}
	return `{"success":true}`
}