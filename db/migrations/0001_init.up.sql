CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS subscriptions (
    id             uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    service_name   text NOT NULL,
    monthly_price  numeric(12,2) NOT NULL CHECK (monthly_price >= 0),
    user_id        uuid NOT NULL,
    start_month    date NOT NULL CHECK (date_trunc('month', start_month) = start_month),
    end_month      date NULL CHECK (
        (end_month IS NULL) OR
        (date_trunc('month', end_month) = end_month AND end_month >= start_month)
    ),
    created_at     timestamptz NOT NULL DEFAULT now(),
    updated_at     timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_subscriptions_user ON subscriptions(user_id);
CREATE INDEX IF NOT EXISTS idx_subscriptions_service ON subscriptions(service_name);
CREATE INDEX IF NOT EXISTS idx_subscriptions_period ON subscriptions(start_month, end_month);
