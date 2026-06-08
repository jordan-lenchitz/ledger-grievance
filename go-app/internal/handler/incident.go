package handler
import (
	"database/sql"
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"
	"github.com/jordan-lenchitz/ledger-grievance/go-app/internal/domain"
	"github.com/jordan-lenchitz/ledger-grievance/go-app/internal/service"
)

type IncidentHandler struct {
	svc service.IncidentService
}

func NewIncidentHandler(svc service.IncidentService) *IncidentHandler {
	return &IncidentHandler{svc: svc}
}

func (h *IncidentHandler) RegisterRoutes(r *gin.Engine) {
	r.POST("/incidents", h.CreateIncident)
	r.GET("/incidents", h.ListIncidents)
	r.GET("/incidents/:id", h.GetIncident)
	r.PATCH("/incidents/:id", h.PatchIncident)
	r.DELETE("/incidents/:id", h.ArchiveIncident)
	r.GET("/compliments", h.GetCompliment)
	r.GET("/wisdom", h.GetWisdom)
	r.GET("/bouquet", h.GetBouquet)
	r.POST("/incidents/:id/vouch", h.VouchIncident)
	r.GET("/health/deep", h.GetDeepHealth)
}

// CreateIncident creates a new incident
// @Summary Create a new incident
// @Tags incidents
// @Accept json
// @Produce json
// @Param incident body domain.IncidentCreate true "Incident request"
// @Success 201 {object} domain.Incident
// @Failure 400 {object} gin.H
// @Failure 418 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /incidents [post]
func (h *IncidentHandler) CreateIncident(c *gin.Context) {
	var payload domain.IncidentCreate
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "path": c.Request.URL.Path})
		return
	}

	incident, err := h.svc.CreateIncident(c.Request.Context(), payload)
	if err != nil {
		if err == service.ErrAssumeGoodIntentions || err == service.ErrMustBeKind {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "path": c.Request.URL.Path})
			return
		}
		if err == service.ErrAccommodationLost {
			c.JSON(http.StatusTeapot, gin.H{"error": err.Error(), "path": c.Request.URL.Path})
			return
		}
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, incident)
}

// ListIncidents lists incidents with pagination
// @Summary List incidents
// @Tags incidents
// @Produce json
// @Param reporter_id query string false "Reporter ID"
// @Param status query string false "Status"
// @Param category query string false "Category"
// @Param q query string false "Search query"
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Success 200 {object} domain.ListResult
// @Failure 500 {object} gin.H
// @Router /incidents [get]
func (h *IncidentHandler) ListIncidents(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	if limit < 1 {
		limit = 1
	}
	if limit > 200 {
		limit = 200
	}

	params := domain.ListParams{
		ReporterID: c.Query("reporter_id"),
		Status:     c.Query("status"),
		Category:   c.Query("category"),
		Query:      c.Query("q"),
		Limit:      limit,
		Offset:     offset,
	}

	result, err := h.svc.ListIncidents(c.Request.Context(), params)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": result.Data,
		"meta": map[string]interface{}{
			"total":  result.Total,
			"limit":  limit,
			"offset": offset,
			"count":  len(result.Data),
		},
	})
}

func (h *IncidentHandler) GetIncident(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	incident, err := h.svc.GetIncident(c.Request.Context(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Incident not found", "path": c.Request.URL.Path})
			return
		}
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, incident)
}

func (h *IncidentHandler) PatchIncident(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var payload domain.IncidentPatch
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "path": c.Request.URL.Path})
		return
	}

	incident, err := h.svc.PatchIncident(c.Request.Context(), id, payload)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Incident not found", "path": c.Request.URL.Path})
			return
		}
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, incident)
}

func (h *IncidentHandler) ArchiveIncident(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	err = h.svc.ArchiveIncident(c.Request.Context(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Incident not found", "path": c.Request.URL.Path})
			return
		}
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": id, "status": domain.StatusArchived})
}

func (h *IncidentHandler) GetCompliment(c *gin.Context) {
	compliment, err := h.svc.GetWholesomeCompliment(c.Request.Context())
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"compliment": compliment})
}

func (h *IncidentHandler) GetWisdom(c *gin.Context) {
	wisdom, err := h.svc.GetGopherWisdom(c.Request.Context())
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"wisdom": wisdom})
}

func (h *IncidentHandler) GetBouquet(c *gin.Context) {
	bouquet, err := h.svc.GetWholesomeBouquet(c.Request.Context())
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, bouquet)
}

func (h *IncidentHandler) VouchIncident(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	err = h.svc.VouchIncident(c.Request.Context(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Incident not found", "path": c.Request.URL.Path})
			return
		}
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": id, "status": "vouched"})
}

func (h *IncidentHandler) GetDeepHealth(c *gin.Context) {
	status := h.svc.CheckHealth(c.Request.Context())
	
	isHealthy := true
	for _, v := range status {
		if v != "healthy" {
			isHealthy = false
			break
		}
	}

	code := http.StatusOK
	if !isHealthy {
		code = http.StatusServiceUnavailable
	}

	c.JSON(code, gin.H{
		"status":  status,
		"overall": isHealthy,
	})
}
