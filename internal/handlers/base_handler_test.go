package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type TestStruct struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}

func setupBaseHandlerTest() (*gin.Engine, *BaseHandler) {
	gin.SetMode(gin.TestMode)

	handler := NewBaseHandler()
	router := gin.New()

	return router, handler
}

func TestBaseHandler_ParseIntParam_Success(t *testing.T) {
	router, handler := setupBaseHandlerTest()

	router.GET("/test/:id", func(c *gin.Context) {
		id, err := handler.ParseIntParam(c, "id")
		if err != nil {
			handler.HandleError(c, err)
			return
		}
		c.JSON(http.StatusOK, gin.H{"id": id})
	})

	req := httptest.NewRequest(http.MethodGet, "/test/123", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"id":123`)
}

func TestBaseHandler_ParseIntParam_MissingParam(t *testing.T) {
	router, handler := setupBaseHandlerTest()

	router.GET("/test", func(c *gin.Context) {
		id, err := handler.ParseIntParam(c, "id")
		if err != nil {
			handler.HandleError(c, err)
			return
		}
		c.JSON(http.StatusOK, gin.H{"id": id})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Missing parameter: id")
}

func TestBaseHandler_ParseIntParam_InvalidParam(t *testing.T) {
	router, handler := setupBaseHandlerTest()

	router.GET("/test/:id", func(c *gin.Context) {
		id, err := handler.ParseIntParam(c, "id")
		if err != nil {
			handler.HandleError(c, err)
			return
		}
		c.JSON(http.StatusOK, gin.H{"id": id})
	})

	req := httptest.NewRequest(http.MethodGet, "/test/invalid", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid parameter: id")
}

func TestBaseHandler_ParseIntQuery_WithValue(t *testing.T) {
	router, handler := setupBaseHandlerTest()

	router.GET("/test", func(c *gin.Context) {
		page := handler.ParseIntQuery(c, "page", 1)
		c.JSON(http.StatusOK, gin.H{"page": page})
	})

	req := httptest.NewRequest(http.MethodGet, "/test?page=5", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"page":5`)
}

func TestBaseHandler_ParseIntQuery_DefaultValue(t *testing.T) {
	router, handler := setupBaseHandlerTest()

	router.GET("/test", func(c *gin.Context) {
		page := handler.ParseIntQuery(c, "page", 1)
		c.JSON(http.StatusOK, gin.H{"page": page})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"page":1`)
}

func TestBaseHandler_ParseIntQuery_InvalidValue(t *testing.T) {
	router, handler := setupBaseHandlerTest()

	router.GET("/test", func(c *gin.Context) {
		page := handler.ParseIntQuery(c, "page", 1)
		c.JSON(http.StatusOK, gin.H{"page": page})
	})

	req := httptest.NewRequest(http.MethodGet, "/test?page=invalid", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"page":1`) // Should return default value
}

func TestBaseHandler_GetStringQuery(t *testing.T) {
	router, handler := setupBaseHandlerTest()

	router.GET("/test", func(c *gin.Context) {
		search := handler.GetStringQuery(c, "search")
		c.JSON(http.StatusOK, gin.H{"search": search})
	})

	req := httptest.NewRequest(http.MethodGet, "/test?search=example", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"search":"example"`)
}

func TestBaseHandler_GetStringQuery_EmptyValue(t *testing.T) {
	router, handler := setupBaseHandlerTest()

	router.GET("/test", func(c *gin.Context) {
		search := handler.GetStringQuery(c, "search")
		c.JSON(http.StatusOK, gin.H{"search": search})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"search":""`)
}

func TestBaseHandler_BindJSON_Success(t *testing.T) {
	router, handler := setupBaseHandlerTest()

	router.POST("/test", func(c *gin.Context) {
		var data TestStruct
		if err := handler.BindJSON(c, &data); err != nil {
			handler.HandleError(c, err)
			return
		}
		c.JSON(http.StatusOK, data)
	})

	reqBody := TestStruct{Name: "test", Value: 42}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/test", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"name":"test"`)
	assert.Contains(t, w.Body.String(), `"value":42`)
}

func TestBaseHandler_BindJSON_InvalidJSON(t *testing.T) {
	router, handler := setupBaseHandlerTest()

	router.POST("/test", func(c *gin.Context) {
		var data TestStruct
		if err := handler.BindJSON(c, &data); err != nil {
			handler.HandleError(c, err)
			return
		}
		c.JSON(http.StatusOK, data)
	})

	req := httptest.NewRequest(http.MethodPost, "/test", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid JSON format")
}

func TestBaseHandler_HandleError_ValidationError(t *testing.T) {
	router, handler := setupBaseHandlerTest()

	router.GET("/test", func(c *gin.Context) {
		err := NewValidationError("test validation error", "details")
		handler.HandleError(c, err)
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "test validation error")
}

func TestBaseHandler_HandleError_NotFoundError(t *testing.T) {
	router, handler := setupBaseHandlerTest()

	router.GET("/test", func(c *gin.Context) {
		err := NewNotFoundError("resource not found")
		handler.HandleError(c, err)
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "resource not found")
}

func TestBaseHandler_HandleError_UnknownError(t *testing.T) {
	router, handler := setupBaseHandlerTest()

	router.GET("/test", func(c *gin.Context) {
		err := errors.New("unknown error")
		handler.HandleError(c, err)
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "An unexpected error occurred")
}
