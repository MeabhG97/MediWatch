package routes

import (
	"mediwatch/server/controllers"

	"github.com/gin-gonic/gin"
)

func UserScheduleRouter(router *gin.Engine) {
	router.GET("/user/:id/schedule", controllers.GetAllUserSchedule)
	router.PUT("/user/:id/schedule", controllers.CreateUserSchedule)
	router.DELETE("/user/:id/schedule/:sId", controllers.DeleteUserSchedule)
}
