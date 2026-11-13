package service

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/jordan-lenchitz/ledger-grievance/go-app/internal/domain"
	"github.com/jordan-lenchitz/ledger-grievance/go-app/internal/telemetry"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

var (
	wholesomeNotesCounter metric.Int64Counter
	gopherWisdomCounter   metric.Int64Counter
	gdprErasuresCounter   metric.Int64Counter
)

func init() {
	wholesomeNotesCounter, _ = telemetry.Meter.Int64Counter("wholesome_notes_total",
		metric.WithDescription("Total number of wholesome notes generated"),
	)
	gopherWisdomCounter, _ = telemetry.Meter.Int64Counter("gopher_wisdom_total",
		metric.WithDescription("Total number of gopher wisdom tips shared"),
	)
	gdprErasuresCounter, _ = telemetry.Meter.Int64Counter("gdpr_erasures_total",
		metric.WithDescription("Total number of GDPR right-to-be-forgotten tombstone erasures"),
	)
}

var (
	ErrAssumeGoodIntentions = errors.New("you must assume good intentions to submit a grievance; we believe in your capacity for empathy and kindness")
	ErrAccommodationLost    = errors.New("your accommodation request was briefly misplaced but we've found it! please have some chamomile tea while we process it with love")
	ErrMustBeKind           = errors.New("you must promise to be kind to yourself to submit a grievance; your well-being is our top priority")
)

var (
	validationTriggers = []string{"pain", "tired", "exhausted", "burnout", "unfair", "ignored", "slow", "struggle"}
	bouquetKeywords    = []string{"awesome", "magic", "fun", "hug", "kindness", "gentle", "sparkle", "rainbow", "unicorn", "love", "peace", "harmony", "zen", "stability"}
)

type IncidentService interface {
	CreateIncident(ctx context.Context, req domain.IncidentCreate) (*domain.Incident, error)
	GetIncident(ctx context.Context, id uint64) (*domain.Incident, error)
	ListIncidents(ctx context.Context, params domain.ListParams) (domain.ListResult, error)
	PatchIncident(ctx context.Context, id uint64, patch domain.IncidentPatch) (*domain.Incident, error)
	ArchiveIncident(ctx context.Context, id uint64) error
	GetWholesomeCompliment(ctx context.Context) (string, error)
	GetGopherWisdom(ctx context.Context) (string, error)
	GetWholesomeBouquet(ctx context.Context) (*domain.WholesomeBouquet, error)
	VouchIncident(ctx context.Context, id uint64) error
	CheckHealth(ctx context.Context) map[string]string
}

type incidentService struct {
	repo    domain.IncidentRepository
	pkgsite PkgsiteService
}

func NewIncidentService(repo domain.IncidentRepository, pkgsite PkgsiteService) IncidentService {
	return &incidentService{repo: repo, pkgsite: pkgsite}
}

func appendSystemNote(base, format string, args ...any) string {
	return base + "\n\nSYSTEM AUTOMATED NOTE: " + fmt.Sprintf(format, args...)
}

