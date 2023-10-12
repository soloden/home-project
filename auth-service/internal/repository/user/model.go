package user

import "time"

type User struct {
	UUID         string    `json:"uuid" bson:"uuid"`
	Email        string    `json:"email" bson:"email"`
	Username     string    `json:"username" bson:"username"`
	PasswordHash []byte    `json:"password" bson:"password"`
	Roles        []string  `json:"roles" bson:"roles"`
	CreatedAt    time.Time `json:"created_at" bson:"createdAt"`
	UpdatedAt    time.Time `json:"updated_at" bson:"updatedAt"`
}
