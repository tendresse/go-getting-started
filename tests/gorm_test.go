package tests

import (
	// "encoding/json"
	"gopkg.in/pg.v5"
	"log"
	"fmt"
	"github.com/tendresse/go-getting-started/app/models"
	"github.com/tendresse/go-getting-started/app/dao"
)

func executeQueries(queries []string) error {
	for _, q := range queries {
		_, err := Global.DB.Exec(q)
		if err != nil {
			return err
		}
	}
	return nil
}

type GlobalStruct struct{
	DB	*pg.DB
}

var Global GlobalStruct

func init() {
	initDB()
	loadDB()
	test()
	testAchievement()
	testBlog()
	testGif()
	testTag()
	testRole()
	testUser()
	//
}

func test(){
	// TAGS
	tags_dao := dao.Tag{DB:Global.DB}
	gifs_dao := dao.Gif{DB:Global.DB}

	tag_milf := models.Tag{
		Title: "milf",
		Banned: false,
	}
	tag_boobs := models.Tag{
		Title: "boobs",
		Banned: false,
	}
	tag_gif := models.Tag{
		Title: "gif",
		Banned: true,
	}
	tag_cowboy := models.Tag{
		Title: "cowboy",
		Banned: false,
	}
	if err := tags_dao.CreateTags(&[]models.Tag{tag_milf,tag_cowboy,tag_boobs,tag_gif}); err != nil {
		fmt.Println(err)
	}

	// ACHIEVEMENTS
	achievements_dao := dao.Achievement{DB:Global.DB}

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
	if err := achievements_dao.CreateAchievements(&[]models.Achievement{achievement_mommy, achievement_rider}); err != nil {
		fmt.Println(err)
	}

	// BLOGS
	blog_dao := dao.Blog{DB : Global.DB}

	blog_milf := models.Blog{
		Title: "milf_tumblr",
		Url: "https://milf.tumblr.com",
	}
	blog_cowboy:= models.Blog{
		Title: "super_cowboy",
		Url: "https://super_cowboy.tumblr.com",
	}
	if err := blog_dao.CreateBlogs(&[]models.Blog{blog_milf,blog_cowboy}); err != nil {
		fmt.Println("CreateBlog")
		fmt.Println(err)
	}

	// GIFS
	gifs_dao := dao.Gif{DB:Global.DB}

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
	if err := gifs_dao.CreateGifs(&[]models.Gif{gif_milf_pov,gif_gays}); err != nil {
		fmt.Println("CreateGifs")
		fmt.Println(err)
	}
	// RELATION GIF TAGS
	if err := gifs_dao.AddTagsToGif(&[]models.Tag{tag_milf,tag_boobs}, &gif_milf_pov); err != nil {
		fmt.Println("AddTagsToGif")
		fmt.Println(err)
	}
	if err := gifs_dao.AddTagToGif(&tag_cowboy, &gif_gays); err != nil {
		fmt.Println("AddTagToGif")
		fmt.Println(err)
	}

	// ROLES
	role_dao := dao.Role{DB:Global.DB}
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
	users_dao := dao.User{DB:Global.DB}
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

	// RELATIONS USERS ACHIEVEMENTS

	if err := users_dao.AddAchievementsToUser(&[]models.Achievement{achievement_mommy,achievement_rider}, &user_alice); err != nil {
		fmt.Println("AddAchievementsToUser")
		fmt.Println(err)
	}

}

func initDB() {
	Global.DB = pg.Connect(&pg.Options{
		User: "postgres",
		Database: "postgres",
		Password: "postgres",
		Addr: "127.0.0.1:32768",
	})
	if err := executeQueries(delete_schema); err != nil {
		log.Fatalln(err)
	}
	fmt.Println("deleted schema")

	if err := executeQueries(create_schema); err != nil {
		log.Fatalln(err)
	}
	fmt.Println("created schema")
}

