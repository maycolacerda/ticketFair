-- migration/database-init.sql

-- ─────────────────────────────────────────────
-- TABLES
-- ─────────────────────────────────────────────

CREATE TABLE IF NOT EXISTS users (
    user_id    UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    email      VARCHAR(100) NOT NULL UNIQUE,
    password   VARCHAR(255) NOT NULL,
    username   VARCHAR(100) NOT NULL UNIQUE,
    active     BOOLEAN      NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ  NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS merchants (
    merchant_id UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    name        VARCHAR(100) NOT NULL,
    email       VARCHAR(100) NOT NULL UNIQUE,
    password    VARCHAR(255) NOT NULL,
    phone       VARCHAR(20)  NOT NULL,
    description TEXT,
    active      BOOLEAN      NOT NULL DEFAULT true,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT now(),
    deleted_at  TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS merchant_reps (
    merchant_rep_id UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    merchant_id     UUID         NOT NULL REFERENCES merchants(merchant_id) ON DELETE CASCADE,
    name            VARCHAR(100) NOT NULL,
    email           VARCHAR(100) NOT NULL UNIQUE,
    password        VARCHAR(255) NOT NULL,
    phone           VARCHAR(20)  NOT NULL,
    role            VARCHAR(20)  NOT NULL DEFAULT 'staff',
    active          BOOLEAN      NOT NULL DEFAULT true,
    created_at      TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ  NOT NULL DEFAULT now(),
    deleted_at      TIMESTAMPTZ,
    CONSTRAINT chk_merchant_rep_role CHECK (role IN ('admin', 'manager', 'staff'))
);

CREATE TABLE IF NOT EXISTS profiles (
    profile_id     UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id        UUID         NOT NULL UNIQUE REFERENCES users(user_id) ON DELETE CASCADE,
    first_name     VARCHAR(100) NOT NULL,
    last_name      VARCHAR(100) NOT NULL,
    phone_number   VARCHAR(20)  NOT NULL UNIQUE,
    verified_email BOOLEAN      NOT NULL DEFAULT false,
    verified_phone BOOLEAN      NOT NULL DEFAULT false,
    created_at     TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at     TIMESTAMPTZ  NOT NULL DEFAULT now(),
    deleted_at     TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS addresses (
    address_id UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    profile_id UUID         NOT NULL UNIQUE REFERENCES profiles(profile_id) ON DELETE CASCADE,
    street     VARCHAR(255) NOT NULL,
    city       VARCHAR(100) NOT NULL,
    state      VARCHAR(100) NOT NULL,
    country    CHAR(2)      NOT NULL DEFAULT 'BR',
    zip_code   VARCHAR(20)  NOT NULL,
    created_at TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ  NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ,
    CONSTRAINT chk_country_code CHECK (length(country) = 2)
);

CREATE TABLE IF NOT EXISTS events (
    event_id    UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    merchant_id UUID         NOT NULL REFERENCES merchants(merchant_id) ON DELETE CASCADE,
    name        VARCHAR(100) NOT NULL,
    description TEXT,
    location    VARCHAR(255) NOT NULL,
    start_time  TIMESTAMPTZ  NOT NULL,
    end_time    TIMESTAMPTZ  NOT NULL,
    capacity    INT          NOT NULL DEFAULT 0,
    active      BOOLEAN      NOT NULL DEFAULT true,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT now(),
    deleted_at  TIMESTAMPTZ,
    CONSTRAINT chk_event_times    CHECK (end_time > start_time),
    CONSTRAINT chk_event_capacity CHECK (capacity >= 0)
);

CREATE TABLE IF NOT EXISTS transactions (
    transaction_id UUID          PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id        UUID          NOT NULL REFERENCES users(user_id) ON DELETE RESTRICT,
    event_id       UUID          NOT NULL REFERENCES events(event_id) ON DELETE RESTRICT,
    amount         DECIMAL(10,2) NOT NULL,
    status         VARCHAR(20)   NOT NULL DEFAULT 'pending',
    created_at     TIMESTAMPTZ   NOT NULL DEFAULT now(),
    updated_at     TIMESTAMPTZ   NOT NULL DEFAULT now(),
    deleted_at     TIMESTAMPTZ,
    CONSTRAINT chk_transaction_status CHECK (status IN ('pending', 'completed', 'failed', 'refunded')),
    CONSTRAINT chk_transaction_amount CHECK (amount > 0)
);

-- ─────────────────────────────────────────────
-- INDEXES
-- ─────────────────────────────────────────────

CREATE INDEX IF NOT EXISTS idx_users_deleted_at          ON users(deleted_at);
CREATE INDEX IF NOT EXISTS idx_merchants_deleted_at      ON merchants(deleted_at);
CREATE INDEX IF NOT EXISTS idx_merchant_reps_deleted_at  ON merchant_reps(deleted_at);
CREATE INDEX IF NOT EXISTS idx_merchant_reps_merchant_id ON merchant_reps(merchant_id);
CREATE INDEX IF NOT EXISTS idx_profiles_deleted_at       ON profiles(deleted_at);
CREATE INDEX IF NOT EXISTS idx_events_deleted_at         ON events(deleted_at);
CREATE INDEX IF NOT EXISTS idx_events_merchant_id        ON events(merchant_id);
CREATE INDEX IF NOT EXISTS idx_events_start_time         ON events(start_time);
CREATE INDEX IF NOT EXISTS idx_transactions_user_id      ON transactions(user_id);
CREATE INDEX IF NOT EXISTS idx_transactions_event_id     ON transactions(event_id);
CREATE INDEX IF NOT EXISTS idx_transactions_status       ON transactions(status);

-- ─────────────────────────────────────────────
-- FUNCTIONS
-- ─────────────────────────────────────────────

CREATE OR REPLACE FUNCTION create_profile_with_address(
    p_user_id    UUID,
    p_first_name VARCHAR,
    p_last_name  VARCHAR,
    p_phone      VARCHAR,
    p_street     VARCHAR,
    p_city       VARCHAR,
    p_state      VARCHAR,
    p_country    CHAR(2),
    p_zip_code   VARCHAR
) RETURNS UUID AS $$
DECLARE
    v_profile_id UUID;
BEGIN
    INSERT INTO profiles (user_id, first_name, last_name, phone_number)
    VALUES (p_user_id, p_first_name, p_last_name, p_phone)
    RETURNING profile_id INTO v_profile_id;

    INSERT INTO addresses (profile_id, street, city, state, country, zip_code)
    VALUES (v_profile_id, p_street, p_city, p_state, p_country, p_zip_code);

    RETURN v_profile_id;
END;
$$ LANGUAGE PLpgSQL;

CREATE OR REPLACE FUNCTION purchase_ticket(
    p_user_id  UUID,
    p_event_id UUID,
    p_amount   DECIMAL
) RETURNS UUID AS $$
DECLARE
    v_transaction_id UUID;
    v_capacity       INT;
BEGIN
    SELECT capacity INTO v_capacity
    FROM events
    WHERE event_id = p_event_id
      AND active = true
      AND deleted_at IS NULL
    FOR UPDATE;

    IF v_capacity IS NULL THEN
        RAISE EXCEPTION 'event_not_found';
    END IF;

    IF v_capacity <= 0 THEN
        RAISE EXCEPTION 'event_sold_out';
    END IF;

    UPDATE events
    SET capacity   = capacity - 1,
        updated_at = now()
    WHERE event_id = p_event_id;

    INSERT INTO transactions (user_id, event_id, amount, status)
    VALUES (p_user_id, p_event_id, p_amount, 'completed')
    RETURNING transaction_id INTO v_transaction_id;

    RETURN v_transaction_id;
END;
$$ LANGUAGE PLpgSQL;

CREATE OR REPLACE FUNCTION refund_ticket(
    p_transaction_id UUID
) RETURNS VOID AS $$
DECLARE
    v_event_id UUID;
    v_status   VARCHAR;
BEGIN
    SELECT event_id, status INTO v_event_id, v_status
    FROM transactions
    WHERE transaction_id = p_transaction_id
      AND deleted_at IS NULL
    FOR UPDATE;

    IF v_event_id IS NULL THEN
        RAISE EXCEPTION 'transaction_not_found';
    END IF;

    IF v_status != 'completed' THEN
        RAISE EXCEPTION 'transaction_not_refundable';
    END IF;

    UPDATE events
    SET capacity   = capacity + 1,
        updated_at = now()
    WHERE event_id = v_event_id;

    UPDATE transactions
    SET status     = 'refunded',
        updated_at = now()
    WHERE transaction_id = p_transaction_id;
END;
$$ LANGUAGE PLpgSQL;