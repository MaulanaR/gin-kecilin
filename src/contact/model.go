package contact

import (
	"time"

	"github.com/maulanar/gin-kecilin/database"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Contact struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	FirstName *string            `json:"first_name"              validate:"required,min=2,max=100" bson:"first_name,omitempty"`
	LastName  *string            `json:"last_name,omitempty"     validate:""                       bson:"last_name,omitempty"`
	Email     *string            `json:"email"                   validate:"required,email,min=2"   bson:"email,omitempty"`
	Phone     *string            `json:"phone,omitempty"         validate:""                       bson:"phone,omitempty"`
	Address   *string            `json:"address"                 validate:"required,min=2"         bson:"address,omitempty"`
	CreatedAt time.Time          `json:"created_at"              bson:"created_at,omitempty"`
	UpdatedAt time.Time          `json:"updated_at"              bson:"updated_at,omitempty"`
	ContactID string             `json:"contact_id"              bson:"contact_id,omitempty"`

	// relate to cctvs
	CCTVs []ContactCctv `json:"cctvs,omitempty"`
}

type ContactCctv struct {
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
	return database.OpenCollection("contacts")
}

func CctvCollection() *mongo.Collection {
	return database.OpenCollection("cctvs")
}
