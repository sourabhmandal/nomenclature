package translation

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Handler exposes translation APIs.
type TranslationHandler struct {
	service *TranslationService
}

func NewTranslationHandler(service *TranslationService) *TranslationHandler {
	return &TranslationHandler{
		service: service,
	}
}

// TranslateHandler handles translation requests.
func (h *TranslationHandler) TranslateHandler(c *gin.Context) {
	var req TranslationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	res, err := h.service.Translate(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, res)
}
