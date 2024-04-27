package server

import (
	"net/http"

	"dxkite.cn/meownest/src/service"
	"github.com/gin-gonic/gin"
)

func NewEndpoint(s service.Endpoint) *Endpoint {
	return &Endpoint{s: s}
}

type Endpoint struct {
	s service.Endpoint
}

func (s *Endpoint) Create(c *gin.Context) {
	var param service.CreateEndpointParam

	if err := c.ShouldBind(&param); err != nil {
		Error(c, http.StatusBadRequest, "invalid_parameter", err.Error())
		return
	}

	rst, err := s.s.Create(c, &param)
	if err != nil {
		Error(c, http.StatusInternalServerError, "internal_error", err.Error())
		return
	}

	Result(c, http.StatusCreated, rst)
}

func (s *Endpoint) Get(c *gin.Context) {
	var param service.GetEndpointParam

	if err := c.ShouldBindUri(&param); err != nil {
		Error(c, http.StatusBadRequest, "invalid_parameter", err.Error())
		return
	}

	if err := c.ShouldBindQuery(&param); err != nil {
		Error(c, http.StatusBadRequest, "invalid_parameter", err.Error())
		return
	}

	rst, err := s.s.Get(c, &param)
	if err != nil {
		Error(c, http.StatusInternalServerError, "internal_error", err.Error())
		return
	}
	Result(c, http.StatusOK, rst)
}

func WithEndpoint(path string, server *Endpoint) func(s *HttpServer) {
	return func(s *HttpServer) {
		group := s.engine.Group(path)
		{
			group.POST("", server.Create)
			group.GET("/:id", server.Get)
		}
	}
}