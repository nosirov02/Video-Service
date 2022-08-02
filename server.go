package main

import (
	"awesomeProject/controller"
	"awesomeProject/middlewares"
	"awesomeProject/repository"
	"awesomeProject/service"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"os"
)

var (
	videoRepository repository.VideoRepository = repository.NewVideoRepository()
	videoService    service.VideoService       = service.NewVideoService(videoRepository)
	videoController controller.VideController  = controller.NewVideoController(videoService)
	jwtService      service.JwtService         = service.NewJwtService()
	loginService    service.LoginService       = service.NewLoginService()
	loginController controller.LoginController = controller.NewLoginController(loginService, jwtService)
)

func setupLogOutput() {
	f, _ := os.Create("gin.log")
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
}

func main() {
	server := gin.Default()

	//server.Use(gin.Recovery(), middlewares.Logger(), middlewares.BasicAuth(), gindump.Dump())
	//server.Use(gin.Recovery(), middlewares.Logger(), middlewares.BasicAuth())
	server.Use(gin.Recovery(), gin.Logger())

	server.Static("/css", "./templates/css")
	server.LoadHTMLGlob("./templates/*.html")

	server.POST("/login", func(ctx *gin.Context) {
		token := loginController.Login(ctx)
		if token != "" {
			ctx.JSON(http.StatusOK, gin.H{"token": token})
		} else {
			ctx.JSON(http.StatusUnauthorized, nil)
		}
	})

	apiRoutes := server.Group("/api", middlewares.AuthorizeJWT())
	{
		apiRoutes.GET("/videos", func(ctx *gin.Context) {
			ctx.JSON(200, videoController.FindAll())
		})

		apiRoutes.POST("/videos", func(ctx *gin.Context) {
			err := videoController.Save(ctx)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			} else {
				ctx.JSON(http.StatusOK, gin.H{"message": "Video input is Valid"})
			}
		})

		apiRoutes.PUT("/videos/:id", func(ctx *gin.Context) {
			err := videoController.Update(ctx)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			} else {
				ctx.JSON(http.StatusOK, gin.H{"message": "Video input is Valid"})
			}
		})

		apiRoutes.DELETE("/videos/:id", func(ctx *gin.Context) {
			err := videoController.Delete(ctx)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			} else {
				ctx.JSON(http.StatusOK, gin.H{"message": "Video input is Valid"})
			}
		})
	}
	viewRouts := server.Group("/view")
	{
		viewRouts.GET("/videos", videoController.ShowAll)
	}
	server.Run(":8080")
}
