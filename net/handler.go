package net

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sean0427/micro-service-pratice-auth-domain/model"
)

type service interface {
	Login(ctx context.Context, params *model.LoginInfo) (*model.Authentication, error)
	Verify(ctx context.Context, params *model.Authentication) (bool, error)
	Refresh(ctx context.Context, params *model.Authentication) (*model.Authentication, error)
}

type handler struct {
	service service
}

func New(service service) *handler {
	return &handler{
		service: service}
}

func (h *handler) Verify(c *gin.Context) {
	value := c.Params.ByName("Authorization")

	token := strings.TrimPrefix(value, "Bearer ") // header Authorization bare
	params := model.Authentication{
		Token: token,
	}

	auth, err := h.service.Verify(c, &params)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}

	c.JSONP(http.StatusOK, auth)
}

func (h *handler) Login(c *gin.Context) {
	var params model.LoginInfo

	if err := c.ShouldBind(&params); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}

	auth, err := h.service.Login(c, &params)
	if err != nil || auth == nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}

	c.Header("token", auth.Token)
	c.JSONP(http.StatusOK, auth)
}

func (h *handler) Refresh(c *gin.Context) {
	value := c.Params.ByName("Authorization")

	token := strings.TrimPrefix(value, "Bearer ") // header Authorization bare
	params := model.Authentication{
		Token: token,
	}

	auth, err := h.service.Refresh(c, &params)
	if err != nil {
		// TODO return unAuth
		c.AbortWithError(http.StatusInternalServerError, err)
	}

	c.Header("token", auth.Token)

	c.JSONP(http.StatusOK, auth)
}

func (h *handler) Register(e *gin.Engine) {
	r := e.Group("")
	{
		r.POST("/verify", h.Verify)
		r.POST("/access_token", h.Login)
		r.POST("/refresh", h.Refresh)

	}
}
