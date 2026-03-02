package translation

import "github.com/gin-gonic/gin"

// Handler exposes translation APIs.
type TranslationHandler interface {
	TranslateHandler(c *gin.Context)
}
