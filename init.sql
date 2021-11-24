CREATE TABLE IF NOT EXISTS wallet (
    id VARCHAR(64) PRIMARY KEY,
    balance INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TYPE operation AS ENUM ('deposit', 'withdraw');

CREATE TABLE IF NOT EXISTS operations (
    id BIGSERIAL PRIMARY KEY,
    wallet_id VARCHAR(64) NOT NULL REFERENCES wallet (id),
    operation_type operation NOT NULL,
    amount INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX wallet_operation_idx ON operations USING BTREE (wallet_id, operation_type);
CREATE INDEX created_at_idx ON operations USING BTREE (created_at);