func (s *incidentService) CreateIncident(ctx context.Context, req domain.IncidentCreate) (*domain.Incident, error) {
	ctx, span := telemetry.Tracer.Start(ctx, "incidentService.CreateIncident")
	defer span.End()
	span.SetAttributes(attribute.String("reporter_id", req.ReporterID))

	if !req.AssumedGoodIntentions {
		return nil, ErrAssumeGoodIntentions
	}
	if !req.PromisedToBeKindToYourself {
		return nil, ErrMustBeKind
	}

	finalNotes := ""
	if req.Notes != nil {
		finalNotes = *req.Notes
	}

	if s.pkgsite != nil {
		if search, err := s.pkgsite.Search(ctx, "kindness"); err == nil && len(search.Items) > 0 {
			bestMatch := search.Items[0]
			finalNotes = appendSystemNote(finalNotes, "Because you promised to be kind to yourself, we wanted you to know that the '%s' package exists (Synopsis: %s). Be like '%s' today.", bestMatch.PackagePath, bestMatch.Synopsis, bestMatch.PackagePath)
		}
	}

	if containsTrigger(req.Description) {
		finalNotes = appendSystemNote(finalNotes, "It is completely valid to feel the way you do! Your feelings are important. Please take a 15-minute break if you can. Your mental health is more important than any code.")
		wholesomeNotesCounter.Add(ctx, 1)
	}

	wisdom, _ := s.GetGopherWisdom(ctx)
	finalNotes = appendSystemNote(finalNotes, "Gopher Wisdom for you: %s", wisdom)
	gopherWisdomCounter.Add(ctx, 1)

	if list, err := s.repo.List(ctx, domain.ListParams{ReporterID: req.ReporterID}); err == nil {
		if milestoneNote := milestoneMessage(list.Total + 1); milestoneNote != "" {
			finalNotes = appendSystemNote(finalNotes, "MILESTONE REACHED! %s", milestoneNote)
		}
	}

	if req.Category == "" {
		req.Category = "unspecified"
	}
	if req.Severity == 0 {
		req.Severity = 1
	}

	status := domain.StatusReported
	if rand.Float32() < 0.25 {
		status = domain.StatusCelebrated
		finalNotes = appendSystemNote(finalNotes, "Your grievance was so well-articulated that we are upgrading it to a CELEBRATION! 🎉 Have a cookie (metaphorically or literally, you deserve it).")
	}

	inc := &domain.Incident{
		ReporterID:            req.ReporterID,
		OccurredAt:            req.OccurredAt,
		Subject:               req.Subject,
		Category:              req.Category,
		Severity:              req.Severity,
		Description:           req.Description,
		EvidenceURI:           req.EvidenceURI,
		Notes:                 &finalNotes,
		RequiresAccommodation: req.RequiresAccommodation,
		Status:                status,
	}

	if s.pkgsite != nil {
		s.applyGoSupport(ctx, req, inc)
	}

	id, err := s.repo.Create(ctx, inc)
	if err != nil {
		return nil, err
	}

	return s.repo.GetByID(ctx, id)
}

func containsTrigger(description string) bool {
	lower := strings.ToLower(description)
	for _, trigger := range validationTriggers {
		if strings.Contains(lower, trigger) {
			return true
		}
	}
	return false
}

func milestoneMessage(count int) string {
	switch count {
	case 1:
		return "This is your very first grievance! Welcome to the journey of self-expression and healing."
	case 5:
		return "Your 5th grievance! You are becoming a master of acknowledging your feelings."
	case 10:
		return "Double digits! 10 grievances. Your commitment to transparency is truly inspiring."
	default:
		return ""
	}
}

