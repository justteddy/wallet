CREATE DATABASE wallets;
\connect wallets;

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

CREATE INDEX wallet_idx ON operations USING HASH (wallet_id);
CREATE INDEX operation_type_idx ON operations USING HASH (operation_type);
CREATE INDEX created_at_idx ON operations USING BTREE (created_at);
