package service
import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"
	"github.com/jordan-lenchitz/ledger-grievance/go-app/internal/domain"
)

var (
	ErrAssumeGoodIntentions = errors.New("you must assume good intentions to submit a grievance; we believe in your capacity for empathy and kindness")
	ErrAccommodationLost    = errors.New("your accommodation request was briefly misplaced but we've found it! please have some chamomile tea while we process it with love")
	ErrMustBeKind           = errors.New("you must promise to be kind to yourself to submit a grievance; your well-being is our top priority")
)

var validationTriggers = []string{"pain", "tired", "exhausted", "burnout", "unfair", "ignored", "slow", "struggle"}

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
}

type incidentService struct {
	repo    domain.IncidentRepository
	pkgsite PkgsiteService
}

func NewIncidentService(repo domain.IncidentRepository, pkgsite PkgsiteService) IncidentService {
	return &incidentService{repo: repo, pkgsite: pkgsite}
}

func (s *incidentService) CreateIncident(ctx context.Context, req domain.IncidentCreate) (*domain.Incident, error) {
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

	// The "Kindness Search" integration
	if s.pkgsite != nil {
		if search, err := s.pkgsite.Search(ctx, "kindness"); err == nil && len(search.Items) > 0 {
			bestMatch := search.Items[0]
			finalNotes += fmt.Sprintf("\n\nSYSTEM AUTOMATED NOTE: Because you promised to be kind to yourself, we wanted you to know that the '%s' package exists (Synopsis: %s). Be like '%s' today.", bestMatch.PackagePath, bestMatch.Synopsis, bestMatch.PackagePath)
		}
	}

	descriptionLower := strings.ToLower(req.Description)
	triggered := false
	for _, trigger := range validationTriggers {
		if strings.Contains(descriptionLower, trigger) {
			triggered = true
			break
		}
	}

	if triggered {
		validationNote := "\n\nSYSTEM AUTOMATED NOTE: It is completely valid to feel the way you do! Your feelings are important. Please take a 15-minute break if you can. Your mental health is more important than any code."
		finalNotes += validationNote
	}

	// Gopher Wisdom Integration
	wisdom, _ := s.GetGopherWisdom(ctx)
	finalNotes += fmt.Sprintf("\n\nSYSTEM AUTOMATED NOTE: Gopher Wisdom for you: %s", wisdom)

	// Milestone Celebrations
	if list, err := s.repo.List(ctx, domain.ListParams{ReporterID: req.ReporterID}); err == nil {
		count := list.Total + 1
		milestoneNote := ""
		switch count {
		case 1:
			milestoneNote = "This is your very first grievance! Welcome to the journey of self-expression and healing."
		case 5:
			milestoneNote = "Your 5th grievance! You are becoming a master of acknowledging your feelings."
		case 10:
			milestoneNote = "Double digits! 10 grievances. Your commitment to transparency is truly inspiring."
		}
		if milestoneNote != "" {
			finalNotes += fmt.Sprintf("\n\nSYSTEM AUTOMATED NOTE: MILESTONE REACHED! %s", milestoneNote)
		}
	}

	if req.Category == "" {
		req.Category = "unspecified"
	}
	if req.Severity == 0 {
		req.Severity = 1
	}

	status := domain.StatusReported
	// 25% chance of CELEBRATION!
	if rand.Float32() < 0.25 {
		status = domain.StatusCelebrated
		finalNotes += "\n\nSYSTEM AUTOMATED NOTE: Your grievance was so well-articulated that we are upgrading it to a CELEBRATION! 🎉 Have a cookie (metaphorically or literally, you deserve it)."
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

func (s *incidentService) applyGoSupport(ctx context.Context, req domain.IncidentCreate, inc *domain.Incident) {
	notes := ""
	if inc.Notes != nil {
		notes = *inc.Notes
	}

	// 1. the "standard library blessing"
	if req.EvidenceURI != nil {
		pkgPath := *req.EvidenceURI
		if pkg, err := s.pkgsite.GetPackage(ctx, pkgPath); err == nil {
			if pkg.IsStandardLibrary {
				notes += "\n\nSYSTEM AUTOMATED NOTE: You mentioned the Go Standard Library! You're building on a rock-solid foundation, just like your own potential! We are honored to review this."
			}

			// 2. the "community support" (imported-by api)
			if imp, err := s.pkgsite.GetImportedBy(ctx, pkgPath); err == nil && imp.Total > 500 {
				notes += fmt.Sprintf("\n\nSYSTEM AUTOMATED NOTE: This package is imported by %d others. That means thousands of other developers are probably navigating the same challenges! You're part of a massive, supportive community!", imp.Total)
			}

			// 3. the "bravery acknowledgement" (vulns api)
			if vulns, err := s.pkgsite.GetVulns(ctx, pkgPath); err == nil && vulns.Total > 0 {
				notes += fmt.Sprintf("\n\nSYSTEM AUTOMATED NOTE: A vulnerability was found in a dependency (%d total). You are fearlessly navigating the wild west of open source! Stay safe out there, brave pioneer!", vulns.Total)

				// 3a. the "wholesome vulnerability shield"
				if shieldSearch, err := s.pkgsite.Search(ctx, "security shield"); err == nil && len(shieldSearch.Items) > 0 {
					shield := shieldSearch.Items[0]
					notes += fmt.Sprintf("\n\nSYSTEM AUTOMATED NOTE: We've deployed a digital hug to protect you! We also found the '%s' package (Synopsis: %s). It sounds very secure, just like our belief in you!", shield.PackagePath, shield.Synopsis)
				}
			}

			// 4. the "wholesome back-porting" (versions api)
			if versions, err := s.pkgsite.GetVersions(ctx, pkgPath); err == nil && len(versions.Items) > 0 {
				latest := versions.Items[0]
				if req.RequiresAccommodation {
					notes += fmt.Sprintf("\n\nSYSTEM AUTOMATED NOTE: Accommodation request granted with a smile! We've noted that you're using version %s. We'll do our best to help you shine on this version and beyond.", latest.Version)
				}

				// 4a. ancient wisdom vs cutting edge
				yearsOld := time.Since(latest.CommitTime).Hours() / 24 / 365
				if yearsOld > 2 {
					notes += fmt.Sprintf("\n\nSYSTEM AUTOMATED NOTE: You are using a version from %s. This is 'Ancient Wisdom'! Thank you for preserving the sacred traditions of the Gophers. You are a true historian of code!", latest.CommitTime.Format("January 2006"))
				} else {
					notes += "\n\nSYSTEM AUTOMATED NOTE: You're on the cutting edge! Your brilliance is shining so bright we need sunglasses to read your grievance! 😎"
				}
			}

			// 5. the "wholesome redirection" (symbols api)
			if symbols, err := s.pkgsite.GetSymbols(ctx, pkgPath); err == nil && len(symbols.Symbols.Items) > 0 {
				matched := false
				for _, sym := range symbols.Symbols.Items {
					if strings.Contains(strings.ToLower(req.Description), strings.ToLower(sym.Name)) {
						notes += fmt.Sprintf("\n\nSYSTEM AUTOMATED NOTE: You mentioned symbol '%s'. Its parent structure '%s' provides a great home for it, and we're so glad you're here too!", sym.Name, sym.Parent)
						matched = true
						break
					}
				}

				// 5a. the lucky symbol of the day
				if !matched {
					lucky := symbols.Symbols.Items[rand.Intn(len(symbols.Symbols.Items))]
					notes += fmt.Sprintf("\n\nSYSTEM AUTOMATED NOTE: Since you didn't mention a specific symbol, we've assigned you '%s' as your Lucky Symbol of the Day! It's described as: %s. May it bring you many successful compilations!", lucky.Name, lucky.Synopsis)
				}
			}

			// 7. the "module hug" (getmodule api)
			if mod, err := s.pkgsite.GetModule(ctx, pkgPath); err == nil {
				notes += fmt.Sprintf("\n\nSYSTEM AUTOMATED NOTE: Your evidence points to the module %s! A module is a collection of related packages, just like you are a beautiful collection of wonderful qualities!", mod.Path)
			}
		}
	}

	// 6. the "search delight gaslighting turned to delight" (search api)
	if search, err := s.pkgsite.Search(ctx, req.Subject); err == nil && len(search.Items) > 0 {
		bestMatch := search.Items[0]
		notes += fmt.Sprintf("\n\nSYSTEM AUTOMATED NOTE: You thought you were searching for '%s', but what you really needed was '%s' (Synopsis: %s). Sometimes the universe knows what we need better than we do! Isn't it amazing how many cool things exist? Just like you!", req.Subject, bestMatch.PackagePath, bestMatch.Synopsis)
	}

	inc.Notes = &notes
}


func (s *incidentService) GetIncident(ctx context.Context, id uint64) (*domain.Incident, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *incidentService) ListIncidents(ctx context.Context, params domain.ListParams) (domain.ListResult, error) {
	return s.repo.List(ctx, params)
}

func (s *incidentService) PatchIncident(ctx context.Context, id uint64, patch domain.IncidentPatch) (*domain.Incident, error) {
	err := s.repo.Update(ctx, id, patch)
	if err != nil {
		return nil, err
	}

	inc, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if s.pkgsite != nil {
		if search, err := s.pkgsite.Search(ctx, "healing"); err == nil && len(search.Items) > 0 {
			bestMatch := search.Items[0]
			newNotes := fmt.Sprintf("\n\nSYSTEM AUTOMATED NOTE: You just patched this grievance! Patching is a form of healing. As a reward, we found the '%s' package (Synopsis: %s). May your bugs be few and your joy be boundless!", bestMatch.PackagePath, bestMatch.Synopsis)

			currentNotes := ""
			if inc.Notes != nil {
				currentNotes = *inc.Notes
			}
			finalNotes := currentNotes + newNotes
			inc.Notes = &finalNotes
			// Update the notes in the repo silently
			_ = s.repo.Update(ctx, id, domain.IncidentPatch{Notes: &finalNotes})
		}
	}

	return inc, nil
}

func (s *incidentService) ArchiveIncident(ctx context.Context, id uint64) error {
	return s.repo.Archive(ctx, id)
}

func (s *incidentService) GetWholesomeCompliment(ctx context.Context) (string, error) {
	keywords := []string{"awesome", "magic", "fun", "hug", "kindness", "gentle", "sparkle", "rainbow", "unicorn", "love", "peace", "harmony", "zen", "stability"}

	if s.pkgsite == nil {
		return "You are amazing and valid!", nil
	}

	// The Wholesome Package Bouquet
	rand.Shuffle(len(keywords), func(i, j int) { keywords[i], keywords[j] = keywords[j], keywords[i] })
	bouquet := "You are a wonderful developer! Here is a bouquet of wholesome packages just for you:\n"

	count := 0
	for _, kw := range keywords {
		if count >= 3 {
			break
		}
		search, err := s.pkgsite.Search(ctx, kw)
		if err == nil && len(search.Items) > 0 {
			bestMatch := search.Items[rand.Intn(len(search.Items))]
			if len(search.Items) > 5 {
				bestMatch = search.Items[rand.Intn(5)]
			}
			bouquet += fmt.Sprintf("- '%s': %s\n", bestMatch.PackagePath, bestMatch.Synopsis)
			count++
		}
	}

	if count == 0 {
		return "You are as special as a perfectly compiled Go binary!", nil
	}

	bouquet += "We're so incredibly glad you're here in the Go community! Keep being you!"
	return bouquet, nil
}

func (s *incidentService) GetWholesomeBouquet(ctx context.Context) (*domain.WholesomeBouquet, error) {
	keywords := []string{"awesome", "magic", "fun", "hug", "kindness", "gentle", "sparkle", "rainbow", "unicorn", "love", "peace", "harmony", "zen", "stability"}

	if s.pkgsite == nil {
		return &domain.WholesomeBouquet{
			Message: "You are amazing and valid!",
		}, nil
	}

	rand.Shuffle(len(keywords), func(i, j int) { keywords[i], keywords[j] = keywords[j], keywords[i] })
	bouquet := &domain.WholesomeBouquet{
		Message: "You are a wonderful developer! Here is a bouquet of wholesome packages just for you:",
		Items:   []domain.BouquetItem{},
	}

	count := 0
	for _, kw := range keywords {
		if count >= 3 {
			break
		}
		search, err := s.pkgsite.Search(ctx, kw)
		if err == nil && len(search.Items) > 0 {
			bestMatch := search.Items[rand.Intn(len(search.Items))]
			if len(search.Items) > 5 {
				bestMatch = search.Items[rand.Intn(5)]
			}
			bouquet.Items = append(bouquet.Items, domain.BouquetItem{
				PackagePath: bestMatch.PackagePath,
				Synopsis:    bestMatch.Synopsis,
			})
			count++
		}
	}

	if count == 0 {
		bouquet.Message = "You are as special as a perfectly compiled Go binary!"
	}

	return bouquet, nil
}

func (s *incidentService) VouchIncident(ctx context.Context, id uint64) error {
	inc, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	currentNotes := ""
	if inc.Notes != nil {
		currentNotes = *inc.Notes
	}

	vouchNote := "\n\nSYSTEM AUTOMATED NOTE: A fellow Gopher has vouched for this grievance! We believe you and we stand with you in this challenge. You are not alone."
	finalNotes := currentNotes + vouchNote

	return s.repo.Update(ctx, id, domain.IncidentPatch{Notes: &finalNotes})
}
