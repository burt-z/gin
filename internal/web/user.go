package web

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"jike_gin/consts"
	"jike_gin/internal/domain"
	"jike_gin/internal/service"
	"net/http"
	"regexp"
	"time"
)

type UserHandler struct {
	svc         service.UserService
	emailRegexp *regexp.Regexp // 验证邮箱
	codeSvc     service.CodeService
}

func NewUserHandler(svc service.UserService, codeSvc service.CodeService) *UserHandler {
	return &UserHandler{
		svc:         svc,
		emailRegexp: regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`),
		codeSvc:     codeSvc,
	}
}

func (u *UserHandler) RegisterRoutes(server *gin.Engine) {
	userGroup := server.Group("/users")
	userGroup.POST("/signup", u.SingUp)
	//userGroup.POST("/login", u.Login)
	userGroup.POST("/login", u.LoginJWT)
	userGroup.GET("/profile", u.Profile)
	userGroup.POST("/edit", u.ProfileEdit)
	userGroup.POST("/login_sms/code/send", u.SendLoginSMSCode)
	userGroup.POST("/login_sms", u.LoginSMS)
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
		ctx.JSON(http.StatusOK, gin.H{"status": 200, "code": 50010, "msg": "邮箱格式错误", "email": p.Email})
		return
	}
	err = u.svc.Signup(ctx, domain.User{Email: p.Email, Password: p.Password})
	if err != nil {
		ctx.String(http.StatusOK, "系统异常")
		return
	}

	ctx.String(http.StatusOK, "注册成功")
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
	err = u.setJWTToken(ctx, member.Id)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{Code: 4, Msg: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, Result{Code: 200, Msg: "success"})
}

func (u *UserHandler) setJWTToken(ctx *gin.Context, Uid int64) error {
	userClaims := UserClaims{
		UId:              Uid,
		RegisteredClaims: jwt.RegisteredClaims{ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 30))},
		UserAgent:        ctx.Request.UserAgent(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, userClaims)
	tokenStr, err := token.SignedString([]byte(consts.GetAuthSecret()))
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"status": 200, "code": 50010, "msg": err.Error()})
		return err
	}
	ctx.Header("x-jwt-token", tokenStr)
	fmt.Println("member", Uid)
	return nil
}

func (u *UserHandler) Logout(ctx *gin.Context) {
	sess := sessions.Default(ctx)
	// 退出登录
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

const biz = "login"

func (u *UserHandler) SendLoginSMSCode(ctx *gin.Context) {
	type Req struct {
		Phone string `json:"phone"`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}
	// 是不是一个合法的手机号码
	// 考虑正则表达式
	if req.Phone == "" {
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "输入有误",
		})
		return
	}
	err := u.codeSvc.Send(ctx, biz, req.Phone)
	switch err {
	case nil:
		ctx.JSON(http.StatusOK, Result{
			Msg: "发送成功",
		})
	case service.ErrCodeSendTooMany:
		ctx.JSON(http.StatusOK, Result{
			Msg: "发送太频繁，请稍后再试",
		})
	default:
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
	}
}

func (u *UserHandler) LoginSMS(ctx *gin.Context) {
	type Req struct {
		Phone string `json:"phone"`
		Code  string `json:"code"`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}
	// 这边，可以加上各种校验
	err := u.codeSvc.Verify(ctx, biz, req.Phone, req.Code)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	// 查找用户,设置 token
	user, err := u.svc.FindOrCreate(ctx, req.Phone)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{5, "注册/登录失败"})
		return
	}
	err = u.setJWTToken(ctx, user.Id)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{5, "注册/登录失败"})
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Msg: "验证码校验通过",
	})
}

type Result struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}