func loadDB(){
	// TENDRESSES
	alice_bob_1 := models.Tendresse{
		SenderID: user_alice.ID,
		ReceiverID: user_bob.ID,
		GifID: gif_gays.ID,
		Viewed: false,
	}
	if err := Global.DB.Insert(&alice_bob_1); err != nil {
		fmt.Println(err)
	}
	alice_bob_2 := models.Tendresse{
		SenderID: user_alice.ID,
		ReceiverID: user_bob.ID,
		GifID: gif_milf_pov.ID,
		Viewed: true,
	}
	if err := Global.DB.Insert(&alice_bob_2); err != nil {
		fmt.Println(err)
	}
	bob_alice_1 := models.Tendresse{
		SenderID: user_bob.ID,
		ReceiverID: user_alice.ID,
		GifID: gif_gays.ID,
		Viewed: false,
	}
	if err := Global.DB.Insert(&bob_alice_1); err != nil {
		fmt.Println(err)
	}

	// FRIENDS
	relation_users_friends_1 := models.UsersFriends{
		UserID:user_alice.ID,
		FriendID:user_bob.ID,
	}
	if err := Global.DB.Insert(&relation_users_friends_1); err != nil {
		fmt.Println(err)
	}
	relation_users_friends_2 := models.UsersFriends{
		UserID:user_bob.ID,
		FriendID:user_alice.ID,
	}
	if err := Global.DB.Insert(&relation_users_friends_2); err != nil {
		fmt.Println(err)
	}

	// RELATION USERS ROLES
	relation_users_roles_1 := models.UsersRoles{
		UserID:user_alice.ID,
		RoleID:role_admin.ID,
	}
	if err := Global.DB.Insert(&relation_users_roles_1); err != nil {
		fmt.Println(err)
	}
	relation_users_roles_2 := models.UsersRoles{
		UserID:user_bob.ID,
		RoleID:role_modo.ID,
	}
	if err := Global.DB.Insert(&relation_users_roles_2); err != nil {
		fmt.Println(err)
	}
}

func testAchievement() {
	achievements_dao := dao.Achievement{DB:Global.DB}

	achievement_create := models.Achievement{
		Condition: 10,
		Icon: "ach5.jpg",
		TagID: 4,
		Title: "i like dolphin",
		TypeOf: "receive",
		XP: 10,
	}
	if err := achievements_dao.CreateAchievement(&achievement_create); err != nil {
		fmt.Println("CreateAchievement")
		fmt.Println(err)
	}

	achievement_create.Title = "i like trains"
	if err := achievements_dao.UpdateAchievement(&achievement_create); err != nil {
		fmt.Println("UpdateAchievement")
		fmt.Println(err)
	}

	achievement := models.Achievement{ID:1}
	if err := achievements_dao.GetAchievement(&achievement); err != nil {
		fmt.Println("GetAchievement")
		fmt.Println(err)
	}

	achievement_full := models.Achievement{ID:1}
	if err := achievements_dao.GetFullAchievement(&achievement_full); err != nil {
		fmt.Println("GetFullAchievement")
		fmt.Println(err)
	}

	achievement_by_title := models.Achievement{}
	if err := achievements_dao.GetAchievementByTitle("mommy", &achievement_by_title); err != nil {
		fmt.Println("GetAchievementByTitle")
		fmt.Println(err)
	}

	achievements := []models.Achievement{{ID:1}, {ID:2}}
	if err := achievements_dao.GetAchievements(&achievements); err != nil {
		fmt.Println("GetAchievements")
		fmt.Println(err)
	}

	achievements_full := []models.Achievement{{ID:1}, {ID:2}}
	if err := achievements_dao.GetFullAchievements(&achievements_full); err != nil {
		fmt.Println("GetFullAchievements")
		fmt.Println(err)
	}
}

