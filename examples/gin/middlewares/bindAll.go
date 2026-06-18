package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New(validator.WithRequiredStructEnabled())

func init() {
	validate.SetTagName("binding")
	// Disable Gin's default validation since we're doing it manually after binding all sources
	binding.Validator = nil
}

func Bind[T any](c *gin.Context) {
	var req T
	// URI
	if err := c.ShouldBindUri(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid uri", "detail": err.Error()})
		return
	}
	// Query/form
	if err := c.ShouldBindQuery(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid query/form", "detail": err.Error()})
		return
	}
	// Headers
	if err := c.ShouldBindHeader(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid headers", "detail": err.Error()})
		return
	}
	// JSON
	if c.Request.Method == http.MethodPost || c.Request.Method == http.MethodPut || c.Request.Method == http.MethodPatch {
		if err := c.ShouldBindJSON(&req); err != nil {
			// Skip if body is empty
			if !isEmptyBody(err) {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid json", "detail": err.Error()})
				return
			}
		}
	}
	// Validate once, after every source has been bound
	if err := validate.Struct(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "validation failed", "detail": err.Error()})
		return
	}
	c.Set("data", req)
	c.Next()
}

func isEmptyBody(err error) bool {
	return err.Error() == "EOF"
}
