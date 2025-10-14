package web

import (
	"github.com/gin-gonic/gin"
	"jike_gin/internal/service"
	"jike_gin/internal/service/wechat"
	"net/http"
)

type OAuth2WechatHandler struct {
	svc     wechat.Service
	userSvc service.UserService
}

func NewOAuth2WechatHandler(s wechat.Service, u service.UserService) *OAuth2WechatHandler {
	return &OAuth2WechatHandler{
		svc:     s,
		userSvc: u,
	}
}

func (h *OAuth2WechatHandler) RegisterRouter(s *gin.Engine) {
	g := s.Group("/oauth2/wechat")
	g.GET("authurl", h.AuthURL)
	g.GET("callback", h.CallBack)
}

func (h *OAuth2WechatHandler) AuthURL(ctx *gin.Context) {
	url, err := h.svc.AuthURL(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, Result{Code: 7, Msg: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, Result{Data: url})
}

func (h *OAuth2WechatHandler) CallBack(ctx *gin.Context) {

}
