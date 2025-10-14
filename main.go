package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"jike_gin/config"
	"jike_gin/internal/repository"
	"jike_gin/internal/repository/cache"
	"jike_gin/internal/repository/dao"
	"jike_gin/internal/service"
	"jike_gin/internal/service/sms/memory"
	"jike_gin/internal/service/wechat"
	"jike_gin/internal/web"
	ijwt "jike_gin/internal/web/jwt"
	"jike_gin/internal/web/middleware"
	"net/http"
	"strings"
)

func main() {
	db := initDb()
	rdb := initRedis()
	//u := initUser(db, rdb)
	ud := dao.NewUserDao(db)
	uc := cache.NewUserCache(rdb)
	repo := repository.NewUserRepository(ud, uc)
	svc := service.NewUserService(repo)
	codeCache := cache.NewCodeCache(rdb)
	codeRepo := repository.NewCodeRepository(codeCache)
	smsService := memory.NewService()
	codeSvc := service.NewCodeService(codeRepo, smsService)
	redisJwt := ijwt.NewRedisJWTHandler(rdb)
	u := web.NewUserHandler(svc, codeSvc, redisJwt)

	server := initWebServer()
	u.RegisterRoutes(server)

	wechatService := wechat.NewService("", "", nil)
	oAuthHandker := web.NewOAuth2WechatHandler(wechatService, svc)
	oAuthHandker.RegisterRouter(server)

	//server := gin.Default()
	server.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"error": "", "msg": "ping..."})
	})
	server.Run(":8080")
}

func initDb() *gorm.DB {
	db, err := gorm.Open(mysql.Open(config.Config.Db.DSN), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	//CreateUser(db)
	return db
}

func initRedis() redis.Cmdable {
	redisClient := redis.NewClient(&redis.Options{
		Addr: config.Config.Redis.Addr,
	})
	return redisClient
}

//func initUser(db *gorm.DB, rdb redis.Cmdable) *web.UserHandler {
//	ud := dao.NewUserDao(db)
//	uc := cache.NewUserCache(rdb)
//	repo := repository.NewUserRepository(ud, uc)
//	svc := service.NewUserService(repo)
//	codeCache := cache.NewCodeCache(rdb)
//	codeRepo := repository.NewCodeRepository(codeCache)
//	//c, _ := sms.NewClient(common.NewCredential("", ""), "", profile.NewClientProfile())
//	//smsService := tencent.NewService(context.Background(), "", "", c)
//	smsService := memory.NewService()
//	codeSvc := service.NewCodeService(codeRepo, smsService)
//	return web.NewUserHandler(svc, codeSvc)
//}

func initWebServer() *gin.Engine {
	server := gin.Default()

	server.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:3000"},
		AllowMethods: []string{"PUT", "GET", "POST", "DELETE"},
		AllowHeaders: []string{"Origin", "Content-Type", "Authorization"},
		// 不加前端拿不到,后端返回的
		ExposeHeaders:    []string{"x-jwt-token"},
		AllowCredentials: true, //允许携带凭证
		AllowOriginFunc: func(origin string) bool {
			if strings.HasPrefix(origin, "http://localhost") {
				return true
			}
			return strings.Contains(origin, "live.webook.com")
		},
	}))
	//store := cookie.NewStore([]byte("secret"))
	// 多个参最大空闲连接数
	//store, err := redis.NewStore(16, "tcp", config.Config.Redis.Addr, "", "", []byte("eaba3041e2aa440b9b5e05dbab6163"), []byte("eaba1db08e1a0e421eb636d5b98b7f78"))
	//if err != nil {
	//	panic(err)
	//}

	//redisClient := redisDb.NewClient(&redisDb.Options{
	//	Addr: config.Config.Redis.Addr,
	//})
	//// 限流,1秒钟限制 100
	//server.Use(ratelimit.NewBuilder(redisClient, time.Second, 100).Build())

	//server.Use(sessions.Sessions("mysession", store))
	//server.Use(middleware.NewLoginMiddlewareBuilder().IgnorePath("/users/login").IgnorePath("/users/signup").Build())
	server.Use(middleware.NewLoginJWTMiddlewareBuilder().
		IgnorePath("/users/login_sms").
		IgnorePath("/users/login_sms/code/send").
		IgnorePath("/users/login").
		IgnorePath("/users/signup").
		IgnorePath("/users/refresh_token").
		IgnorePath("/oauth2/wechat/authurl").
		Build())

	return server
}

func CreateUser(db *gorm.DB) {
	err := dao.CreateTale(db)
	if err != nil {
		panic(err)
	}
}
