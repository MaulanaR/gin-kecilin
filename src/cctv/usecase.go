package cctv

import (
	"context"
	"errors"
	"math"
	"time"

	"github.com/maulanar/gin-kecilin/src/contact"
	"github.com/maulanar/gin-kecilin/utils"

	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// adjustable depending on usecase
type UsecaseHandler struct {
	Ctx           context.Context
	Page          int64
	Limit         int64
	TotalData     int64
	FilterAndSort utils.HelperUsecaseHandler
}

var valildator = validator.New()

func (uc *UsecaseHandler) Get() ([]Cctv, error) {
	if uc.Page < 1 {
		uc.Page = 1
	}
	if uc.Limit < 1 {
		uc.Limit = 10
	}

	filter := uc.FilterAndSort.SetFilter() // dynamic filter by query param
	sort := uc.FilterAndSort.SetSort()     // dynamic sort by query param
	skip := (uc.Page - 1) * uc.Limit       // offset
	opts := options.Find().
		SetProjection(bson.M{ // block sensitive content
		}).
		SetSort(sort).
		SetSkip(skip).
		SetLimit(uc.Limit)

	// total docs
	total, err := Collection().CountDocuments(uc.Ctx, filter)
	if err != nil {
		return nil, err
	}

	cur, err := Collection().Find(uc.Ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cur.Close(uc.Ctx)

	var datas []Cctv
	if err := cur.All(uc.Ctx, &datas); err != nil {
		return nil, err
	}

	totalPages := int64(math.Ceil(float64(total) / float64(uc.Limit)))
	if totalPages > 0 && uc.Page > totalPages {
		datas = []Cctv{}
	}

	contactUC := contact.UsecaseHandler{
		Ctx: uc.Ctx,
	}
	for k := range datas {
		v := &datas[k]

		contact, err := contactUC.GetByID(v.ContactID)
		if err != nil {
			return nil, err
		}
		v.Contact = contact
		v.Contact.CCTVs = nil
	}

	uc.TotalData = total
	return datas, nil
}

func (uc *UsecaseHandler) Create(param *Cctv) error {
	// validate input
	if err := valildator.Struct(param); err != nil {
		return err
	}

	// validate ip_address is unique
	count, err := Collection().CountDocuments(uc.Ctx, bson.M{"ip_address": param.IPAddress})
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("Duplicate Ip Address")
	}

	// validate contact id is valid
	contactUC := contact.UsecaseHandler{
		Ctx: uc.Ctx,
	}
	_, err = contactUC.GetByID(param.ContactID)
	if err != nil {
		return err
	}

	param.ID = primitive.NewObjectID()
	param.CctvID = param.ID.Hex()
	param.CreatedAt = time.Now()
	param.UpdatedAt = time.Now()

	_, err = Collection().InsertOne(uc.Ctx, param)
	if err != nil {
		return err
	}

	return nil
}

func (uc *UsecaseHandler) GetByID(id string) (*Cctv, error) {
	var data Cctv
	err := Collection().FindOne(uc.Ctx, bson.M{"cctv_id": id}).Decode(&data)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("Data " + ModuleName + " with id " + id + " is not found")
		}
		return nil, err
	}

	// get data contact
	contactUC := contact.UsecaseHandler{
		Ctx: uc.Ctx,
	}
	data.Contact, err = contactUC.GetByID(data.ContactID)
	if err != nil {
		return nil, err
	}
	data.Contact.CCTVs = nil

	return &data, nil
}

func (uc *UsecaseHandler) UpdateByID(id string, param *Cctv) error {
	// validate id exists
	oldData, err := uc.GetByID(id)
	if err != nil {
		return err
	}

	param.ID = oldData.ID
	param.CctvID = oldData.CctvID
	param.UpdatedAt = time.Now()

	// validate ip address is unique
	if oldData.IPAddress != nil && *oldData.IPAddress != *param.IPAddress {
		count, err := Collection().CountDocuments(uc.Ctx, bson.M{"ip_address": param.IPAddress})
		if err != nil {
			return err
		}

		if count > 0 {
			return errors.New("IP Address already exists")
		}
	}

	// validate contact id is valid
	contactUC := contact.UsecaseHandler{
		Ctx: uc.Ctx,
	}
	_, err = contactUC.GetByID(param.ContactID)
	if err != nil {
		return err
	}

	filter := bson.M{"cctv_id": id}
	update := bson.M{"$set": param}
	_, err = Collection().UpdateOne(uc.Ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (uc *UsecaseHandler) DeleteByID(id string) error {
	// validate id exists
	_, err := uc.GetByID(id)
	if err != nil {
		return err
	}

	filter := bson.M{"cctv_id": id}
	_, err = Collection().DeleteOne(uc.Ctx, filter)
	if err != nil {
		return err
	}

	return nil
}
