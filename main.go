package main

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/atomic"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"jike_gin/internal/repository"
	articleRepo "jike_gin/internal/repository/article"
	"jike_gin/internal/repository/cache"
	"jike_gin/internal/repository/dao"
	"jike_gin/internal/repository/dao/article"
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
	//InitViper()
	InitLogger()
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

	articleDao := article.NewGORMArticleDAO(db)
	articleRep := articleRepo.NewArticleRepository(articleDao, nil, nil, db)
	articleSvc := service.NewArticleService(articleRep)
	//articleSvc := service.NewArticleServiceV1(articleRepo)
	articleHandler := web.NewArticleHandler(articleSvc)
	articleHandler.RegisterRoutes(server)

	//server := gin.Default()
	server.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"error": "", "msg": "ping..."})
	})
	server.Run(":8080")
}

func initDb() *gorm.DB {
	dsn := viper.GetString("db.mysql.dsn")
	fmt.Println("dsn", dsn)
	db, err := gorm.Open(mysql.Open("root:root@tcp(gin-webook-mysql:3309)/webook"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	//CreateUser(db)
	return db
}

func initRedis() redis.Cmdable {
	//addr := viper.GetString("redis.addr")
	redisClient := redis.NewClient(&redis.Options{
		Addr: "gin-book-redis:10379",
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

var tplId atomic.String

func InitViper() {
	//配置的名字,没有后缀
	viper.SetConfigName("dev")
	// 配置文件类型
	viper.SetConfigType("yaml")

	// 配置文件的路径,当前工作目录下的 config的子目录,可以有读个路径,扫描多个路径,允许 main 函数的时候是在
	//jike 目录下,但是 golang 的 ide 的工具里面配置了允许路径是 jike/gin,所以找的还是 gin下的 config
	viper.AddConfigPath("./config") // k8s 上找不到文件加,所以先将初始化方法去掉
	//viper.AddConfigPath("$HOME/.appname")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

}

// InitViperV1 如何根据不同环境读取不同配置,使用参数
func InitViperV1() {
	//在go ide里面的运行配置里面程序实参增加  --config=config/dev.yaml
	cfile := pflag.String("config", "config/dev.yaml", "配置文件路径")
	pflag.Parse()
	viper.SetConfigFile(*cfile)
	// 监听文件变化
	viper.WatchConfig()
	// 设置默认值
	tplId.Store("123")
	viper.OnConfigChange(func(e fsnotify.Event) {
		// e 里面不含有变化内容,配置需要重新获取赋值
		fmt.Println("配置变化===>", e)
		// 从配置中读取
		tplId.Store(viper.Get("tpl.id").(string))
	})

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

}

func InitViperRemote() {
	viper.SetConfigType("yaml")
	// 端口和compose.yaml里面一样,webook 是和其他使用 etcd的区分开
	err := viper.AddRemoteProvider("etcd3", "127.0.0.1:12379", "gin")
	if err != nil {
		panic(err)
	}

}

func InitLogger() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	zap.ReplaceGlobals(logger)
	zap.L().Info("create logger success")
}
func InitLogger2() {
	myCore := MyCore{}
	logger := zap.New(myCore)
	zap.ReplaceGlobals(logger)
	zap.L().Info("create logger success")
}

type MyCore struct {
	zapcore.Core
}

// 全局替换
func (c MyCore) Write(entry zapcore.Entry, fields []zapcore.Field) error {
	for _, v := range fields {
		if v.Key == "phone" {
			phone := v.String
			v.String = phone[:3] + "****" + phone[7:]
		}
	}
	return c.Core.Write(entry, fields)
}
