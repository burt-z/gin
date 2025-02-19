package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"go_project/gin/internal/repository"
	"go_project/gin/internal/repository/dao"
	"go_project/gin/internal/service"
	"go_project/gin/internal/web"
	"go_project/gin/internal/web/middleware"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {

	db := initDb()

	u := initUser(db)

	server := initWebServer()

	u.RegisterRoutes(server)

	server.Run(":8080")
}

func initDb() *gorm.DB {
	db, err := gorm.Open(mysql.Open("root:root@tcp(localhost:13316)/webook"))
	if err != nil {
		panic(err)
	}
	CreateUser(db)
	return db
}

func initUser(db *gorm.DB) *web.UserHandler {
	return web.NewUserHandler(service.NewUserService(repository.NewUserRepository(dao.NewUserDao(db))))
}

func initWebServer() *gin.Engine {
	server := gin.Default()

	server.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"PUT", "GET", "POST", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true, //允许携带凭证
	}))
	//store := cookie.NewStore([]byte("secret"))
	// 多个参最大空闲连接数
	store, err := redis.NewStore(16, "tcp", "localhost:6379", "", []byte("eaba3041-e2aa-440b-9b5e-05dbab61630e"), []byte("eaba1db08e-1a0e-421e-b636-d5b98b7f78e0"))
	if err != nil {
		panic(err)
	}

	server.Use(sessions.Sessions("mysession", store))
	server.Use(middleware.NewLoginMiddlewareBuilder().IgnorePath("/users/login").IgnorePath("/users/signup").Build())
	return server
}

func CreateUser(db *gorm.DB) {
	err := dao.CreateTale(db)
	if err != nil {
		panic(err)
	}
}
