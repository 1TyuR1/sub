package subscription

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"time"

	"crud_ef/internal/domain"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, in domain.CreateInput) (domain.Subscription, error)
	Get(ctx context.Context, id uuid.UUID) (domain.Subscription, error)
	List(ctx context.Context, f domain.ListFilter) ([]domain.Subscription, error)
	Update(ctx context.Context, id uuid.UUID, in domain.UpdateInput) (domain.Subscription, error)
	Delete(ctx context.Context, id uuid.UUID) (bool, error)
	Total(ctx context.Context, from, to time.Time, userID *uuid.UUID, service *string) (string, error)
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func parseMonth(s string) (time.Time, error) {
	t, err := time.Parse("2006-01", s)
	if err != nil {
		return time.Time{}, err
	}
	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.UTC), nil
}

func validPrice(p string) bool {
	if p == "" {
		return false
	}
	if !strings.Contains(p, ".") {
		_, err := strconv.Atoi(p)
		return err == nil
	}
	parts := strings.Split(p, ".")
	if len(parts) != 2 {
		return false
	}
	if _, err := strconv.Atoi(parts[0]); err != nil {
		return false
	}
	if len(parts[1]) == 0 || len(parts[1]) > 2 {
		return false
	}
	_, err := strconv.Atoi(parts[1])
	return err == nil
}

func (s *Service) Create(ctx context.Context, in domain.CreateInput) (domain.Subscription, error) {
	if strings.TrimSpace(in.ServiceName) == "" || !validPrice(in.MonthlyPrice) {
		return domain.Subscription{}, errors.New("invalid service_name or monthly_price")
	}
	if _, err := parseMonth(in.StartMonth); err != nil {
		return domain.Subscription{}, errors.New("invalid start_month (YYYY-MM)")
	}
	if in.EndMonth != nil {
		if _, err := parseMonth(*in.EndMonth); err != nil {
			return domain.Subscription{}, errors.New("invalid end_month (YYYY-MM)")
		}
	}
	return s.repo.Create(ctx, in)
}

func (s *Service) Get(ctx context.Context, id uuid.UUID) (domain.Subscription, error) {
	return s.repo.Get(ctx, id)
}

func (s *Service) List(ctx context.Context, f domain.ListFilter) ([]domain.Subscription, error) {
	return s.repo.List(ctx, f)
}

func (s *Service) Update(ctx context.Context, id uuid.UUID, in domain.UpdateInput) (domain.Subscription, error) {
	if in.MonthlyPrice != nil && !validPrice(*in.MonthlyPrice) {
		return domain.Subscription{}, errors.New("invalid monthly_price")
	}
	if in.StartMonth != nil {
		if _, err := parseMonth(*in.StartMonth); err != nil {
			return domain.Subscription{}, errors.New("invalid start_month (YYYY-MM)")
		}
	}
	if in.EndMonth != nil && *in.EndMonth != "" {
		if _, err := parseMonth(*in.EndMonth); err != nil {
			return domain.Subscription{}, errors.New("invalid end_month (YYYY-MM)")
		}
	}
	return s.repo.Update(ctx, id, in)
}

func (s *Service) Delete(ctx context.Context, id uuid.UUID) (bool, error) {
	return s.repo.Delete(ctx, id)
}

func (s *Service) Total(ctx context.Context, fromStr, toStr string, userID *uuid.UUID, service *string) (string, error) {
	from, err := parseMonth(fromStr)
	if err != nil {
		return "", errors.New("invalid from (YYYY-MM)")
	}
	to, err := parseMonth(toStr)
	if err != nil {
		return "", errors.New("invalid to (YYYY-MM)")
	}
	if to.Before(from) {
		return "", errors.New("to must be >= from")
	}
	return s.repo.Total(ctx, from, to, userID, service)
}
