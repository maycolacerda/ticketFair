-- migration/database-init.sql

CREATE DATABASE IF NOT EXISTS ticketfair;
USE ticketfair;
-- ─────────────────────────────────────────────
-- TABLES
-- ─────────────────────────────────────────────

CREATE TABLE IF NOT EXISTS users (
    user_id    UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    email      VARCHAR(100) NOT NULL UNIQUE,
    CONSTRAINT uni_users_email UNIQUE(email), 
    password   VARCHAR(255) NOT NULL,
    username   VARCHAR(100) NOT NULL,
    CONSTRAINT uni_users_username UNIQUE(username),
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
    profile_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    UUID NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    CONSTRAINT uni_profiles_user_id UNIQUE(user_id),
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
CREATE TABLE IF NOT EXISTS tickets (
    ticket_id      UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    transaction_id UUID        NOT NULL REFERENCES transactions(transaction_id) ON DELETE RESTRICT,
    user_id        UUID        NOT NULL REFERENCES users(user_id) ON DELETE RESTRICT,
    event_id       UUID        NOT NULL REFERENCES events(event_id) ON DELETE RESTRICT,
    status         VARCHAR(20) NOT NULL DEFAULT 'active',
    created_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at     TIMESTAMPTZ,
    CONSTRAINT chk_ticket_status CHECK (status IN ('active', 'used', 'refunded', 'cancelled'))
);

CREATE TABLE IF NOT EXISTS verifications (
    verification_id UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID        NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    type            VARCHAR(10) NOT NULL,
    code            VARCHAR(10) NOT NULL,
    expires_at      TIMESTAMPTZ NOT NULL,
    used_at         TIMESTAMPTZ,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at      TIMESTAMPTZ,
    CONSTRAINT chk_verification_type CHECK (type IN ('email', 'phone'))
);

CREATE TABLE IF NOT EXISTS admins (
    admin_id   UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    name       VARCHAR(100) NOT NULL,
    email      VARCHAR(100) NOT NULL UNIQUE,
    password   VARCHAR(255) NOT NULL,
    active     BOOLEAN      NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ  NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
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
CREATE INDEX IF NOT EXISTS idx_tickets_user_id        ON tickets(user_id);
CREATE INDEX IF NOT EXISTS idx_tickets_event_id       ON tickets(event_id);
CREATE INDEX IF NOT EXISTS idx_tickets_transaction_id ON tickets(transaction_id);
CREATE INDEX IF NOT EXISTS idx_tickets_status         ON tickets(status);
CREATE INDEX IF NOT EXISTS idx_verifications_user_id ON verifications(user_id);
CREATE INDEX IF NOT EXISTS idx_verifications_type    ON verifications(type);
CREATE INDEX IF NOT EXISTS idx_admins_deleted_at ON admins(deleted_at);


-- ─────────────────────────────────────────────
-- SEED DATA
-- password for all accounts is: PassW0rd!
-- bcrypt hash generated with cost 10
-- ─────────────────────────────────────────────

INSERT INTO merchants (
    merchant_id,
    name,
    email,
    password,
    phone,
    description,
    active
) VALUES (
    'a1b2c3d4-e5f6-7890-abcd-ef1234567890',
    'TicketFair Produções',
    'contato@ticketfairprod.com',
    '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 
    '44999000000',
    'Produtora oficial de eventos TicketFair.',
    true
) ON CONFLICT (email) DO NOTHING;

INSERT INTO merchant_reps (
    merchant_rep_id,
    merchant_id,
    name,
    email,
    password,
    phone,
    role,
    active
) VALUES (
    'b2c3d4e5-f6a7-8901-bcde-f12345678901',
    'a1b2c3d4-e5f6-7890-abcd-ef1234567890',
    'Carlos Admin',
    'carlos@ticketfairprod.com',
    '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 
    '44999000001',
    'admin',
    true
) ON CONFLICT (email) DO NOTHING;

INSERT INTO events (
    event_id,
    merchant_id,
    name,
    description,
    location,
    start_time,
    end_time,
    capacity,
    active
) VALUES (
    'c3d4e5f6-a7b8-9012-cdef-123456789012',
    'a1b2c3d4-e5f6-7890-abcd-ef1234567890',
    'Festival de Verão 2026',
    'O maior festival de música do Paraná.',
    'Parque de Exposições de Cianorte — PR',
    '2026-12-20T18:00:00Z',
    '2026-12-20T23:59:00Z',
    1000,
    true
) ON CONFLICT DO NOTHING;

INSERT INTO admins (
    admin_id,
    name,
    email,
    password,
    active
) VALUES (
    'd4e5f6a7-b8c9-0123-defa-234567890123',
    'Super Admin',
    'admin@ticketfair.com',
    '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi',
    true
) ON CONFLICT (email) DO NOTHING;


-- ─────────────────────────────────────────────
-- FUNCTIONS
-- ─────────────────────────────────────────────

CREATE OR REPLACE FUNCTION create_profile_with_address(
    p_user_id     UUID,
    p_first_name  VARCHAR,
    p_last_name   VARCHAR,
    p_phone       VARCHAR,
    p_street      VARCHAR,
    p_city        VARCHAR,
    p_state       VARCHAR,
    p_country     CHAR(2),
    p_zip_code    VARCHAR
) RETURNS UUID AS $$
  WITH new_profile AS (
    INSERT INTO profiles (user_id, first_name, last_name, phone_number)
    VALUES (p_user_id, p_first_name, p_last_name, p_phone)
    RETURNING profile_id
  ),
  new_address AS (
    INSERT INTO addresses (profile_id, street, city, state, country, zip_code)
    SELECT profile_id, p_street, p_city, p_state, p_country, p_zip_code FROM new_profile
    RETURNING address_id 
  ) -- <--- This parenthesis was missing!
  SELECT profile_id FROM new_profile;
$$ LANGUAGE SQL;

-- 2. Purchase Ticket (Optimized for Distributed SQL)
CREATE OR REPLACE FUNCTION purchase_ticket(
    p_user_id  UUID,
    p_event_id UUID,
    p_amount   DECIMAL
) RETURNS UUID AS $$
  WITH updated_event AS (
    UPDATE events
    SET capacity = capacity - 1, updated_at = now()
    WHERE event_id = p_event_id 
      AND active = true 
      AND capacity > 0 
      AND deleted_at IS NULL
    RETURNING event_id
  )
  INSERT INTO transactions (user_id, event_id, amount, status)
  SELECT p_user_id, event_id, p_amount, 'completed'
  FROM updated_event
  RETURNING transaction_id;
$$ LANGUAGE SQL;


CREATE OR REPLACE FUNCTION refund_ticket(
    p_transaction_id UUID
) RETURNS UUID AS $$
  WITH refunded_tx AS (
    UPDATE transactions
    SET status = 'refunded', updated_at = now()
    WHERE transaction_id = p_transaction_id 
      AND status = 'completed'
      AND deleted_at IS NULL
    RETURNING event_id, transaction_id
  ),
  restore_capacity AS (
    UPDATE events AS e
    SET capacity = e.capacity + 1, updated_at = now()
    FROM refunded_tx AS r
    WHERE e.event_id = r.event_id
    RETURNING e.event_id 
  )
  SELECT transaction_id FROM refunded_tx;
$$ LANGUAGE SQL;