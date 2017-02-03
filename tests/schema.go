package tests

var create_schema = []string{
	`CREATE TABLE tags (
		id 		SERIAL 	CONSTRAINT pk_tag PRIMARY KEY,
		banned		boolean,
		name		text
	);`,
	`CREATE TABLE achievements (
		id 		SERIAL 	CONSTRAINT pk_achievement PRIMARY KEY,
		condition   	int,
		icon        	text,
		tag_id		int	REFERENCES tags,
		title       	text,
		type_of      	text,
		xp          	int
	);`,
	`CREATE TABLE blogs (
		id 		SERIAL 	CONSTRAINT pk_blog PRIMARY KEY,
		name        	text,
		url       	text
	);`,
	`CREATE TABLE gifs (
		id 		SERIAL 	CONSTRAINT pk_gif PRIMARY KEY,
		blog_id		int	REFERENCES blogs,
		url		text,
		lame_score	int
	);`,
	`CREATE TABLE roles (
		id 		SERIAL 	CONSTRAINT pk_role PRIMARY KEY,
		name		text
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
		sender_id 	int 	REFERENCES users,
		receiver_id 	int	REFERENCES users,
		gif_id 		int	REFERENCES gifs,
		viewed		boolean	DEFAULT false
	);`,
	`CREATE TABLE gifs_tags (
		tag_id 		int	REFERENCES tags,
		gif_id		int	REFERENCES gifs,
		CONSTRAINT pk_gif_tags PRIMARY KEY (tag_id, gif_id)
	);`,
	`CREATE TABLE users_achievements (
		achievement_id 	int	REFERENCES achievements,
		user_id		int	REFERENCES users,
		score        	int,
		CONSTRAINT pk_user_achievements PRIMARY KEY (achievement_id, user_id)
	);`,
	`CREATE TABLE users_friends (
		user_id 	int	REFERENCES users,
		friend_id	int	REFERENCES users,
		CONSTRAINT pk_user_friends PRIMARY KEY (user_id, friend_id)

	);`,
	`CREATE TABLE users_roles (
		role_id		int	REFERENCES roles,
		user_id		int	REFERENCES users,
		CONSTRAINT pk_user_roles PRIMARY KEY (role_id, user_id)
	);`,
}

var delete_schema = []string{
	`DROP SCHEMA public CASCADE;`,
	`CREATE SCHEMA public;`,
}
