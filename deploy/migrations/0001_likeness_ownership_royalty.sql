CREATE TABLE IF NOT EXISTS real_person_models (
    id UUID PRIMARY KEY,
    person_id TEXT NOT NULL,
    display_name TEXT NOT NULL,
    verification_id TEXT NOT NULL,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    revoked_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS consent_agreements (
    id UUID PRIMARY KEY,
    model_id UUID NOT NULL REFERENCES real_person_models(id) ON DELETE CASCADE,
    allowed_uses JSONB NOT NULL DEFAULT '[]'::jsonb,
    restricted_uses JSONB NOT NULL DEFAULT '[]'::jsonb,
    territories JSONB NOT NULL DEFAULT '[]'::jsonb,
    starts_at TIMESTAMPTZ NOT NULL,
    ends_at TIMESTAMPTZ NOT NULL,
    revoked_at TIMESTAMPTZ,
    verification_ref TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS character_ownerships (
    character_id UUID PRIMARY KEY,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS ownership_splits (
    character_id UUID NOT NULL REFERENCES character_ownerships(character_id) ON DELETE CASCADE,
    owner_id UUID NOT NULL,
    percentage NUMERIC(5,2) NOT NULL CHECK (percentage > 0),
    PRIMARY KEY(character_id, owner_id)
);

CREATE TABLE IF NOT EXISTS royalty_transactions (
    id UUID PRIMARY KEY,
    character_id UUID NOT NULL,
    gross_revenue NUMERIC(14,2) NOT NULL CHECK (gross_revenue >= 0),
    currency CHAR(3) NOT NULL,
    payout_by_owner JSONB NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS consent_usage_audit (
    id BIGSERIAL PRIMARY KEY,
    consent_id UUID NOT NULL,
    model_id UUID NOT NULL,
    requester_id TEXT NOT NULL,
    intended_use TEXT NOT NULL,
    allowed BOOLEAN NOT NULL,
    reason TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS immutable_ledger (
    seq BIGSERIAL PRIMARY KEY,
    event_type TEXT NOT NULL,
    aggregate_id TEXT NOT NULL,
    payload JSONB NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE OR REPLACE FUNCTION prevent_immutable_ledger_mutation()
RETURNS TRIGGER AS $$
BEGIN
    RAISE EXCEPTION 'immutable_ledger is append-only';
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS immutable_ledger_no_update ON immutable_ledger;
CREATE TRIGGER immutable_ledger_no_update
BEFORE UPDATE OR DELETE ON immutable_ledger
FOR EACH ROW
EXECUTE FUNCTION prevent_immutable_ledger_mutation();
