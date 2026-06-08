package service

import (
	"context"
	"testing"
	"time"

	"github.com/jordan-lenchitz/ledger-grievance/go-app/internal/domain"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestCreateIncident(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockIncidentRepository(ctrl)
	mockPkgsite := NewMockPkgsiteService(ctrl)
	svc := NewIncidentService(mockRepo, mockPkgsite)

	ctx := context.Background()

	req := domain.IncidentCreate{
		ReporterID:                 "jordan",
		Subject:                    "test incident",
		Description:                "a simple test",
		AssumedGoodIntentions:      true,
		PromisedToBeKindToYourself: true,
	}

	// For kindness promise search
	mockPkgsite.EXPECT().Search(gomock.Any(), "kindness").Return(&domain.PaginatedResponse[domain.SearchResult]{
		Items: []domain.SearchResult{{PackagePath: "kindness", Synopsis: "Be kind"}},
	}, nil)

	// For subject search
	mockPkgsite.EXPECT().Search(gomock.Any(), req.Subject).Return(&domain.PaginatedResponse[domain.SearchResult]{}, nil)

	mockRepo.EXPECT().
		List(gomock.Any(), domain.ListParams{ReporterID: req.ReporterID}).
		Return(domain.ListResult{Total: 0}, nil)

	mockRepo.EXPECT().
		Create(gomock.Any(), gomock.Any()).
		Return(uint64(1), nil)

	mockRepo.EXPECT().
		GetByID(gomock.Any(), uint64(1)).
		Return(&domain.Incident{ID: 1, ReporterID: "jordan"}, nil)

	incident, err := svc.CreateIncident(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, incident)
	assert.Equal(t, uint64(1), incident.ID)
}

func TestCreateIncident_WholesomeValidation(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockIncidentRepository(ctrl)
	mockPkgsite := NewMockPkgsiteService(ctrl)
	svc := NewIncidentService(mockRepo, mockPkgsite)

	ctx := context.Background()

	req := domain.IncidentCreate{
		ReporterID:                 "jordan",
		Subject:                    "test incident",
		Description:                "i am feeling burnout",
		AssumedGoodIntentions:      true,
		PromisedToBeKindToYourself: true,
	}

	mockPkgsite.EXPECT().Search(gomock.Any(), "kindness").Return(&domain.PaginatedResponse[domain.SearchResult]{}, nil)
	mockPkgsite.EXPECT().Search(gomock.Any(), req.Subject).Return(&domain.PaginatedResponse[domain.SearchResult]{}, nil)

	mockRepo.EXPECT().
		List(gomock.Any(), domain.ListParams{ReporterID: req.ReporterID}).
		Return(domain.ListResult{Total: 0}, nil)

	mockRepo.EXPECT().
		Create(gomock.Any(), gomock.Any()).
		Do(func(ctx context.Context, inc *domain.Incident) {
			assert.Contains(t, *inc.Notes, "completely valid to feel the way you do")
		}).
		Return(uint64(1), nil)

	mockRepo.EXPECT().
		GetByID(gomock.Any(), uint64(1)).
		Return(&domain.Incident{ID: 1}, nil)

	_, err := svc.CreateIncident(ctx, req)

	assert.NoError(t, err)
}

func TestListIncidents(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockIncidentRepository(ctrl)
	mockPkgsite := NewMockPkgsiteService(ctrl)
	svc := NewIncidentService(mockRepo, mockPkgsite)

	ctx := context.Background()
	params := domain.ListParams{ReporterID: "jordan", Limit: 10}

	mockRepo.EXPECT().
		List(gomock.Any(), params).
		Return(domain.ListResult{
			Data:  []domain.Incident{{ID: 1}, {ID: 2}},
			Total: 2,
		}, nil)

	result, err := svc.ListIncidents(ctx, params)

	assert.NoError(t, err)
	assert.Equal(t, 2, len(result.Data))
	assert.Equal(t, 2, result.Total)
}

func TestPatchIncident(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockIncidentRepository(ctrl)
	mockPkgsite := NewMockPkgsiteService(ctrl)
	svc := NewIncidentService(mockRepo, mockPkgsite)

	ctx := context.Background()
	id := uint64(1)
	patch := domain.IncidentPatch{Status: (*domain.Status)(&[]domain.Status{domain.StatusResolved}[0])}

	mockRepo.EXPECT().Update(gomock.Any(), id, patch).Return(nil)
	mockRepo.EXPECT().GetByID(gomock.Any(), id).Return(&domain.Incident{ID: id, Status: domain.StatusResolved}, nil)

	mockPkgsite.EXPECT().Search(gomock.Any(), "healing").Return(&domain.PaginatedResponse[domain.SearchResult]{
		Items: []domain.SearchResult{{PackagePath: "healing", Synopsis: "Heal"}},
	}, nil)
	mockRepo.EXPECT().Update(gomock.Any(), id, gomock.Any()).Return(nil)

	incident, err := svc.PatchIncident(ctx, id, patch)

	assert.NoError(t, err)
	assert.Equal(t, domain.StatusResolved, incident.Status)
}

func TestArchiveIncident(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockIncidentRepository(ctrl)
	mockPkgsite := NewMockPkgsiteService(ctrl)
	svc := NewIncidentService(mockRepo, mockPkgsite)

	ctx := context.Background()
	id := uint64(1)

	mockRepo.EXPECT().Archive(gomock.Any(), id).Return(nil)

	err := svc.ArchiveIncident(ctx, id)

	assert.NoError(t, err)
}

func TestCreateIncident_AssumeGoodIntentions_False(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockIncidentRepository(ctrl)
	mockPkgsite := NewMockPkgsiteService(ctrl)
	svc := NewIncidentService(mockRepo, mockPkgsite)

	ctx := context.Background()
	req := domain.IncidentCreate{
		ReporterID:            "jordan",
		AssumedGoodIntentions: false,
	}

	_, err := svc.CreateIncident(ctx, req)

	assert.Error(t, err)
	assert.Equal(t, ErrAssumeGoodIntentions, err)
}

func TestCreateIncident_PromisedToBeKindToYourself_False(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockIncidentRepository(ctrl)
	mockPkgsite := NewMockPkgsiteService(ctrl)
	svc := NewIncidentService(mockRepo, mockPkgsite)

	ctx := context.Background()
	req := domain.IncidentCreate{
		ReporterID:                 "jordan",
		AssumedGoodIntentions:      true,
		PromisedToBeKindToYourself: false,
	}

	_, err := svc.CreateIncident(ctx, req)

	assert.Error(t, err)
	assert.Equal(t, ErrMustBeKind, err)
}

func TestGetWholesomeCompliment(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockIncidentRepository(ctrl)
	mockPkgsite := NewMockPkgsiteService(ctrl)
	svc := NewIncidentService(mockRepo, mockPkgsite)

	ctx := context.Background()

	mockPkgsite.EXPECT().Search(gomock.Any(), gomock.Any()).Return(&domain.PaginatedResponse[domain.SearchResult]{
		Items: []domain.SearchResult{{PackagePath: "github.com/awesome/pkg", Synopsis: "Something great"}},
	}, nil).AnyTimes()

	compliment, err := svc.GetWholesomeCompliment(ctx)

	assert.NoError(t, err)
	assert.Contains(t, compliment, "bouquet")
	assert.Contains(t, compliment, "github.com/awesome/pkg")
}

func TestCreateIncident_GoSupport_Table(t *testing.T) {
	tests := []struct {
		name          string
		evidence      string
		description   string
		accommodation bool
		setupMocks    func(*MockPkgsiteService)
		expectedNotes []string
	}{
		{
			name:     "Standard Library Blessing",
			evidence: "fmt",
			setupMocks: func(m *MockPkgsiteService) {
				m.EXPECT().Search(gomock.Any(), "kindness").Return(&domain.PaginatedResponse[domain.SearchResult]{}, nil)
				m.EXPECT().GetPackage(gomock.Any(), "fmt").Return(&domain.Package{IsStandardLibrary: true}, nil)
				m.EXPECT().GetImportedBy(gomock.Any(), "fmt").Return(&domain.PackageImportedBy{}, nil)
				m.EXPECT().GetVulns(gomock.Any(), "fmt").Return(&domain.PaginatedResponse[domain.Vulnerability]{}, nil)
				m.EXPECT().GetVersions(gomock.Any(), "fmt").Return(&domain.PaginatedResponse[domain.ModuleVersion]{
					Items: []domain.ModuleVersion{{Version: "v1.0.0", CommitTime: time.Now()}},
				}, nil)
				m.EXPECT().GetSymbols(gomock.Any(), "fmt").Return(&domain.PackageSymbols{
					Symbols: struct {
						Items []domain.Symbol `json:"items"`
						Total int             `json:"total"`
					}{Items: []domain.Symbol{{Name: "Println"}}},
				}, nil)
				m.EXPECT().GetModule(gomock.Any(), "fmt").Return(&domain.Module{Path: "std"}, nil)
				m.EXPECT().Search(gomock.Any(), gomock.Any()).Return(&domain.PaginatedResponse[domain.SearchResult]{}, nil).AnyTimes()
			},
			expectedNotes: []string{"rock-solid foundation", "honored to review", "cutting edge"},
		},
		{
			name:     "Community Support",
			evidence: "github.com/popular/pkg",
			setupMocks: func(m *MockPkgsiteService) {
				m.EXPECT().Search(gomock.Any(), "kindness").Return(&domain.PaginatedResponse[domain.SearchResult]{}, nil)
				m.EXPECT().GetPackage(gomock.Any(), "github.com/popular/pkg").Return(&domain.Package{}, nil)
				m.EXPECT().GetImportedBy(gomock.Any(), "github.com/popular/pkg").Return(&domain.PackageImportedBy{Total: 1000}, nil)
				m.EXPECT().GetVulns(gomock.Any(), "github.com/popular/pkg").Return(&domain.PaginatedResponse[domain.Vulnerability]{}, nil)
				m.EXPECT().GetVersions(gomock.Any(), "github.com/popular/pkg").Return(&domain.PaginatedResponse[domain.ModuleVersion]{
					Items: []domain.ModuleVersion{{Version: "v1.0.0", CommitTime: time.Now().AddDate(-5, 0, 0)}},
				}, nil)
				m.EXPECT().GetSymbols(gomock.Any(), "github.com/popular/pkg").Return(&domain.PackageSymbols{
					Symbols: struct {
						Items []domain.Symbol `json:"items"`
						Total int             `json:"total"`
					}{Items: []domain.Symbol{{Name: "PopularFunc"}}},
				}, nil)
				m.EXPECT().GetModule(gomock.Any(), "github.com/popular/pkg").Return(&domain.Module{Path: "github.com/popular/pkg"}, nil)
				m.EXPECT().Search(gomock.Any(), gomock.Any()).Return(&domain.PaginatedResponse[domain.SearchResult]{}, nil).AnyTimes()
			},
			expectedNotes: []string{"thousands of other developers", "supportive community", "Ancient Wisdom"},
		},
		{
			name:     "Bravery Acknowledgement",
			evidence: "github.com/insecure/pkg",
			setupMocks: func(m *MockPkgsiteService) {
				m.EXPECT().Search(gomock.Any(), "kindness").Return(&domain.PaginatedResponse[domain.SearchResult]{}, nil)
				m.EXPECT().GetPackage(gomock.Any(), "github.com/insecure/pkg").Return(&domain.Package{}, nil)
				m.EXPECT().GetImportedBy(gomock.Any(), "github.com/insecure/pkg").Return(&domain.PackageImportedBy{}, nil)
				m.EXPECT().GetVulns(gomock.Any(), "github.com/insecure/pkg").Return(&domain.PaginatedResponse[domain.Vulnerability]{Total: 5}, nil)
				m.EXPECT().Search(gomock.Any(), "security shield").Return(&domain.PaginatedResponse[domain.SearchResult]{
					Items: []domain.SearchResult{{PackagePath: "shield", Synopsis: "Protect"}},
				}, nil)
				m.EXPECT().GetVersions(gomock.Any(), "github.com/insecure/pkg").Return(&domain.PaginatedResponse[domain.ModuleVersion]{
					Items: []domain.ModuleVersion{{Version: "v1.0.0", CommitTime: time.Now()}},
				}, nil)
				m.EXPECT().GetSymbols(gomock.Any(), "github.com/insecure/pkg").Return(&domain.PackageSymbols{
					Symbols: struct {
						Items []domain.Symbol `json:"items"`
						Total int             `json:"total"`
					}{Items: []domain.Symbol{{Name: "InsecureFunc"}}},
				}, nil)
				m.EXPECT().GetModule(gomock.Any(), "github.com/insecure/pkg").Return(&domain.Module{Path: "github.com/insecure/pkg"}, nil)
				m.EXPECT().Search(gomock.Any(), gomock.Any()).Return(&domain.PaginatedResponse[domain.SearchResult]{}, nil).AnyTimes()
			},
			expectedNotes: []string{"fearlessly navigating", "brave pioneer", "deployed a digital hug"},
		},
		{
			name: "Search Delight",
			setupMocks: func(m *MockPkgsiteService) {
				m.EXPECT().Search(gomock.Any(), "kindness").Return(&domain.PaginatedResponse[domain.SearchResult]{}, nil)
				m.EXPECT().Search(gomock.Any(), "my problem").Return(&domain.PaginatedResponse[domain.SearchResult]{
					Items: []domain.SearchResult{{PackagePath: "github.com/helpful/pkg", Synopsis: "Helping is fun"}},
				}, nil)
			},
			expectedNotes: []string{"amazing how many cool things exist", "github.com/helpful/pkg"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := NewMockIncidentRepository(ctrl)
			mockPkgsite := NewMockPkgsiteService(ctrl)
			svc := NewIncidentService(mockRepo, mockPkgsite)

			tt.setupMocks(mockPkgsite)

			req := domain.IncidentCreate{
				ReporterID:                 "jordan",
				Subject:                    "my problem",
				Description:                tt.description,
				EvidenceURI:                &tt.evidence,
				RequiresAccommodation:      tt.accommodation,
				AssumedGoodIntentions:      true,
				PromisedToBeKindToYourself: true,
			}
			if tt.evidence == "" {
				req.EvidenceURI = nil
			}

			mockRepo.EXPECT().
				List(gomock.Any(), domain.ListParams{ReporterID: req.ReporterID}).
				Return(domain.ListResult{Total: 0}, nil)

			mockRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Do(func(ctx context.Context, inc *domain.Incident) {
				for _, note := range tt.expectedNotes {
					assert.Contains(t, *inc.Notes, note)
				}
			}).Return(uint64(1), nil)

			mockRepo.EXPECT().GetByID(gomock.Any(), uint64(1)).Return(&domain.Incident{}, nil)

			_, err := svc.CreateIncident(context.Background(), req)
			assert.NoError(t, err)
		})
	}
}

func TestGetGopherWisdom(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockIncidentRepository(ctrl)
	mockPkgsite := NewMockPkgsiteService(ctrl)
	svc := NewIncidentService(mockRepo, mockPkgsite)

	wisdom, err := svc.GetGopherWisdom(context.Background())
	assert.NoError(t, err)
	assert.NotEmpty(t, wisdom)
}

func TestGetWholesomeBouquet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockIncidentRepository(ctrl)
	mockPkgsite := NewMockPkgsiteService(ctrl)
	svc := NewIncidentService(mockRepo, mockPkgsite)

	mockPkgsite.EXPECT().Search(gomock.Any(), gomock.Any()).Return(&domain.PaginatedResponse[domain.SearchResult]{
		Items: []domain.SearchResult{{PackagePath: "github.com/awesome/pkg", Synopsis: "Something great"}},
	}, nil).AnyTimes()

	bouquet, err := svc.GetWholesomeBouquet(context.Background())

	assert.NoError(t, err)
	assert.NotNil(t, bouquet)
	assert.Equal(t, "You are a wonderful developer! Here is a bouquet of wholesome packages just for you:", bouquet.Message)
	assert.GreaterOrEqual(t, len(bouquet.Items), 1)
}

