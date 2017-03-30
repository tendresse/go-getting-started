package tests

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