func testBlog(){
	blog_dao := dao.Blog{DB : Global.DB}

	blog_create := models.Blog{Title:"test_blog", Url:"http://test.tumblr.com"}
	if err := blog_dao.CreateBlog(&blog_create); err != nil {
		fmt.Println("CreateBlog")
		fmt.Println(err)
	}

	blog_create.Title = "updated_name"
	if err := blog_dao.UpdateBlog(&blog_create); err != nil {
		fmt.Println("UpdateBlog")
		fmt.Println(err)
	}

	blog_get := models.Blog{ID:1}
	if err := blog_dao.GetBlog(&blog_get); err != nil {
		fmt.Println("GetBlog")
		fmt.Println(err)
	}


	blog_full := models.Blog{ID:2}
	if err := blog_dao.GetFullBlog(&blog_full); err != nil {
		fmt.Println("GetFullBlog")
		fmt.Println(err)
	}


	blog_title := models.Blog{}
	if err := blog_dao.GetBlogByTitle("updated_name", &blog_title); err != nil {
		fmt.Println("GetBlogByTitle")
		fmt.Println(err)
	}


	blogs := []models.Blog{{ID:1}, {ID:2}}
	if err := blog_dao.GetBlogs(&blogs); err != nil {
		fmt.Println("GetBlogs")
		fmt.Println(err)
	}

	blogs_full := []models.Blog{{ID:1}, {ID:2}}
	if err := blog_dao.GetFullBlogs(&blogs_full); err != nil {
		fmt.Println("GetFullBlogs")
		fmt.Println(err)
	}
}

func testGif(){
	gifs_dao := dao.Gif{DB:Global.DB}

	gif_create := models.Gif{
		BlogID:2,
		Url:"https://sc1.tumblr.com/gif14",
		LameScore:0,
	}
	if err := gifs_dao.CreateGif(&gif_create); err != nil {
		fmt.Println("CreateGif")
		fmt.Println(err)
	}
	fmt.Println("CreateGif")
	fmt.Println(gif_create)

	gif_create.LameScore = 10
	if err := gifs_dao.UpdateGif(&gif_create); err != nil {
		fmt.Println("UpdateGif")
		fmt.Println(err)
	}

	gif := models.Gif{ID:1}
	if err := gifs_dao.GetGif(&gif); err != nil {
		fmt.Println("GetGif")
		fmt.Println(err)
	}

	gif_full := models.Gif{ID:1}
	if err := gifs_dao.GetFullGif(&gif_full); err != nil {
		fmt.Println("GetFullGif")
		fmt.Println(err)
	}

	gifs := []models.Gif{{ID:1}, {ID:2}}
	if err := gifs_dao.GetGifs(&gifs); err != nil {
		fmt.Println("GetGifs")
		fmt.Println(err)
	}

	gifs_full := []models.Gif{{ID:1}, {ID:2}}
	if err := gifs_dao.GetFullGifs(&gifs_full); err != nil {
		fmt.Println("GetFullGifs")
		fmt.Println(err)
	}

}

func testTag(){
	tags_dao := dao.Tag{DB:Global.DB}
	gifs_dao := dao.Gif{DB:Global.DB}

	tag_create := models.Tag{
		Title: "bbc",
		Banned: false,
	}
	if err := tags_dao.CreateTag(&tag_create); err != nil {
		fmt.Println(err)
	}

	tag_create.Title = "i like trains"
	gif := models.Gif{ID:1}
	gifs_dao.GetGif(&gif)
	tag_create.Gifs = []models.Gif{gif}
	if err := tags_dao.UpdateTag(&tag_create); err != nil {
		fmt.Println(err)
	}

	Tag := models.Tag{ID:1}
	if err := tags_dao.GetTag(&Tag); err != nil {
		fmt.Println(err)
	}

	Tag_full := models.Tag{ID:1}
	if err := tags_dao.GetFullTag(&Tag_full); err != nil {
		fmt.Println(err)
	}

	Tag_by_title := models.Tag{}
	if err := tags_dao.GetTagByTitle("mommy", &Tag_by_title); err != nil {
		fmt.Println(err)
	}

	tags := []models.Tag{{ID:1}, {ID:2}}
	if err := tags_dao.GetTags(&tags); err != nil {
		fmt.Println(err)
	}

	tags_full := []models.Tag{{ID:1}, {ID:2}}
	if err := tags_dao.GetFullTags(&tags_full); err != nil {
		fmt.Println(err)
	}
}

