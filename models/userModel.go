package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type User struct {
	ID           primitive.ObjectID `bson:"_id"`
	firstName    *string            `json: "firstName" validate:"required, min=2, max=100"`
	lastName     *string            `json: "lastName" validate:"required, min=2, max=100"`
	password     *string            `json: "password" validate:"required, min=6"`
	email        *string            `json: "email" validate:"email, required"`
	phone        *string            `json: "phone" validate:"required"`
	token        *string            `json: "token"`
	userType     *string            `json: "userType" validate:"required, eq=ADMIN|eq=USER"`
	refreshToken *string            `json: "refreshToken"`
	cratedAt     *time.Time         `json: "cratedAt"`
	updatedAt    *time.Time         `json: "updatedAt"`
	userId       string             `json: "userId"`
}
