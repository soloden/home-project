package model

import "time"

type User struct {
	UUID      string    `bson:"uuid"`
	Email     string    `bson:"email"`
	Username  string    `bson:"username"`
	Password  string    `bson:"-"`
	Roles     []string  `bson:"roles"`
	CreatedAt time.Time `bson:"createdAt"`
	UpdatedAt time.Time `bson:"updatedAt"`
}
