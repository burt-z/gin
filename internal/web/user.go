package web

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go_project/gin/consts"
	"go_project/gin/internal/domain"
	"go_project/gin/internal/service"
	"net/http"
	"regexp"
	"time"
)

type UserHandler struct {
	svc         *service.UserService
	emailRegexp *regexp.Regexp // 验证邮箱
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{
		svc:         svc,
		emailRegexp: regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`),
	}
}

func (u *UserHandler) RegisterRoutes(server *gin.Engine) {
	userGroup := server.Group("/users")
	userGroup.POST("/signup", u.SingUp)
	//userGroup.POST("/login", u.Login)
	userGroup.POST("/login", u.LoginJWT)
	userGroup.GET("/profile", u.Profile)
	userGroup.POST("/edit", u.ProfileEdit)
}

func (u *UserHandler) SingUp(ctx *gin.Context) {
	type Param struct {
		Email           string `json:"email"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirmPassword"`
	}
	var p Param
	err := ctx.ShouldBind(&p)
	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		ctx.JSON(http.StatusOK, gin.H{"status": 200, "code": 50010, "msg": "参数错误"})
		return
	}
	isSafe := u.emailRegexp.MatchString(p.Email)
	if !isSafe {
		ctx.JSON(http.StatusOK, gin.H{"status": 200, "code": 50010, "msg": "邮箱格式错误"})
		return
	}
	err = u.svc.Signup(ctx, domain.User{Email: p.Email, Password: p.Password})
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"status": 200, "code": 50010, "msg": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": 200, "code": 0, "msg": ""})
}

func (u *UserHandler) Login(ctx *gin.Context) {
	type Param struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var p Param
	err := ctx.ShouldBind(&p)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"status": 200, "code": 50010, "msg": err.Error()})
		return
	}
	isSafe := u.emailRegexp.MatchString(p.Email)
	if !isSafe {
		ctx.JSON(http.StatusOK, gin.H{"status": 200, "code": 50010, "msg": "邮箱格式错误"})
		return
	}
	member, err := u.svc.Login(ctx, domain.User{Email: p.Email, Password: p.Password})
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"status": 200, "code": 50010, "msg": err.Error()})
		return
	}

	// 登录成功设置 session
	sess := sessions.Default(ctx)
	sess.Set("userId", member.Id)
	sess.Options(sessions.Options{
		//HttpOnly: true,
		//Secure:   true,
		MaxAge: 60, //单位秒
	})
	sess.Save()

	ctx.JSON(http.StatusOK, gin.H{"status": 200, "code": 0, "msg": ""})
}

// LoginJWT
func (u *UserHandler) LoginJWT(ctx *gin.Context) {
	type Param struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var p Param
	err := ctx.ShouldBind(&p)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"status": 200, "code": 50010, "msg": err.Error()})
		return
	}
	isSafe := u.emailRegexp.MatchString(p.Email)
	if !isSafe {
		ctx.JSON(http.StatusOK, gin.H{"status": 200, "code": 50010, "msg": "邮箱格式错误"})
		return
	}
	member, err := u.svc.Login(ctx, domain.User{Email: p.Email, Password: p.Password})
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"status": 200, "code": 50010, "msg": err.Error()})
		return
	}
	userClaims := UserClaims{
		UId:              member.Id,
		RegisteredClaims: jwt.RegisteredClaims{ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 2))},
		UserAgent:        ctx.Request.UserAgent(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, userClaims)
	tokenStr, err := token.SignedString([]byte(consts.GetAuthSecret()))
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"status": 200, "code": 50010, "msg": err.Error()})
		return
	}
	ctx.Header("x-jwt-token", tokenStr)
	fmt.Println("member", member.Id)

	ctx.JSON(http.StatusOK, gin.H{"status": 200, "code": 0, "msg": ""})
}

func (u *UserHandler) Logout(ctx *gin.Context) {
	sess := sessions.Default(ctx)
	sess.Options(sessions.Options{MaxAge: -1})
	sess.Save()
	ctx.JSON(http.StatusOK, gin.H{"status": 200, "code": 0, "msg": "登出成功"})
}

func (u *UserHandler) ProfileEdit(ctx *gin.Context) {
	type Param struct {
		Birthday *string `json:"birthday"`
		NickName *string `json:"nickName"`
		AboutMe  *string `json:"about_me"`
	}
	var p Param
	err := ctx.ShouldBind(&p)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"status": 200, "code": 50010, "msg": err.Error()})
		return
	}

	userId := ctx.GetInt64("user_id")
	if userId == 0 {
		ctx.JSON(http.StatusOK, gin.H{"status": 200, "code": 50010, "msg": "用户未登录"})
		return
	}

	user := domain.User{Id: userId}
	if p.Birthday != nil {
		user.Birthday = *p.Birthday
		user.Keys = append(user.Keys, "birthday")
	}

	if p.NickName != nil {
		user.NickName = *p.NickName
		user.Keys = append(user.Keys, "nickname")
	}

	if p.AboutMe != nil {
		user.AboutMe = *p.AboutMe
		user.Keys = append(user.Keys, "about_me")
	}

	err = u.svc.ProfileEdit(ctx, user)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"status": 200, "code": 50010, "msg": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": 200, "code": 0})
}

func (u *UserHandler) Profile(ctx *gin.Context) {
	userId := ctx.GetInt64("user_id")
	if userId == 0 {
		ctx.JSON(http.StatusOK, gin.H{"status": 200, "code": 50010, "msg": "用户未登录"})
		return
	}
	user, err := u.svc.Profile(ctx, userId)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"status": 200, "code": 50010, "msg": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": 200, "Id": user.Id, "Email": user.Email, "Birthday": user.Birthday, "Nickname": user.NickName, "AboutMe": user.AboutMe, "data": map[string]interface{}{"code": 0, "msg": "success"}})
}

type UserClaims struct {
	jwt.RegisteredClaims
	UId       int64  `json:"id"`
	UserAgent string `json:"user_agent"`
}