func TestVouchIncident(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockIncidentRepository(ctrl)
	mockPkgsite := NewMockPkgsiteService(ctrl)
	svc := NewIncidentService(mockRepo, mockPkgsite)

	ctx := context.Background()
	id := uint64(1)
	notes := "original notes"
	mockRepo.EXPECT().GetByID(ctx, id).Return(&domain.Incident{ID: id, Notes: &notes}, nil)
	mockRepo.EXPECT().Update(ctx, id, gomock.Any()).Do(func(ctx context.Context, id uint64, patch domain.IncidentPatch) {
		assert.Contains(t, *patch.Notes, "fellow Gopher has vouched")
	}).Return(nil)

	err := svc.VouchIncident(ctx, id)
	assert.NoError(t, err)
}

func TestCreateIncident_Milestones(t *testing.T) {
	tests := []struct {
		name          string
		totalExisting int
		expectedNote  string
	}{
		{
			name:          "First Incident",
			totalExisting: 0,
			expectedNote:  "This is your very first grievance",
		},
		{
			name:          "Fifth Incident",
			totalExisting: 4,
			expectedNote:  "Your 5th grievance",
		},
		{
			name:          "Tenth Incident",
			totalExisting: 9,
			expectedNote:  "Double digits! 10 grievances",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := NewMockIncidentRepository(ctrl)
			mockPkgsite := NewMockPkgsiteService(ctrl)
			svc := NewIncidentService(mockRepo, mockPkgsite)

			req := domain.IncidentCreate{
				ReporterID:                 "jordan",
				Subject:                    "test",
				Description:                "test",
				AssumedGoodIntentions:      true,
				PromisedToBeKindToYourself: true,
			}

			mockPkgsite.EXPECT().Search(gomock.Any(), "kindness").Return(&domain.PaginatedResponse[domain.SearchResult]{}, nil)
			mockPkgsite.EXPECT().Search(gomock.Any(), req.Subject).Return(&domain.PaginatedResponse[domain.SearchResult]{}, nil)

			mockRepo.EXPECT().
				List(gomock.Any(), domain.ListParams{ReporterID: req.ReporterID}).
				Return(domain.ListResult{Total: tt.totalExisting}, nil)

			mockRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Do(func(ctx context.Context, inc *domain.Incident) {
				assert.Contains(t, *inc.Notes, tt.expectedNote)
			}).Return(uint64(1), nil)

			mockRepo.EXPECT().GetByID(gomock.Any(), uint64(1)).Return(&domain.Incident{}, nil)

			_, err := svc.CreateIncident(context.Background(), req)
			assert.NoError(t, err)
		})
	}
}
