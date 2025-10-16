package web

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"jike_gin/internal/domain"
	"jike_gin/internal/service"
	"net/http"
)

type ArticleHandler struct {
	svc service.ArticleService
}

func NewArticleHandler(svc service.ArticleService) *ArticleHandler {
	return &ArticleHandler{
		svc: svc,
	}
}

func (u *ArticleHandler) RegisterRoutes(server *gin.Engine) {
	articleGroup := server.Group("/articles")
	articleGroup.POST("/edit", u.Edit)
}
func (h *ArticleHandler) Edit(ctx *gin.Context) {
	type Req struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}

	userId := ctx.GetInt64("user_id")

	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}
	id, err := h.svc.Save(ctx, domain.Article{Title: req.Title, Content: req.Content, Author: domain.Author{Id: userId}})
	if err != nil {
		zap.L().Error("保存文章错误", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, Result{
			Code: 5001,
			Msg:  "保存错误",
			Data: "",
		})
		return
	}

	ctx.JSON(http.StatusOK, Result{Code: 0, Msg: "success", Data: id})

}
