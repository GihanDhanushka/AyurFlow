-- The hash avoids persisting attempted login identifiers in the throttle ledger.
CREATE TABLE auth_login_attempts (
    identifier_hash bytea PRIMARY KEY,
    failures integer NOT NULL CHECK (failures >= 0),
    window_started_at timestamptz NOT NULL,
    blocked_until timestamptz
);
