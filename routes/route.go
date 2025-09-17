package routes

import (
	"net/http"

	"github.com/maulanar/gin-kecilin/middleware"
	"github.com/maulanar/gin-kecilin/src/cctv"
	"github.com/maulanar/gin-kecilin/src/contact"
	"github.com/maulanar/gin-kecilin/src/user"

	"github.com/gin-gonic/gin"
)

func SetRouter(r *gin.Engine) {
	// Server status
	r.GET("/api/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong, Server Online !")
	})

	r.POST("/api/signup", user.SignUp())
	r.POST("/api/login", user.Login())

	// This endpoint requires login first
	protec := r.Group("/")
	protec.Use(middleware.Authenticate())
	{
		protec.GET("/api/user/me", user.GetUser())

		// Users
		protec.GET("/api/users", user.GetHandler())
		protec.GET("/api/users/:id", user.GetByIDHandler())
		protec.POST("/api/users", user.SignUp())
		protec.PUT("/api/users/:id", user.UpdateHandler())
		protec.PATCH("/api/users/:id", user.UpdateHandler())
		protec.DELETE("/api/users/:id", user.DeleteHandler())

		// Contacts
		protec.GET("/api/contacts", contact.GetHandler())
		protec.GET("/api/contacts/:id", contact.GetByIDHandler())
		protec.POST("/api/contacts", contact.CreateHandler())
		protec.PUT("/api/contacts/:id", contact.UpdateHandler())
		protec.PATCH("/api/contacts/:id", contact.UpdateHandler())
		protec.DELETE("/api/contacts/:id", contact.DeleteHandler())

		// CCTVS
		protec.GET("/api/cctvs", cctv.GetHandler())
		protec.GET("/api/cctvs/:id", cctv.GetByIDHandler())
		protec.POST("/api/cctvs", cctv.CreateHandler())
		protec.PUT("/api/cctvs/:id", cctv.UpdateHandler())
		protec.PATCH("/api/cctvs/:id", cctv.UpdateHandler())
		protec.DELETE("/api/cctvs/:id", cctv.DeleteHandler())
	}
}
