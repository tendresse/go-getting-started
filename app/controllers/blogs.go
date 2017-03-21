// controllers/blogs_controllers.go

package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
	"os"

	"github.com/tendresse/go-getting-started/app/config"
	"github.com/tendresse/go-getting-started/app/models"

	log "github.com/Sirupsen/logrus"
	_ "github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)


type BlogsController struct {
}

type errorString struct {
	s string
}
func (e *errorString) Error() string {
	return e.s
}

var http_client = &http.Client{Timeout: 10 * time.Second}

func fetchTumblr(url string) (models.Tumblr,error) {
	resp, err := http_client.Get(url)
	defer resp.Body.Close()
	if err != nil {
		log.Error(err)
		return nil, err
	}

	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	tumblr := models.Tumblr{}
	err = json.Unmarshal(contents, &tumblr)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	if tumblr.Meta.Status != 200 {
		log.Error(err)
		return nil, &errorString{`{"success":false, "error":"API request is invalid"}`}
	}
	return tumblr, nil
}

// wrap with admin rights
func (c BlogsController) GetBlogs() string {
	// TODO : CHOOSE WHAT JSON TO RETURN
	// on pourrait use un select ou un omit sur la query
	blogs := []models.Blog{}
	if err := config.Global.DB.Find(&blogs).Error; err != nil {
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
func (c BlogsController) GetBlog(blog_id int) string {
	// TODO : CHOOSE WHAT JSON TO RETURN
	blog := models.Blog{}
	if err := config.Global.DB.First(&blog, blog_id).Error; err != nil {
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
func (c BlogsController) AddBlog(blog_url string) string {
    // regex, user input data validation

	blog := models.Blog{}
	if config.Global.DB.Where("Url = ?",blog_url).First(&blog).Error != nil {
		blog = models.Blog{Url: blog_url}
		config.Global.DB.Create(&blog)
	} else {
		return `{"success":false, "error":"blog already exists"}`
	}

	known_tags := []models.Tag{}
	banned_tags := []models.Tag{}
	if config.Global.DB.Where("Banned = ?","true").Find(&banned_tags).Error != nil {
		banned_tags = []models.Tag{}
	}

	blog_url = strings.Join([]string{"https://api.tumblr.com/v2/blog/",blog_url,"/posts/photo?api_key=",config.Global.TumblrAPIKey},"")
	tumblr,err := fetchTumblr(blog_url)
	if err != nil{
		log.Error(err)
		return `{"success":false, "error":"error while fetching the Tumblr"}`
	}

	blog_url = strings.Join([]string{blog_url,"&offset="},"")

	for i := 0; i < tumblr.Response.Blog.TotalPosts; i+=20 {
		posts_url := strings.Join([]string{blog_url,strconv.Itoa(i)},"")
		
		tumblr_posts, err := fetchTumblr(posts_url)
		if err != nil {
			log.Error(err)
			return `{"success":false, "error":"error while getting contents"}`
		}

		for _,post := range tumblr_posts.Response.Posts {
			if strings.Compare(post.Type,"photo") != 0 {
				continue
			}
			gif_url := post.Photos[0].OriginalSize.Url
			if strings.Compare(gif_url[len(gif_url)-3:],"gif") != 0 {
				continue
			}
			gif := models.Gif{}
			config.Global.DB.Where(models.Gif{Url: gif_url}).FirstOrCreate(&gif)
			config.Global.DB.Model(&gif).Association("Blog").Append(blog)
			for _,post_tags := range post.Tags {
				Tags:
				for _,tag_title := range strings.Split(post_tags, " ") {
					tag := models.Tag{Name: tag_title}
                    			// check if tag is not too short
					if len(tag_title) < 2{
						continue 
					}
					for _,banned_tag := range banned_tags {
						if banned_tag.Name == tag.Name {
							continue Tags
						}
					}
					for _,known_tag := range known_tags {
						if known_tag.Name == tag.Name {
							config.Global.DB.Model(&gif).Association("Tags").Append(known_tag)
							config.Global.DB.Model(&gif).Association("Blog").Append(blog)
							continue Tags
						}
					}
					config.Global.DB.Create(&tag)
					append(known_tags,tag)
					config.Global.DB.Model(&gif).Association("Tags").Append(tag)
				}
			}
		}
	}
	return `{"success":true}`
}

//admin
func (c BlogsController) DeleteBlog(blog_id int) string {
	blog := models.Blog{}
	if err := config.Global.DB.First(&blog, blog_id).Error; err != nil {
		log.Error(err)
		return `{"success":false, "error":"blog not found"}`
	}
	config.Global.DB.Delete(&blog)
	return `{"success":true}`
}