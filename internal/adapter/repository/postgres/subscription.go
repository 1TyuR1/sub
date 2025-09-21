package postgres

import (
	"context"
	"strconv"
	"strings"
	"time"

	"crud_ef/internal/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SubscriptionRepo struct {
	pool *pgxpool.Pool
}

func NewSubscriptionRepo(pool *pgxpool.Pool) *SubscriptionRepo {
	return &SubscriptionRepo{pool: pool}
}

func (r *SubscriptionRepo) Create(ctx context.Context, in domain.CreateInput) (domain.Subscription, error) {
	var s domain.Subscription
	var end any
	if in.EndMonth != nil {
		t, _ := time.Parse("2006-01", *in.EndMonth)
		end = time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.UTC)
	} else {
		end = nil
	}
	startT, _ := time.Parse("2006-01", in.StartMonth)
	start := time.Date(startT.Year(), startT.Month(), 1, 0, 0, 0, 0, time.UTC)

	q := `
INSERT INTO subscriptions (service_name, monthly_price, user_id, start_month, end_month)
VALUES ($1, $2::numeric(12,2), $3, $4, $5)
RETURNING id, service_name, monthly_price::text, user_id,
         to_char(start_month, 'YYYY-MM') AS start_month,
         CASE WHEN end_month IS NULL THEN NULL ELSE to_char(end_month, 'YYYY-MM') END AS end_month,
         created_at, updated_at;
`
	err := r.pool.QueryRow(ctx, q, in.ServiceName, in.MonthlyPrice, in.UserID, start, end).Scan(
		&s.ID, &s.ServiceName, &s.MonthlyPrice, &s.UserID, &s.StartMonth, &s.EndMonth, &s.CreatedAt, &s.UpdatedAt,
	)
	return s, err
}

func (r *SubscriptionRepo) Get(ctx context.Context, id uuid.UUID) (domain.Subscription, error) {
	q := `
SELECT id, service_name, monthly_price::text, user_id,
       to_char(start_month, 'YYYY-MM') AS start_month,
       CASE WHEN end_month IS NULL THEN NULL ELSE to_char(end_month, 'YYYY-MM') END AS end_month,
       created_at, updated_at
FROM subscriptions WHERE id = $1;
`
	var s domain.Subscription
	err := r.pool.QueryRow(ctx, q, id).Scan(
		&s.ID, &s.ServiceName, &s.MonthlyPrice, &s.UserID, &s.StartMonth, &s.EndMonth, &s.CreatedAt, &s.UpdatedAt,
	)
	return s, err
}

func (r *SubscriptionRepo) List(ctx context.Context, f domain.ListFilter) ([]domain.Subscription, error) {
	var args []any
	var whr []string
	idx := 1
	if f.UserID != nil {
		whr = append(whr, "user_id = $"+strconv.Itoa(idx))
		args = append(args, *f.UserID)
		idx++
	}
	if f.ServiceName != nil && *f.ServiceName != "" {
		whr = append(whr, "service_name ILIKE $"+strconv.Itoa(idx))
		args = append(args, "%"+*f.ServiceName+"%")
		idx++
	}
	where := ""
	if len(whr) > 0 {
		where = "WHERE " + strings.Join(whr, " AND ")
	}

	q := `
SELECT id, service_name, monthly_price::text, user_id,
       to_char(start_month, 'YYYY-MM') AS start_month,
       CASE WHEN end_month IS NULL THEN NULL ELSE to_char(end_month, 'YYYY-MM') END AS end_month,
       created_at, updated_at
FROM subscriptions
` + where + `
ORDER BY created_at DESC
LIMIT $` + strconv.Itoa(idx) + ` OFFSET $` + strconv.Itoa(idx+1) + `;
`
	args = append(args, f.Limit, f.Offset)

	rows, err := r.pool.Query(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []domain.Subscription
	for rows.Next() {
		var s domain.Subscription
		if err := rows.Scan(&s.ID, &s.ServiceName, &s.MonthlyPrice, &s.UserID, &s.StartMonth, &s.EndMonth, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, s)
	}
	return out, rows.Err()
}

func (r *SubscriptionRepo) Update(ctx context.Context, id uuid.UUID, in domain.UpdateInput) (domain.Subscription, error) {
	set := []string{}
	args := []any{}
	i := 1

	if in.ServiceName != nil {
		set = append(set, "service_name = $"+strconv.Itoa(i))
		args = append(args, *in.ServiceName)
		i++
	}
	if in.MonthlyPrice != nil {
		set = append(set, "monthly_price = $"+strconv.Itoa(i)+"::numeric(12,2)")
		args = append(args, *in.MonthlyPrice)
		i++
	}
	if in.StartMonth != nil {
		smT, _ := time.Parse("2006-01", *in.StartMonth)
		sm := time.Date(smT.Year(), smT.Month(), 1, 0, 0, 0, 0, time.UTC)
		set = append(set, "start_month = $"+strconv.Itoa(i))
		args = append(args, sm)
		i++
	}
	if in.EndMonth != nil {
		if *in.EndMonth == "" {
			set = append(set, "end_month = NULL")
		} else {
			emT, _ := time.Parse("2006-01", *in.EndMonth)
			em := time.Date(emT.Year(), emT.Month(), 1, 0, 0, 0, 0, time.UTC)
			set = append(set, "end_month = $"+strconv.Itoa(i))
			args = append(args, em)
			i++
		}
	}
	set = append(set, "updated_at = now()")
	args = append(args, id)

	q := `
UPDATE subscriptions
SET ` + strings.Join(set, ", ") + `
WHERE id = $` + strconv.Itoa(i) + `
RETURNING id, service_name, monthly_price::text, user_id,
          to_char(start_month, 'YYYY-MM') AS start_month,
          CASE WHEN end_month IS NULL THEN NULL ELSE to_char(end_month, 'YYYY-MM') END AS end_month,
          created_at, updated_at;
`
	var s domain.Subscription
	err := r.pool.QueryRow(ctx, q, args...).Scan(
		&s.ID, &s.ServiceName, &s.MonthlyPrice, &s.UserID, &s.StartMonth, &s.EndMonth, &s.CreatedAt, &s.UpdatedAt,
	)
	return s, err
}

func (r *SubscriptionRepo) Delete(ctx context.Context, id uuid.UUID) (bool, error) {
	cmd, err := r.pool.Exec(ctx, `DELETE FROM subscriptions WHERE id = $1`, id)
	if err != nil {
		return false, err
	}
	return cmd.RowsAffected() > 0, nil
}

func (r *SubscriptionRepo) Total(ctx context.Context, from, to time.Time, userID *uuid.UUID, service *string) (string, error) {
	q := `
WITH months AS (
  SELECT generate_series($1::date, $2::date, interval '1 month')::date AS m
)
SELECT COALESCE(TO_CHAR(SUM(s.monthly_price)::numeric(12,2), 'FM9999999990D00'), '0.00')
FROM months mo
JOIN subscriptions s
  ON s.start_month <= mo.m
 AND (s.end_month IS NULL OR s.end_month >= mo.m)
WHERE ($3::uuid IS NULL OR s.user_id = $3)
  AND ($4::text IS NULL OR s.service_name ILIKE '%'||$4||'%');
`
	var userArg any = nil
	if userID != nil {
		userArg = *userID
	}
	var srvArg any = nil
	if service != nil {
		srvArg = *service
	}
	var total string
	err := r.pool.QueryRow(ctx, q, from, to, userArg, srvArg).Scan(&total)
	return total, err
}
