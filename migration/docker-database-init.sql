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
    image_url   TEXT,
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
    ticket_type_id UUID REFERENCES ticket_types(ticket_type_id) ON DELETE RESTRICT,
    quantity INT NOT NULL DEFAULT 1,
    status         VARCHAR(20)   NOT NULL DEFAULT 'pending',
    created_at     TIMESTAMPTZ   NOT NULL DEFAULT now(),
    updated_at     TIMESTAMPTZ   NOT NULL DEFAULT now(),
    deleted_at     TIMESTAMPTZ,
    CONSTRAINT chk_transaction_status CHECK (status IN ('pending', 'completed', 'failed', 'refunded')),
    CONSTRAINT chk_transaction_amount CHECK (amount > 0)
);

CREATE TABLE IF NOT EXISTS ticket_types (
    ticket_type_id  UUID          PRIMARY KEY DEFAULT gen_random_uuid(),
    event_id        UUID          NOT NULL REFERENCES events(event_id) ON DELETE CASCADE,
    name            VARCHAR(100)  NOT NULL,
    description     TEXT,
    category        VARCHAR(50)   NOT NULL DEFAULT 'general',
    price_cents     BIGINT        NOT NULL,
    capacity        INT           NOT NULL,
    available       INT           NOT NULL,
    min_per_order   INT           NOT NULL DEFAULT 1,
    max_per_order   INT           NOT NULL DEFAULT 10,
    sale_starts_at  TIMESTAMPTZ,
    sale_ends_at    TIMESTAMPTZ,
    active          BOOLEAN       NOT NULL DEFAULT true,
    sort_order      INT           NOT NULL DEFAULT 0,
    created_at      TIMESTAMPTZ   NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ   NOT NULL DEFAULT now(),
    deleted_at      TIMESTAMPTZ,
    CONSTRAINT chk_ticket_type_category CHECK (
        category IN (
            'general',
            'vip',
            'early_bird',
            'reserved',
            'group',
            'day_pass',
            'tiered',
            'complimentary',
            'demographic'
        )
    ),
    CONSTRAINT chk_ticket_type_price     CHECK (price_cents >= 0),
    CONSTRAINT chk_ticket_type_capacity  CHECK (capacity > 0),
    CONSTRAINT chk_ticket_type_available CHECK (available >= 0 AND available <= capacity),
    CONSTRAINT chk_ticket_type_order     CHECK (min_per_order >= 1 AND max_per_order >= min_per_order),
    CONSTRAINT chk_ticket_type_sale_window CHECK (
        sale_ends_at IS NULL OR sale_starts_at IS NULL OR sale_ends_at > sale_starts_at
    )
);

CREATE TABLE IF NOT EXISTS tickets (
    ticket_id      UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    transaction_id UUID        NOT NULL REFERENCES transactions(transaction_id) ON DELETE RESTRICT,
    user_id        UUID        NOT NULL REFERENCES users(user_id) ON DELETE RESTRICT,
    event_id       UUID        NOT NULL REFERENCES events(event_id) ON DELETE RESTRICT,
    ticket_type_id UUID REFERENCES ticket_types(ticket_type_id) ON DELETE RESTRICT,
    ticket_type_name VARCHAR(100),
    price_paid_cents BIGINT NOT NULL DEFAULT 0,
    status         VARCHAR(20) NOT NULL DEFAULT 'active',
    created_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at     TIMESTAMPTZ,
    CONSTRAINT chk_ticket_status CHECK (status IN ('active', 'used', 'refunded', 'cancelled'))
);


