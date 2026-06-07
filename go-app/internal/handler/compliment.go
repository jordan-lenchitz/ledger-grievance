package handler

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

// GetCompliment returns a random wholesome compliment
// @Summary Get a wholesome compliment
// @Tags compliments
// @Produce json
// @Success 200 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /compliments [get]
func (h *IncidentHandler) GetCompliment(c *gin.Context) {
	compliment, err := h.svc.GetWholesomeCompliment(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Even though we couldn't find a package, you're still great!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"compliment": compliment})
}