func testRole(){
	roles_dao := dao.Role{DB : Global.DB}

	role_create := models.Role{Title:"admin"}
	roles_dao.CreateRole(&role_create)
	fmt.Println("CreateRole")
	fmt.Println(role_create)

	role_create.Title = "admin_updated"
	roles_dao.UpdateRole(&role_create)
	fmt.Println("UpdateRole")
	fmt.Println(role_create)

	role := models.Role{ID:1}
	roles_dao.GetRole(&role)
	fmt.Println("GetRole")
	fmt.Println(role)

	role_full := models.Role{ID:1}
	roles_dao.GetFullRole(&role_full)
	fmt.Println("GetFullRole")
	fmt.Println(role_full)

	roles := []models.Role{{ID:1}, {ID:2}}
	roles_dao.GetRoles(&roles)
	fmt.Println("GetRoles")
	fmt.Println(roles)

	roles_full := []models.Role{{ID:1}, {ID:2}}
	roles_dao.GetFullRoles(&roles_full)
	fmt.Println("GetFullRoles")
	fmt.Println(roles_full)
}

func testUser(){
	users_dao := dao.User{DB:Global.DB}

	user_create := models.User{
		Device: "token_android_123",
		NSFW: true,
		Username: "create",
		Passhash: "md5_32_123",
		Premium: false,
	}
	users_dao.CreateUser(&user_create)
	fmt.Println("CreateUser")
	fmt.Println(user_create)

	user_create.Username = "updated"
	users_dao.UpdateUser(&user_create)
	fmt.Println("UpdateUser")
	fmt.Println(user_create)

	user := models.User{ID:1}
	users_dao.GetUser(&user)
	fmt.Println("GetUser")
	fmt.Println(user)

	user_full := models.User{ID:1}
	users_dao.GetFullUser(&user_full)
	fmt.Println("GetFullUser")
	fmt.Println(user_full)

	users_dao.DeleteUser(&user_full)
	fmt.Println("DeleteUser")
	fmt.Println(user_full)

	users := []models.User{{ID:1},{ID:2}}
	users_dao.GetUsers(&users)
	fmt.Println("GetUsers")
	fmt.Println(users)

	users_full := []models.User{{ID:1},{ID:2}}
	users_dao.GetFullUsers(&users_full)
	fmt.Println("GetFullUsers")
	fmt.Println(users_full)

	users_dao.DeleteUsers(&users_full)
	fmt.Println("DeleteUsers")
	fmt.Println(users_full)

	username := "alice"
	user_alice := models.User{}
	users_dao.GetUserByUsername(&user_alice, username)
	fmt.Println("GetUserByUsername")
	fmt.Println(user_alice)


	role := models.Role{ID:2}
	users_dao.AddRoleToUser(&role, &user_alice)
	fmt.Println("AddRoleToUser")
	users_dao.GetFullUser(&user_alice)
	fmt.Println(&user_alice)

	role = models.Role{ID:2}
	users_dao.DeleteRoleFromUser(&role, &user_alice)
	fmt.Println("DeleteRoleFromUser")
	users_dao.GetFullUser(&user_alice)
	fmt.Println(&user_alice)

	roles := []models.Role{{ID:1},{ID:2}}
	users_dao.DeleteRolesFromUser(&roles, &user_alice)
	fmt.Println("DeleteRolesFromUser")
	users_dao.GetFullUser(&user_alice)
	fmt.Println(&user_alice)

	roles = []models.Role{{ID:1},{ID:2}}
	users_dao.AddRolesToUser(&roles, &user_alice)
	fmt.Println("AddRolesToUser")
	users_dao.GetFullUser(&user_alice)
	fmt.Println(&user_alice)

	friend := models.User{ID:2}
	users_dao.AddFriendToUser(friend *models.User, user *models.User)


	users_dao.DeleteFriendFromUser(friend *models.User, user *models.User)


	users_dao.AddFriendsToUser(friends *[]models.User, user *models.User)


	users_dao.DeleteFriendsFromUser(friends *[]models.User, user *models.User)

}