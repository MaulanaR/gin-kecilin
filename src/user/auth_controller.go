package user

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/maulanar/gin-kecilin/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var valildator = validator.New()

func SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		user := User{}

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// validate input
		if err := valildator.Struct(user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// validate email is unique
		count, err := Collection().CountDocuments(ctx, bson.M{"email": user.Email})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Email already exists"})
			return
		}

		// encrypt password
		user.Password, err = utils.HashPassword(user.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// set param
		user.CreatedAt = time.Now()
		user.UpdatedAt = time.Now()
		user.ID = primitive.NewObjectID()
		user.UserID = user.ID.Hex()

		// insert to db
		_, err = Collection().InsertOne(ctx, user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		user.Password = nil
		resp := utils.Response{
			Status:     http.StatusText(http.StatusOK),
			Message:    "User created successfully",
			Data:       user,
			Pagination: utils.Pagination{},
		}
		c.JSON(http.StatusOK, resp.BuildSingleResponse())
	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		user := User{}
		FoundUser := User{}

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// validate user login
		err := Collection().FindOne(ctx, bson.M{"email": user.Email}).Decode(&FoundUser)
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Email not found, please check again!"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// validate password
		pwValid, err := utils.VerifyPassword(*user.Password, *FoundUser.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if !pwValid {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid password"})
			return
		}

		token, refreshToken, err := utils.GenerateToken(*FoundUser.Email, FoundUser.UserID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		FoundUser.Token = &token
		FoundUser.RefreshToken = &refreshToken

		// update token
		_, err = Collection().UpdateOne(ctx, bson.M{"email": user.Email}, bson.M{"$set": bson.M{"token": token, "refreshToken": refreshToken, "updated_at": time.Now()}})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message":       "User logged in successfully",
			"user":          FoundUser,
			"token":         FoundUser.Token,
			"refresh_token": FoundUser.RefreshToken,
		})
	}
}

func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		claims, _ := c.Get("claims")
		tokenClaim, ok := claims.(*utils.Claims)
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid token claims"})
			return
		}

		userID := tokenClaim.UserID
		log.Printf("Claims: %v\n", claims)
		var user User
		err := Collection().FindOne(ctx, bson.M{"user_id": userID}).Decode(&user)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization token"})
			return
		}
		user.Password = nil
		c.JSON(http.StatusOK, user)
	}
}
