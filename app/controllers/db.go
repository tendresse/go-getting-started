package controllers

import (
	// "encoding/json"
	_ "gopkg.in/pg.v5"
	log "github.com/Sirupsen/logrus"
	"fmt"
	"github.com/tendresse/go-getting-started/app/models"
	"github.com/tendresse/go-getting-started/app/dao"
	"github.com/tendresse/go-getting-started/app/config"
)
type DB struct{}

func executeQueries(queries []string) error {
	for _, q := range queries {
		_, err := config.Global.DB.Exec(q)
		if err != nil {
			return err
		}
	}
	return nil
}

var create_schema = []string{
	`CREATE TABLE tags (
		id 		SERIAL 	CONSTRAINT pk_tag PRIMARY KEY,
		banned		boolean DEFAULT FALSE,
		title		text NOT NULL UNIQUE
	);`,
	`CREATE TABLE achievements (
		id 		SERIAL 	CONSTRAINT pk_achievement PRIMARY KEY,
		condition   	int     DEFAULT 10,
		icon        	text,
		tag_id		int	REFERENCES tags,
		title       	text,
		type_of      	text,
		xp          	int     DEFAULT 10
	);`,
	`CREATE TABLE blogs (
		id 		SERIAL 	CONSTRAINT pk_blog PRIMARY KEY,
		title        	text,
		url       	text NOT NULL UNIQUE
	);`,
	`CREATE TABLE gifs (
		id 		SERIAL 	CONSTRAINT pk_gif PRIMARY KEY,
		blog_id		int	REFERENCES blogs ON DELETE SET NULL,
		url		text    NOT NULL UNIQUE,
		lame_score	int     DEFAULT 0
	);`,
	`CREATE TABLE roles (
		id 		SERIAL 	CONSTRAINT pk_role PRIMARY KEY,
		title		text
	);
	`,
	`CREATE TABLE users (
		id 		SERIAL 	CONSTRAINT pk_user PRIMARY KEY,
		device 		text,
		nsfw		boolean	DEFAULT false,
		username 	text,
		passhash 	text,
		premium		boolean DEFAULT false
	);
	`,
	`CREATE TABLE tendresses (
		id 		SERIAL 	CONSTRAINT pk_tendresse PRIMARY KEY,
		sender_id 	int 	REFERENCES users ON DELETE CASCADE,
		receiver_id 	int	REFERENCES users ON DELETE CASCADE,
		gif_id 		int	REFERENCES gifs ON DELETE CASCADE,
		viewed		boolean	DEFAULT false
	);`,
	`CREATE TABLE gifs_tags (
		tag_id 		int	REFERENCES tags ON DELETE CASCADE,
		gif_id		int	REFERENCES gifs ON DELETE CASCADE,
		CONSTRAINT pk_gif_tags  PRIMARY KEY (tag_id, gif_id)
	);`,
	`CREATE TABLE users_achievements (
		achievement_id 	int	REFERENCES achievements ON DELETE CASCADE,
		user_id		int	REFERENCES users ON DELETE CASCADE,
		score        	int     DEFAULT 0,
		unlocked	boolean	DEFAULT false,
		CONSTRAINT pk_user_achievements PRIMARY KEY (achievement_id, user_id)
	);`,
	`CREATE TABLE users_friends (
		user_id 	int	REFERENCES users ON DELETE CASCADE,
		friend_id	int	REFERENCES users ON DELETE CASCADE,
		CONSTRAINT pk_user_friends PRIMARY KEY (user_id, friend_id)

	);`,
	`CREATE TABLE users_roles (
		role_id		int	REFERENCES roles ON DELETE CASCADE,
		user_id		int	REFERENCES users ON DELETE CASCADE,
		CONSTRAINT pk_user_roles PRIMARY KEY (role_id, user_id)
	);`,
}

var delete_schema = []string{
	`DROP SCHEMA public CASCADE;`,
	`CREATE SCHEMA public;`,
}


