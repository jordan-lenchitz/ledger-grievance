package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"github.com/jordan-lenchitz/ledger-grievance/go-app/internal/domain"
)

type mysqlIncidentRepository struct {
	db *sql.DB
}

func NewMySQLIncidentRepository(db *sql.DB) domain.IncidentRepository {
	return &mysqlIncidentRepository{db: db}
}

func (r *mysqlIncidentRepository) Create(ctx context.Context, inc *domain.Incident) (uint64, error) {
	res, err := r.db.ExecContext(ctx, `
		INSERT INTO incidents 
		(reporter_id, occurred_at, subject, category, severity, description, evidence_uri, notes, requires_accommodation)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		inc.ReporterID, inc.OccurredAt, inc.Subject, inc.Category, inc.Severity, inc.Description, inc.EvidenceURI, inc.Notes, inc.RequiresAccommodation,
	)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return uint64(id), nil
}

func (r *mysqlIncidentRepository) GetByID(ctx context.Context, id uint64) (*domain.Incident, error) {
	var inc domain.Incident
	err := r.db.QueryRowContext(ctx, "SELECT id, reporter_id, occurred_at, recorded_at, subject, category, severity, description, evidence_uri, requires_accommodation, status, notes FROM incidents WHERE id = ?", id).
		Scan(&inc.ID, &inc.ReporterID, &inc.OccurredAt, &inc.RecordedAt, &inc.Subject, &inc.Category, &inc.Severity, &inc.Description, &inc.EvidenceURI, &inc.RequiresAccommodation, &inc.Status, &inc.Notes)
	if err != nil {
		return nil, err
	}
	return &inc, nil
}

func (r *mysqlIncidentRepository) List(ctx context.Context, params domain.ListParams) (domain.ListResult, error) {
	var clauses []string
	var args []interface{}

	if params.ReporterID != "" {
		clauses = append(clauses, "reporter_id = ?")
		args = append(args, params.ReporterID)
	}
	if params.Status != "" {
		clauses = append(clauses, "status = ?")
		args = append(args, params.Status)
	}
	if params.Category != "" {
		clauses = append(clauses, "category = ?")
		args = append(args, params.Category)
	}
	if params.Query != "" {
		clauses = append(clauses, "(subject LIKE ? OR description LIKE ? OR notes LIKE ?)")
		like := "%" + params.Query + "%"
		args = append(args, like, like, like)
	}

	where := ""
	if len(clauses) > 0 {
		where = " WHERE " + strings.Join(clauses, " AND ")
	}

	countQuery := "SELECT COUNT(*) FROM incidents" + where
	var total int
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return domain.ListResult{}, err
	}

	query := fmt.Sprintf("SELECT id, reporter_id, occurred_at, recorded_at, subject, category, severity, description, evidence_uri, requires_accommodation, status, notes FROM incidents %s ORDER BY recorded_at DESC, id DESC LIMIT ? OFFSET ?", where)
	rows, err := r.db.QueryContext(ctx, query, append(args, params.Limit, params.Offset)...)
	if err != nil {
		return domain.ListResult{}, err
	}
	defer rows.Close()

	incidents := []domain.Incident{}
	for rows.Next() {
		var inc domain.Incident
		err := rows.Scan(&inc.ID, &inc.ReporterID, &inc.OccurredAt, &inc.RecordedAt, &inc.Subject, &inc.Category, &inc.Severity, &inc.Description, &inc.EvidenceURI, &inc.RequiresAccommodation, &inc.Status, &inc.Notes)
		if err != nil {
			return domain.ListResult{}, err
		}
		incidents = append(incidents, inc)
	}

	return domain.ListResult{
		Data:  incidents,
		Total: total,
	}, nil
}

func (r *mysqlIncidentRepository) Update(ctx context.Context, id uint64, patch domain.IncidentPatch) error {
	var assignments []string
	var args []interface{}

	if patch.Status != nil {
		assignments = append(assignments, "status = ?")
		args = append(args, *patch.Status)
	}
	if patch.Notes != nil {
		assignments = append(assignments, "notes = ?")
		args = append(args, *patch.Notes)
	}
	if patch.Category != nil {
		assignments = append(assignments, "category = ?")
		args = append(args, *patch.Category)
	}
	if patch.Severity != nil {
		assignments = append(assignments, "severity = ?")
		args = append(args, *patch.Severity)
	}
	if patch.EvidenceURI != nil {
		assignments = append(assignments, "evidence_uri = ?")
		args = append(args, *patch.EvidenceURI)
	}

	if len(assignments) == 0 {
		return nil
	}

	query := fmt.Sprintf("UPDATE incidents SET %s WHERE id = ?", strings.Join(assignments, ", "))
	args = append(args, id)

	res, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *mysqlIncidentRepository) Archive(ctx context.Context, id uint64) error {
	res, err := r.db.ExecContext(ctx, "UPDATE incidents SET status = 'archived' WHERE id = ?", id)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}