CREATE TABLE IF NOT EXISTS payments (
    payment_id        UUID          PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id           UUID          NOT NULL REFERENCES users(user_id) ON DELETE RESTRICT,
    event_id          UUID          NOT NULL REFERENCES events(event_id) ON DELETE RESTRICT,
    transaction_id    UUID,         -- set after Stripe payment succeeds
    stripe_payment_id VARCHAR(255)  NOT NULL UNIQUE,
    amount            BIGINT        NOT NULL,      -- in cents
    currency          VARCHAR(10)   NOT NULL DEFAULT 'brl',
    quantity          INT           NOT NULL DEFAULT 1,
    status            VARCHAR(20)   NOT NULL DEFAULT 'pending',
    created_at        TIMESTAMPTZ   NOT NULL DEFAULT now(),
    updated_at        TIMESTAMPTZ   NOT NULL DEFAULT now(),
    deleted_at        TIMESTAMPTZ,
    CONSTRAINT chk_payment_status CHECK (
        status IN ('pending', 'succeeded', 'failed', 'canceled', 'refunded')
    ),
    CONSTRAINT chk_payment_amount CHECK (amount > 0)
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

CREATE TABLE IF NOT EXISTS password_resets (
    reset_id    UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     UUID        NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    code        VARCHAR(10) NOT NULL,
    expires_at  TIMESTAMPTZ NOT NULL,
    used_at     TIMESTAMPTZ,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at  TIMESTAMPTZ
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
CREATE INDEX IF NOT EXISTS idx_ticket_types_event_id   ON ticket_types(event_id);
CREATE INDEX IF NOT EXISTS idx_ticket_types_category   ON ticket_types(category);
CREATE INDEX IF NOT EXISTS idx_ticket_types_active     ON ticket_types(active);
CREATE INDEX IF NOT EXISTS idx_ticket_types_deleted_at ON ticket_types(deleted_at);
CREATE INDEX IF NOT EXISTS idx_verifications_user_id ON verifications(user_id);
CREATE INDEX IF NOT EXISTS idx_verifications_type    ON verifications(type);
CREATE INDEX IF NOT EXISTS idx_admins_deleted_at ON admins(deleted_at);
CREATE INDEX IF NOT EXISTS idx_payments_user_id           ON payments(user_id);
CREATE INDEX IF NOT EXISTS idx_payments_event_id          ON payments(event_id);
CREATE INDEX IF NOT EXISTS idx_payments_transaction_id    ON payments(transaction_id);
CREATE INDEX IF NOT EXISTS idx_payments_stripe_payment_id ON payments(stripe_payment_id);
CREATE INDEX IF NOT EXISTS idx_payments_status            ON payments(status);
CREATE INDEX IF NOT EXISTS idx_payments_deleted_at        ON payments(deleted_at);
CREATE INDEX IF NOT EXISTS idx_password_resets_user_id    ON password_resets(user_id);
CREATE INDEX IF NOT EXISTS idx_password_resets_deleted_at ON password_resets(deleted_at);


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
    'kP9vL2!Z&mR5*xN9?uW6@bY1_jT4qS',
    true
) ON CONFLICT (email) DO NOTHING;

INSERT INTO ticket_types (
    ticket_type_id,
    event_id,
    name,
    description,
    category,
    price_cents,
    capacity,
    available,
    min_per_order,
    max_per_order,
    sale_starts_at,
    sale_ends_at,
    sort_order
) VALUES
(
    'e1f2a3b4-c5d6-7890-ef12-345678901234',
    'c3d4e5f6-a7b8-9012-cdef-123456789012',
    'Early Bird',
    'Limited early access tickets at a discounted price.',
    'early_bird',
    2500,   -- R$ 25,00
    100,
    100,
    1,
    4,
    now(),
    '2026-10-01T00:00:00Z',
    1
),
(
    'f2a3b4c5-d6e7-8901-f012-456789012345',
    'c3d4e5f6-a7b8-9012-cdef-123456789012',
    'General Admission',
    'Standard access to all event areas.',
    'general',
    5000,   -- R$ 50,00
    700,
    700,
    1,
    10,
    now(),
    '2026-12-19T23:59:00Z',
    2
),
(
    'a3b4c5d6-e7f8-9012-0123-567890123456',
    'c3d4e5f6-a7b8-9012-cdef-123456789012',
    'VIP',
    'Priority entry, exclusive lounge access and meet-and-greet.',
    'vip',
    15000,  -- R$ 150,00
    150,
    150,
    1,
    4,
    now(),
    '2026-12-19T23:59:00Z',
    3
),
(
    'b4c5d6e7-f8a9-0123-1234-678901234567',
    'c3d4e5f6-a7b8-9012-cdef-123456789012',
    'Group Pack (4 tickets)',
    'Buy 4 General Admission tickets and save 20%.',
    'group',
    16000,  -- R$ 160,00 for 4 (vs R$ 200,00)
    50,
    50,
    1,
    1,    -- sold as a pack, quantity = 1 pack
    now(),
    '2026-12-19T23:59:00Z',
    4
)
ON CONFLICT DO NOTHING;
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
  ) 
  SELECT profile_id FROM new_profile;
