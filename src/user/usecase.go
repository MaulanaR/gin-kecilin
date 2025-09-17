package user

import (
	"context"
	"errors"
	"math"
	"time"

	"github.com/maulanar/gin-kecilin/utils"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// adjustable depending on usecase
type UsecaseHandler struct {
	GinCtx        *gin.Context
	Ctx           context.Context
	Page          int64
	Limit         int64
	TotalData     int64
	FilterAndSort utils.HelperUsecaseHandler
}

func (uc *UsecaseHandler) Get() ([]User, error) {
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
			"password":      0,
			"refresh_token": 0,
			"token":         0,
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

	var users []User
	if err := cur.All(uc.Ctx, &users); err != nil {
		return nil, err
	}

	totalPages := int64(math.Ceil(float64(total) / float64(uc.Limit)))
	if totalPages > 0 && uc.Page > totalPages {
		users = []User{}
	}

	uc.TotalData = total
	return users, nil
}

func (uc *UsecaseHandler) GetByID(id string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var user User
	err := Collection().FindOne(ctx, bson.M{"user_id": id}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("Data " + ModuleName + " with id " + id + " is not found")
		}
		return nil, err
	}

	return &user, nil
}

func (uc *UsecaseHandler) UpdateByID(id string, param *User) error {
	// validate id exists
	oldData, err := uc.GetByID(id)
	if err != nil {
		return err
	}

	param.ID = oldData.ID
	param.UserID = oldData.UserID
	param.UpdatedAt = time.Now()

	// if email changed
	if oldData.Email != nil && *oldData.Email != *param.Email {
		count, err := Collection().CountDocuments(uc.Ctx, bson.M{"email": param.Email})
		if err != nil {
			return err
		}

		if count > 0 {
			return errors.New("Email already exists")
		}
	}

	// if password changed
	if param.Password != nil {
		// encrypt password
		param.Password, err = utils.HashPassword(param.Password)
		if err != nil {
			return err
		}

	}

	filter := bson.M{"user_id": id}
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

	claims, _ := uc.GinCtx.Get("claims")
	tokenClaim, ok := claims.(*utils.Claims)
	if !ok {
		return errors.New("Invalid token claims")
	}

	userID := tokenClaim.UserID

	if id == userID {
		return errors.New("Cannot delete your own account")
	}

	filter := bson.M{"user_id": id}
	_, err = Collection().DeleteOne(uc.Ctx, filter)
	if err != nil {
		return err
	}

	return nil
}
