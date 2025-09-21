package domain

import (
	"time"

	"github.com/google/uuid"
)

type Subscription struct {
	ID           uuid.UUID
	ServiceName  string
	MonthlyPrice string
	UserID       uuid.UUID
	StartMonth   string
	EndMonth     *string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type CreateInput struct {
	ServiceName  string
	MonthlyPrice string
	UserID       uuid.UUID
	StartMonth   string
	EndMonth     *string
}

type UpdateInput struct {
	ServiceName  *string
	MonthlyPrice *string
	StartMonth   *string
	EndMonth     *string
}

type ListFilter struct {
	UserID      *uuid.UUID
	ServiceName *string
	Limit       int
	Offset      int
}
