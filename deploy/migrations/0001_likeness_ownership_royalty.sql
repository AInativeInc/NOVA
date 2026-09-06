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

CREATE OR REPLACE FUNCTION validate_ownership_split_total()
RETURNS TRIGGER AS $$
DECLARE
    target_character_id UUID;
    total NUMERIC(7,2);
BEGIN
    target_character_id := COALESCE(NEW.character_id, OLD.character_id);
    SELECT COALESCE(SUM(percentage), 0) INTO total
    FROM ownership_splits
    WHERE character_id = target_character_id;

    IF total <> 0 AND total <> 100.00::NUMERIC(7,2) THEN
        RAISE EXCEPTION 'ownership split total for character % must equal 100%% when splits exist', target_character_id;
    END IF;

    RETURN COALESCE(NEW, OLD);
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS ownership_splits_validate_total ON ownership_splits;
CREATE CONSTRAINT TRIGGER ownership_splits_validate_total
AFTER INSERT OR UPDATE OR DELETE ON ownership_splits
DEFERRABLE INITIALLY DEFERRED
FOR EACH ROW
EXECUTE FUNCTION validate_ownership_split_total();

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
