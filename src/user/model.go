package user

import (
	"time"

	"github.com/maulanar/gin-kecilin/database"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type User struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	FirstName    *string            `json:"first_name"              validate:"required,min=2,max=100" bson:"first_name,omitempty"`
	LastName     *string            `json:"last_name,omitempty"     validate:""                       bson:"last_name,omitempty"`
	Email        *string            `json:"email"                   validate:"required,email,min=2"   bson:"email,omitempty"`
	Password     *string            `json:"password"                validate:"required,min=2,max=100" bson:"password,omitempty"`
	Phone        *string            `json:"phone,omitempty"         validate:""                       bson:"phone,omitempty"`
	Token        *string            `json:"token,omitempty"         validate:""                       bson:"token,omitempty"`
	RefreshToken *string            `json:"refresh_token,omitempty" validate:""                       bson:"refresh_token,omitempty"`
	CreatedAt    time.Time          `json:"created_at"              bson:"created_at,omitempty"`
	UpdatedAt    time.Time          `json:"updated_at"              bson:"updated_at,omitempty"`
	UserID       string             `json:"user_id"                 bson:"user_id,omitempty"`
}

// whitelist field can be sorted
var AllowedSortFields = map[string]bool{
	"first_name": true,
	"last_name":  true,
	"email":      true,
	"created_at": true,
	"updated_at": true,
}

func Collection() *mongo.Collection {
	return database.OpenCollection("users")
}
