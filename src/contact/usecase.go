package contact

import (
	"context"
	"errors"
	"math"
	"time"

	"github.com/maulanar/gin-kecilin/utils"

	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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

func (uc *UsecaseHandler) Get() ([]Contact, error) {
	if uc.Page < 1 {
		uc.Page = 1
	}
	if uc.Limit < 1 {
		uc.Limit = 10
	}

	filter := uc.FilterAndSort.SetFilter() // dynamic filter by query param
	sort := uc.FilterAndSort.SetSort()     // dynamic sort by query param
	skip := (uc.Page - 1) * uc.Limit       // offset

	// total docs
	total, err := Collection().CountDocuments(uc.Ctx, filter)
	if err != nil {
		return nil, err
	}

	// get related CCTV
	pipeline := mongo.Pipeline{
		bson.D{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "cctvs"},
			{Key: "localField", Value: "contact_id"},
			{Key: "foreignField", Value: "contact_id"},
			{Key: "as", Value: "cctvs"},
		}}},
	}
	// add filters
	pipeline = append(pipeline, bson.D{{Key: "$match", Value: filter}})
	// add sorting
	if len(sort) > 0 {
		pipeline = append(pipeline, bson.D{{Key: "$sort", Value: sort}})
	}
	// pagination
	pipeline = append(pipeline,
		bson.D{{Key: "$skip", Value: skip}},
		bson.D{{Key: "$limit", Value: uc.Limit}},
	)

	cur, err := Collection().Aggregate(uc.Ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cur.Close(uc.Ctx)

	var datas []Contact
	if err := cur.All(uc.Ctx, &datas); err != nil {
		return nil, err
	}

	totalPages := int64(math.Ceil(float64(total) / float64(uc.Limit)))
	if totalPages > 0 && uc.Page > totalPages {
		datas = []Contact{}
	}

	uc.TotalData = total
	return datas, nil
}

func (uc *UsecaseHandler) Create(param *Contact) error {
	// validate input
	if err := valildator.Struct(param); err != nil {
		return err
	}

	// validate email is unique
	count, err := Collection().CountDocuments(uc.Ctx, bson.M{"email": param.Email})
	if err != nil {
		return err
	}

	if count > 0 {
		return errors.New("Email already exists")
	}

	param.ID = primitive.NewObjectID()
	param.ContactID = param.ID.Hex()
	param.CreatedAt = time.Now()
	param.UpdatedAt = time.Now()

	_, err = Collection().InsertOne(uc.Ctx, param)
	if err != nil {
		return err
	}

	return nil
}

func (uc *UsecaseHandler) GetByID(id string) (*Contact, error) {
	var data []Contact

	// get related CCTV
	pipeline := mongo.Pipeline{
		bson.D{{Key: "$match", Value: bson.M{"contact_id": id}}},
		bson.D{{Key: "$limit", Value: 1}},
		bson.D{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "cctvs"},
			{Key: "localField", Value: "contact_id"},
			{Key: "foreignField", Value: "contact_id"},
			{Key: "as", Value: "cctvs"},
		}}},
	}

	dt, err := Collection().Aggregate(uc.Ctx, pipeline)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("Data " + ModuleName + " with id " + id + " is not found")
		}
		return nil, err
	}
	defer dt.Close(uc.Ctx)

	if err := dt.All(uc.Ctx, &data); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("Data " + ModuleName + " with id " + id + " is not found")
		}
		return nil, err
	}

	return &data[0], nil
}

func (uc *UsecaseHandler) UpdateByID(id string, param *Contact) error {
	// validate id exists
	oldData, err := uc.GetByID(id)
	if err != nil {
		return err
	}

	param.ID = oldData.ID
	param.ContactID = oldData.ContactID
	param.UpdatedAt = time.Now()

	// validate email is unique
	if oldData.Email != nil && *oldData.Email != *param.Email {
		count, err := Collection().CountDocuments(uc.Ctx, bson.M{"email": param.Email})
		if err != nil {
			return err
		}

		if count > 0 {
			return errors.New("Email already exists")
		}
	}

	filter := bson.M{"contact_id": id}
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

	filter := bson.M{"contact_id": id}
	_, err = Collection().DeleteOne(uc.Ctx, filter)
	if err != nil {
		return err
	}

	return nil
}