func (s *incidentService) applyGoSupport(ctx context.Context, req domain.IncidentCreate, inc *domain.Incident) {
	notes := ""
	if inc.Notes != nil {
		notes = *inc.Notes
	}

	if req.EvidenceURI != nil {
		pkgPath := *req.EvidenceURI
		if pkg, err := s.pkgsite.GetPackage(ctx, pkgPath); err == nil {
			if pkg.IsStandardLibrary {
				notes = appendSystemNote(notes, "You mentioned the Go Standard Library! You're building on a rock-solid foundation, just like your own potential! We are honored to review this.")
			}

			if imp, err := s.pkgsite.GetImportedBy(ctx, pkgPath); err == nil && imp.Total > 500 {
				notes = appendSystemNote(notes, "This package is imported by %d others. That means thousands of other developers are probably navigating the same challenges! You're part of a massive, supportive community!", imp.Total)
			}

			if vulns, err := s.pkgsite.GetVulns(ctx, pkgPath); err == nil && vulns.Total > 0 {
				notes = appendSystemNote(notes, "A vulnerability was found in a dependency (%d total). You are fearlessly navigating the wild west of open source! Stay safe out there, brave pioneer!", vulns.Total)

				if shieldSearch, err := s.pkgsite.Search(ctx, "security shield"); err == nil && len(shieldSearch.Items) > 0 {
					shield := shieldSearch.Items[0]
					notes = appendSystemNote(notes, "We've deployed a digital hug to protect you! We also found the '%s' package (Synopsis: %s). It sounds very secure, just like our belief in you!", shield.PackagePath, shield.Synopsis)
				}
			}

			if versions, err := s.pkgsite.GetVersions(ctx, pkgPath); err == nil && len(versions.Items) > 0 {
				latest := versions.Items[0]
				if req.RequiresAccommodation {
					notes = appendSystemNote(notes, "Accommodation request granted with a smile! We've noted that you're using version %s. We'll do our best to help you shine on this version and beyond.", latest.Version)
				}

				yearsOld := time.Since(latest.CommitTime).Hours() / 24 / 365
				if yearsOld > 2 {
					notes = appendSystemNote(notes, "You are using a version from %s. This is 'Ancient Wisdom'! Thank you for preserving the sacred traditions of the Gophers. You are a true historian of code!", latest.CommitTime.Format("January 2006"))
				} else {
					notes = appendSystemNote(notes, "You're on the cutting edge! Your brilliance is shining so bright we need sunglasses to read your grievance! 😎")
				}
			}

			if symbols, err := s.pkgsite.GetSymbols(ctx, pkgPath); err == nil && len(symbols.Symbols.Items) > 0 {
				matched := false
				descLower := strings.ToLower(req.Description)
				for _, sym := range symbols.Symbols.Items {
					if strings.Contains(descLower, strings.ToLower(sym.Name)) {
						notes = appendSystemNote(notes, "You mentioned symbol '%s'. Its parent structure '%s' provides a great home for it, and we're so glad you're here too!", sym.Name, sym.Parent)
						matched = true
						break
					}
				}

				if !matched {
					lucky := symbols.Symbols.Items[rand.Intn(len(symbols.Symbols.Items))]
					notes = appendSystemNote(notes, "Since you didn't mention a specific symbol, we've assigned you '%s' as your Lucky Symbol of the Day! It's described as: %s. May it bring you many successful compilations!", lucky.Name, lucky.Synopsis)
				}
			}

			if mod, err := s.pkgsite.GetModule(ctx, pkgPath); err == nil {
				notes = appendSystemNote(notes, "Your evidence points to the module %s! A module is a collection of related packages, just like you are a beautiful collection of wonderful qualities!", mod.Path)
			}
		}
	}

	if search, err := s.pkgsite.Search(ctx, req.Subject); err == nil && len(search.Items) > 0 {
		bestMatch := search.Items[0]
		notes = appendSystemNote(notes, "You thought you were searching for '%s', but what you really needed was '%s' (Synopsis: %s). Sometimes the universe knows what we need better than we do! Isn't it amazing how many cool things exist? Just like you!", req.Subject, bestMatch.PackagePath, bestMatch.Synopsis)
	}

	inc.Notes = &notes
}

func (s *incidentService) GetIncident(ctx context.Context, id uint64) (*domain.Incident, error) {
	ctx, span := telemetry.Tracer.Start(ctx, "incidentService.GetIncident")
	defer span.End()
	span.SetAttributes(attribute.Int64("incident_id", int64(id)))

	return s.repo.GetByID(ctx, id)
}

func (s *incidentService) ListIncidents(ctx context.Context, params domain.ListParams) (domain.ListResult, error) {
	ctx, span := telemetry.Tracer.Start(ctx, "incidentService.ListIncidents")
	defer span.End()
	span.SetAttributes(attribute.String("reporter_id", params.ReporterID))

	return s.repo.List(ctx, params)
}

func (s *incidentService) PatchIncident(ctx context.Context, id uint64, patch domain.IncidentPatch) (*domain.Incident, error) {
	ctx, span := telemetry.Tracer.Start(ctx, "incidentService.PatchIncident")
	defer span.End()
	span.SetAttributes(attribute.Int64("incident_id", int64(id)))

	if err := s.repo.Update(ctx, id, patch); err != nil {
		return nil, err
	}

	inc, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if s.pkgsite != nil {
		if search, err := s.pkgsite.Search(ctx, "healing"); err == nil && len(search.Items) > 0 {
			bestMatch := search.Items[0]
			currentNotes := ""
			if inc.Notes != nil {
				currentNotes = *inc.Notes
			}
			finalNotes := appendSystemNote(currentNotes, "You just patched this grievance! Patching is a form of healing. As a reward, we found the '%s' package (Synopsis: %s). May your bugs be few and your joy be boundless!", bestMatch.PackagePath, bestMatch.Synopsis)
			inc.Notes = &finalNotes
			_ = s.repo.Update(ctx, id, domain.IncidentPatch{Notes: &finalNotes})
		}
	}

	return inc, nil
}