$$ LANGUAGE SQL;


CREATE OR REPLACE FUNCTION purchase_ticket(
    p_user_id        UUID,
    p_event_id       UUID,
    p_ticket_type_id UUID,
    p_quantity       INT,
    p_amount         DECIMAL
) RETURNS UUID AS $$
DECLARE
    v_transaction_id UUID;
    v_available      INT;
    v_tt             RECORD;
    v_now            TIMESTAMPTZ := now();
BEGIN
    -- Lock ticket type row
    SELECT * INTO v_tt
    FROM ticket_types
    WHERE ticket_type_id = p_ticket_type_id
      AND event_id       = p_event_id
      AND active         = true
      AND deleted_at     IS NULL
    FOR UPDATE;

    IF NOT FOUND THEN
        RAISE EXCEPTION 'ticket_type_not_found';
    END IF;

    -- Check sale window
    IF v_tt.sale_starts_at IS NOT NULL AND v_now < v_tt.sale_starts_at THEN
        RAISE EXCEPTION 'ticket_sale_not_started';
    END IF;

    IF v_tt.sale_ends_at IS NOT NULL AND v_now > v_tt.sale_ends_at THEN
        RAISE EXCEPTION 'ticket_sale_ended';
    END IF;

    -- Check quantity bounds
    IF p_quantity < v_tt.min_per_order THEN
        RAISE EXCEPTION 'ticket_below_minimum';
    END IF;

    IF p_quantity > v_tt.max_per_order THEN
        RAISE EXCEPTION 'ticket_exceeds_maximum';
    END IF;

    -- Check availability
    IF v_tt.available < p_quantity THEN
        RAISE EXCEPTION 'ticket_type_sold_out';
    END IF;

    -- Decrement availability on ticket type
    UPDATE ticket_types
    SET available  = available - p_quantity,
        updated_at = now()
    WHERE ticket_type_id = p_ticket_type_id;

    -- Also decrement event capacity
    UPDATE events
    SET capacity   = capacity - p_quantity,
        updated_at = now()
    WHERE event_id = p_event_id;

    -- Insert transaction
    INSERT INTO transactions (user_id, event_id, ticket_type_id, amount, quantity, status)
    VALUES (p_user_id, p_event_id, p_ticket_type_id, p_amount, p_quantity, 'completed')
    RETURNING transaction_id INTO v_transaction_id;

    RETURN v_transaction_id;
END;
$$ LANGUAGE plpgsql;



CREATE OR REPLACE FUNCTION refund_ticket(
    p_transaction_id UUID
) RETURNS VOID AS $$
DECLARE
    v_tx RECORD;
BEGIN
    SELECT * INTO v_tx
    FROM transactions
    WHERE transaction_id = p_transaction_id
      AND deleted_at     IS NULL
    FOR UPDATE;

    IF NOT FOUND THEN
        RAISE EXCEPTION 'transaction_not_found';
    END IF;

    IF v_tx.status != 'completed' THEN
        RAISE EXCEPTION 'transaction_not_refundable';
    END IF;

    -- Restore event capacity
    UPDATE events
    SET capacity   = capacity + v_tx.quantity,
        updated_at = now()
    WHERE event_id = v_tx.event_id;

    -- Restore ticket type availability
    IF v_tx.ticket_type_id IS NOT NULL THEN
        UPDATE ticket_types
        SET available  = available + v_tx.quantity,
            updated_at = now()
        WHERE ticket_type_id = v_tx.ticket_type_id;
    END IF;

    -- Mark transaction refunded
    UPDATE transactions
    SET status     = 'refunded',
        updated_at = now()
    WHERE transaction_id = p_transaction_id;
END;
$$ LANGUAGE plpgsql;