func (c DB) LoadDB(){
	if err := executeQueries(delete_schema); err != nil {
		fmt.Println("bug ?")
		log.Fatalln(err)
	}
	fmt.Println("deleted schema")

	if err := executeQueries(create_schema); err != nil {
		fmt.Println("bug ?")
		log.Fatalln(err)
	}
	fmt.Println("created schema")
	// TAGS
	tags_dao := dao.Tag{DB:config.Global.DB}

	tag_milf := models.Tag{
		Title: "milf",
	}
	tag_boobs := models.Tag{
		Title: "boobs",
	}
	tag_gif := models.Tag{
		Title: "gif",
		Banned: true,
	}
	tag_cowboy := models.Tag{
		Title: "cowboy",
	}
	if err := tags_dao.CreateTags([]*models.Tag{&tag_milf,&tag_cowboy,&tag_boobs,&tag_gif}); err != nil {
		fmt.Println("CreateTags")
		fmt.Println(err)
	}

	// ACHIEVEMENTS
	achievements_dao := dao.Achievement{DB:config.Global.DB}

	achievement_mommy := models.Achievement{
		Condition: 10,
		Icon: "ach1.png",
		TagID: 1,
		Title: "mommy",
		TypeOf: "send",
		XP: 100,
	}
	achievement_rider := models.Achievement{
		Condition: 10,
		Icon: "ach2.jpg",
		TagID: 4,
		Title: "rider cowboy",
		TypeOf: "receive",
		XP: 5,
	}
	if err := achievements_dao.CreateAchievements([]*models.Achievement{&achievement_mommy, &achievement_rider}); err != nil {
		fmt.Println("CreateAchievements")
		fmt.Println(err)
	}

	// BLOGS
	blog_dao := dao.Blog{DB : config.Global.DB}

	blog_milf := models.Blog{
		Title: "milf_tumblr",
		Url: "https://milf.tumblr.com",
	}
	blog_cowboy:= models.Blog{
		Title: "super_cowboy",
		Url: "https://super_cowboy.tumblr.com",
	}
	if err := blog_dao.CreateBlogs([]*models.Blog{&blog_milf,&blog_cowboy}); err != nil {
		fmt.Println("CreateBlog")
		fmt.Println(err)
	}

	// GIFS
	gifs_dao := dao.Gif{DB:config.Global.DB}

	gif_milf_pov := models.Gif{
		BlogID:1,
		Url:"https://sc1.tumblr.com/gif1",
		LameScore:0,
	}
	gif_gays := models.Gif{
		BlogID:2,
		Url:"https://sc1.tumblr.com/gif2",
		LameScore:0,
	}
	if err := gifs_dao.CreateGifs([]*models.Gif{&gif_milf_pov, &gif_gays}); err != nil {
		fmt.Println("CreateGifs")
		fmt.Println(err)
	}

	// RELATION GIF TAGS
	if err := gifs_dao.AddTagToGif(&tag_milf, &gif_milf_pov); err != nil {
		fmt.Println("AddTagsToGif")
		fmt.Println(err)
	}
	if err := gifs_dao.AddTagToGif(&tag_boobs, &gif_milf_pov); err != nil {
		fmt.Println("AddTagsToGif")
		fmt.Println(err)
	}
	if err := gifs_dao.AddTagToGif(&tag_cowboy, &gif_gays); err != nil {
		fmt.Println("AddTagToGif")
		fmt.Println(err)
	}

	// ROLES
	role_dao := dao.Role{DB:config.Global.DB}
	role_admin := models.Role{
		Title:"admin",
	}
	if err := role_dao.CreateRole(&role_admin); err != nil {
		fmt.Println("CreateRole")
		fmt.Println(err)
	}
	role_modo := models.Role{
		Title:"modo",
	}
	if err := role_dao.CreateRole(&role_modo); err != nil {
		fmt.Println("CreateRole")
		fmt.Println(err)
	}


	// USERS
	users_dao := dao.User{DB:config.Global.DB}
	user_alice := models.User{
		Device: "token_android_123",
		NSFW: true,
		Username: "alice",
		Passhash: "md5_32_123",
		Premium: false,
	}
	if err := users_dao.CreateUser(&user_alice); err != nil {
		fmt.Println("CreateUser")
		fmt.Println(err)
	}
	user_bob := models.User{
		Device: "token_ios_123",
		NSFW: true,
		Username: "bob",
		Passhash: "md5_32_123",
		Premium: false,
	}
	if err := users_dao.CreateUser(&user_bob); err != nil {
		fmt.Println("CreateUser")
		fmt.Println(err)
	}

	// FRIENDS
	if err := users_dao.AddFriendToUser(&user_bob,&user_alice); err != nil {
		fmt.Println("AddFriendToUser")
		fmt.Println(err)
	}
	if err := users_dao.AddFriendToUser(&user_alice,&user_bob); err != nil {
		fmt.Println("AddFriendToUser")
		fmt.Println(err)
	}

	// RELATIONS USERS ACHIEVEMENTS
	if err := users_dao.AddAchievementToUser(&achievement_mommy,&user_alice); err != nil {
		fmt.Println("AddAchievementsToUser")
		fmt.Println(err)
	}
	if err := users_dao.AddAchievementToUser(&achievement_rider,&user_alice); err != nil {
		fmt.Println("AddAchievementsToUser")
		fmt.Println(err)
	}


	// RELATION USERS ROLES
	if err := users_dao.AddRoleToUser(&role_admin,&user_alice); err != nil {
		fmt.Println("AddRoleToUser")
		fmt.Println(err)
	}
	if err := users_dao.AddRoleToUser(&role_modo,&user_bob); err != nil {
		fmt.Println("AddRoleToUser")
		fmt.Println(err)
	}

	// TENDRESSES
	tendresses_dao := dao.Tendresse{DB:config.Global.DB}
	alice_bob_1 := models.Tendresse{
		SenderID: user_alice.ID,
		ReceiverID: user_bob.ID,
		GifID: gif_gays.ID,
		Viewed: false,
	}
	alice_bob_2 := models.Tendresse{
		SenderID: user_alice.ID,
		ReceiverID: user_bob.ID,
		GifID: gif_milf_pov.ID,
		Viewed: true,
	}
	bob_alice_1 := models.Tendresse{
		SenderID: user_bob.ID,
		ReceiverID: user_alice.ID,
		GifID: gif_gays.ID,
		Viewed: false,
	}
	if err := tendresses_dao.CreateTendresse(&alice_bob_1); err != nil {
		fmt.Println("CreateTendresse")
		fmt.Println(err)
	}
	if err := tendresses_dao.CreateTendresse(&alice_bob_2); err != nil {
		fmt.Println("CreateTendresse")
		fmt.Println(err)
	}
	if err := tendresses_dao.CreateTendresse(&bob_alice_1); err != nil {
		fmt.Println("CreateTendresse")
		fmt.Println(err)
	}
}
