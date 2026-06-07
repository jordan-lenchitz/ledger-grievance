package domain
import (
	"context"
	"time"
)

type Status string

const (
	StatusReported  Status = "reported"
	StatusReviewing Status = "reviewing"
	StatusResolved  Status = "resolved"
	StatusDismissed Status = "dismissed"
	StatusArchived  Status = "archived"
	StatusCelebrated Status = "celebrated"
)

// Incident represents the database model
// @Description Incident model
type Incident struct {
	ID                    uint64     `json:"id" example:"1"`
	ReporterID            string     `json:"reporter_id" example:"jordan"`
	OccurredAt            *time.Time `json:"occurred_at"`
	RecordedAt            time.Time  `json:"recorded_at"`
	Subject               string     `json:"subject" example:"system failure"`
	Category              string     `json:"category" example:"technical"`
	Severity              uint8      `json:"severity" example:"3"`
	Description           string     `json:"description" example:"the system crashed"`
	EvidenceURI           *string    `json:"evidence_uri" example:"https://example.com/log"`
	RequiresAccommodation bool       `json:"requires_accommodation" example:"false"`
	Status                Status     `json:"status" example:"reported"`
	Notes                 *string    `json:"notes" example:"investigating"`
}

// IncidentCreate represents the input for creating an incident
// @Description Incident creation request
type IncidentCreate struct {
	ReporterID            string     `json:"reporter_id" binding:"required,min=1,max=128" example:"jordan"`
	OccurredAt            *time.Time `json:"occurred_at"`
	Subject               string     `json:"subject" binding:"required,min=1,max=255" example:"system failure"`
	Category              string     `json:"category" binding:"min=1,max=128" example:"technical"`
	Severity              uint8      `json:"severity" binding:"gte=1,lte=5" example:"3"`
	Description           string     `json:"description" binding:"required,min=1" example:"the system crashed"`
	EvidenceURI           *string    `json:"evidence_uri" example:"https://example.com/log"`
	Notes                 *string    `json:"notes" example:"investigating"`
	RequiresAccommodation bool       `json:"requires_accommodation" example:"false"`
	AssumedGoodIntentions bool       `json:"assumed_good_intentions" binding:"required" example:"true"`
	PromisedToBeKindToYourself bool    `json:"promised_to_be_kind_to_yourself" binding:"required" example:"true"`
}

type IncidentPatch struct {
	Status      *Status `json:"status"`
	Notes       *string `json:"notes"`
	Category    *string `json:"category" binding:"omitempty,min=1,max=128"`
	Severity    *uint8  `json:"severity" binding:"omitempty,gte=1,lte=5"`
	EvidenceURI *string `json:"evidence_uri"`
}

type ListParams struct {
	ReporterID string
	Status     string
	Category   string
	Query      string
	Limit      int
	Offset     int
}

type ListResult struct {
	Data  []Incident
	Total int
}

type IncidentRepository interface {
	Create(ctx context.Context, inc *Incident) (uint64, error)
	GetByID(ctx context.Context, id uint64) (*Incident, error)
	List(ctx context.Context, params ListParams) (ListResult, error)
	Update(ctx context.Context, id uint64, patch IncidentPatch) error
	Archive(ctx context.Context, id uint64) error
}
