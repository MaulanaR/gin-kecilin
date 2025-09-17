package contact

import (
	"context"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/maulanar/gin-kecilin/utils"

	"github.com/gin-gonic/gin"
)

var ModuleName = "Contact"

func GetHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		page, _ := strconv.ParseInt(c.DefaultQuery("page", "1"), 10, 64)
		limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "10"), 10, 64)

		if limit < 1 {
			limit = 10
		}
		if limit > 200 {
			limit = 200
		}
		if page < 1 {
			page = 1
		}

		ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
		defer cancel()

		filters := map[string][]string{}
		for key, values := range c.Request.URL.Query() {
			if key == "page" || key == "limit" || key == "order_by" {
				continue
			}
			filters[key] = values
		}

		uc := UsecaseHandler{
			Ctx:   ctx,
			Page:  page,
			Limit: limit,
			FilterAndSort: utils.HelperUsecaseHandler{
				Filters:           filters,
				Sort:              c.Query("order_by"),
				AllowedSortFields: AllowedSortFields,
			},
		}

		datas, err := uc.Get()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		totalPages := int(math.Ceil(float64(uc.TotalData) / float64(limit)))

		resp := utils.Response{
			Status:  http.StatusText(http.StatusOK),
			Message: "Successfully get all " + ModuleName,
			Data:    datas,
			Pagination: utils.Pagination{
				Page:       int(page),
				Limit:      int(limit),
				TotalCount: int(uc.TotalData),
				TotalPages: totalPages,
				HasNext:    int(page) < totalPages,
				HasPrev:    page > 1,
			},
		}
		c.JSON(http.StatusOK, resp.BuildResponse())
	}
}

func GetByIDHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
		defer cancel()

		uc := UsecaseHandler{
			Ctx: ctx,
		}

		data, err := uc.GetByID(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		resp := utils.Response{
			Status:     http.StatusText(http.StatusOK),
			Message:    "Successfully get " + ModuleName,
			Data:       data,
			Pagination: utils.Pagination{},
		}
		c.JSON(http.StatusOK, resp.BuildSingleResponse())
	}
}

func CreateHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
		defer cancel()

		uc := UsecaseHandler{
			Ctx: ctx,
		}

		param := Contact{}

		if err := c.BindJSON(&param); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := uc.Create(&param)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		resp := utils.Response{
			Status:     http.StatusText(http.StatusOK),
			Message:    ModuleName + " created successfully",
			Data:       param,
			Pagination: utils.Pagination{},
		}
		c.JSON(http.StatusOK, resp.BuildSingleResponse())
	}
}

func UpdateHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
		defer cancel()

		uc := UsecaseHandler{
			Ctx: ctx,
		}

		param := Contact{}
		if err := c.BindJSON(&param); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := uc.UpdateByID(id, &param)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		resp := utils.Response{
			Status:     http.StatusText(http.StatusOK),
			Message:    ModuleName + " updated successfully",
			Data:       param,
			Pagination: utils.Pagination{},
		}
		c.JSON(http.StatusOK, resp.BuildSingleResponse())
	}
}

func DeleteHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
		defer cancel()

		uc := UsecaseHandler{
			Ctx: ctx,
		}

		err := uc.DeleteByID(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		resp := utils.Response{
			Status:     http.StatusText(http.StatusOK),
			Message:    ModuleName + " deleted successfully",
			Pagination: utils.Pagination{},
		}
		c.JSON(http.StatusOK, resp.BuildSingleResponse())
	}
}
