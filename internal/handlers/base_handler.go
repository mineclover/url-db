package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

type BaseHandler struct{}

func NewBaseHandler() *BaseHandler {
	return &BaseHandler{}
}

func (h *BaseHandler) HandleError(c *gin.Context, err error) {
	HandleError(c, err)
}

func (h *BaseHandler) ParseIntParam(c *gin.Context, param string) (int, error) {
	str := c.Param(param)
	if str == "" {
		return 0, NewValidationError("Missing parameter: "+param, nil)
	}
	
	value, err := strconv.Atoi(str)
	if err != nil {
		return 0, NewValidationError("Invalid parameter: "+param, nil)
	}
	
	return value, nil
}

func (h *BaseHandler) ParseIntQuery(c *gin.Context, query string, defaultValue int) int {
	str := c.Query(query)
	if str == "" {
		return defaultValue
	}
	
	value, err := strconv.Atoi(str)
	if err != nil {
		return defaultValue
	}
	
	return value
}

func (h *BaseHandler) GetStringQuery(c *gin.Context, query string) string {
	return c.Query(query)
}

func (h *BaseHandler) BindJSON(c *gin.Context, obj interface{}) error {
	if err := c.ShouldBindJSON(obj); err != nil {
		return NewValidationError("Invalid JSON format", err.Error())
	}
	return nil
}