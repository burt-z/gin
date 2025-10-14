package web

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"time"
)

type jwtHandler struct {
	atKey []byte
	rtKey []byte
}

func NewJetHandler() jwtHandler {
	return jwtHandler{
		atKey: []byte("95osj3fUD7fo0mlYdDbncXz4VD2igvf0"),
		rtKey: []byte("95osj3fUD7fo0mlYdDbncXz4VD2igvf0"),
	}
}
func (h *jwtHandler) setLoginToken(ctx *gin.Context, uid int64) error {
	ssid := uuid.New().String()
	err := h.setJWTToken(ctx, uid, ssid)
	if err != nil {
		return err
	}
	err = h.setRefreshToken(ctx, uid, ssid)
	return err
}

func (h *jwtHandler) setJWTToken(ctx *gin.Context, uid int64, ssid string) error {
	claims := UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 30)),
		},
		UId:  uid,
		Ssid: ssid,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenStr, err := token.SignedString(h.rtKey)
	if err != nil {
		return err
	}
	ctx.Header("x-jwt-token", tokenStr)
	return nil
}

func (h *jwtHandler) setRefreshToken(ctx *gin.Context, uid int64, ssid string) error {
	claims := RefreshUserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 30)),
		},
		UId:  uid,
		Ssid: ssid,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenStr, err := token.SignedString(h.rtKey)
	if err != nil {
		return err
	}
	ctx.Header("x-refresh-token", tokenStr)
	return nil
}

type UserClaims struct {
	jwt.RegisteredClaims
	UId       int64  `json:"id"`
	UserAgent string `json:"user_agent"`
	Ssid      string `json:"ssid"`
}

type RefreshUserClaims struct {
	jwt.RegisteredClaims
	UId  int64  `json:"id"`
	Ssid string `json:"ssid"`
}
