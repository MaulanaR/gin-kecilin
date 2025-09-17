package cctv

import (
	"time"

	"github.com/maulanar/gin-kecilin/database"
	"github.com/maulanar/gin-kecilin/src/contact"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Cctv struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	CctvID    string             `json:"cctv_id"             bson:"cctv_id,omitempty"`
	ContactID string             `json:"contact_id"          validate:"required" bson:"contact_id,omitempty"`
	Name      string             `json:"name"                validate:"required" bson:"name,omitempty"`
	Location  *string            `json:"location"            bson:"location,omitempty"`
	IPAddress *string            `json:"ip_address"          bson:"ip_address,omitempty"`
	Brand     *string            `json:"brand"               bson:"brand,omitempty"`
	Model     *string            `json:"model"               bson:"model,omitempty"`
	Status    string             `json:"status"              validate:"required,oneof=online offline maintenance" bson:"status,omitempty"`
	CreatedAt time.Time          `json:"created_at"          bson:"created_at,omitempty"`
	UpdatedAt time.Time          `json:"updated_at"          bson:"updated_at,omitempty"`

	Contact *contact.Contact `json:"contact"`
}

// whitelist field can be sorted
var AllowedSortFields = map[string]bool{
	"cctv_id":    true,
	"contact_id": true,
	"ip_address": true,
	"name":       true,
	"created_at": true,
	"updated_at": true,
}

func Collection() *mongo.Collection {
	return database.OpenCollection("cctvs")
}