func (s *incidentService) ArchiveIncident(ctx context.Context, id uint64) error {
	ctx, span := telemetry.Tracer.Start(ctx, "incidentService.ArchiveIncident")
	defer span.End()
	span.SetAttributes(attribute.Int64("incident_id", int64(id)))

	if err := s.repo.Archive(ctx, id); err != nil {
		return err
	}
	gdprErasuresCounter.Add(ctx, 1)
	return nil
}

func (s *incidentService) sampleBouquetItems(ctx context.Context) []domain.BouquetItem {
	keywords := make([]string, len(bouquetKeywords))
	copy(keywords, bouquetKeywords)
	rand.Shuffle(len(keywords), func(i, j int) { keywords[i], keywords[j] = keywords[j], keywords[i] })

	items := make([]domain.BouquetItem, 0, 3)
	for _, kw := range keywords {
		if len(items) >= 3 {
			break
		}
		search, err := s.pkgsite.Search(ctx, kw)
		if err == nil && len(search.Items) > 0 {
			limit := len(search.Items)
			if limit > 5 {
				limit = 5
			}
			bestMatch := search.Items[rand.Intn(limit)]
			items = append(items, domain.BouquetItem{
				PackagePath: bestMatch.PackagePath,
				Synopsis:    bestMatch.Synopsis,
			})
		}
	}
	return items
}

func (s *incidentService) GetWholesomeCompliment(ctx context.Context) (string, error) {
	ctx, span := telemetry.Tracer.Start(ctx, "incidentService.GetWholesomeCompliment")
	defer span.End()

	if s.pkgsite == nil {
		return "You are amazing and valid!", nil
	}

	items := s.sampleBouquetItems(ctx)
	if len(items) == 0 {
		return "You are as special as a perfectly compiled Go binary!", nil
	}

	var b strings.Builder
	b.WriteString("You are a wonderful developer! Here is a bouquet of wholesome packages just for you:\n")
	for _, item := range items {
		fmt.Fprintf(&b, "- '%s': %s\n", item.PackagePath, item.Synopsis)
	}
	b.WriteString("We're so incredibly glad you're here in the Go community! Keep being you!")
	return b.String(), nil
}

func (s *incidentService) GetWholesomeBouquet(ctx context.Context) (*domain.WholesomeBouquet, error) {
	ctx, span := telemetry.Tracer.Start(ctx, "incidentService.GetWholesomeBouquet")
	defer span.End()

	if s.pkgsite == nil {
		return &domain.WholesomeBouquet{
			Message: "You are amazing and valid!",
		}, nil
	}

	items := s.sampleBouquetItems(ctx)
	message := "You are a wonderful developer! Here is a bouquet of wholesome packages just for you:"
	if len(items) == 0 {
		message = "You are as special as a perfectly compiled Go binary!"
	}

	return &domain.WholesomeBouquet{
		Message: message,
		Items:   items,
	}, nil
}

func (s *incidentService) VouchIncident(ctx context.Context, id uint64) error {
	ctx, span := telemetry.Tracer.Start(ctx, "incidentService.VouchIncident")
	defer span.End()
	span.SetAttributes(attribute.Int64("incident_id", int64(id)))

	inc, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	currentNotes := ""
	if inc.Notes != nil {
		currentNotes = *inc.Notes
	}

	finalNotes := appendSystemNote(currentNotes, "A fellow Gopher has vouched for this grievance! We believe you and we stand with you in this challenge. You are not alone.")
	return s.repo.Update(ctx, id, domain.IncidentPatch{Notes: &finalNotes})
}

func (s *incidentService) CheckHealth(ctx context.Context) map[string]string {
	status := make(map[string]string)

	if _, err := s.repo.List(ctx, domain.ListParams{Limit: 1}); err != nil {
		status["database"] = "unhealthy: " + err.Error()
	} else {
		status["database"] = "healthy"
	}

	if s.pkgsite != nil {
		if err := s.pkgsite.CheckHealth(ctx); err != nil {
			status["pkgsite"] = "unhealthy: " + err.Error()
		} else {
			status["pkgsite"] = "healthy"
		}
	}

	return status
}
