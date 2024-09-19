package route

import (
	"TestBackDev/handler"
	"github.com/gin-gonic/gin"
)

func RunAPI(address string) error {
	userHandler := handler.NewUserHandler()

	r := gin.Default()

	apiRoutes := r.Group("/api")
	userRoutes := apiRoutes.Group("/user")
	{
		userRoutes.POST("/register", userHandler.CreateUser)
		userRoutes.POST("/sign_in", userHandler.SignIn)
		userRoutes.POST("/refresh_tokens", userHandler.RefreshTokenPair)
	}

	return r.Run(address)
}
