-- +migrate Up
CREATE TABLE "users" (
	"id" char(255) NOT NULL UNIQUE,
	"full_name" char(255) NOT NULL,
	"email" varchar(255) NOT NULL UNIQUE,
	"password" varchar(255) NOT NULL,
	"created_at" datetime DEFAULT CURRENT_TIMESTAMP,
	"updated_at" datetime DEFAULT CURRENT_TIMESTAMP,
	"deleted_at" datetime,
	PRIMARY KEY("email")
);

-- +migrate Down
DROP TABLE "users";