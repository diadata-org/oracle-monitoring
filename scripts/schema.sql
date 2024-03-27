/*CREATE DATABASE monitoring;*/

CREATE TABLE IF NOT EXISTS oracles (
  contract_address TEXT NOT NULL PRIMARY KEY,
  contract_abi TEXT NOT NULL,
  chain_id TEXT NOT NULL,
  node_url TEXT NOT NULL,
  creation_block BIGINT NOT NULL
);

-- CREATE TABLE IF NOT EXISTS feederupdates (
--   id BIGSERIAL PRIMARY KEY,
--   oracle_address TEXT NOT NULL ,
--   chain_id TEXT NOT NULL ,
--   transaction_hash TEXT NOT NULL UNIQUE,
--   transaction_cost TEXT NOT NULL,
--   asset_key TEXT NOT NULL,
--   asset_price TEXT NOT NULL,
--   update_block BIGINT NOT NULL,
--   update_from TEXT NOT NULL,
--   from_balance TEXT NOT NULL,
--   gas_cost TEXT NOT NULL,
--   gas_used TEXT NOT NULL,
--   update_time TIMESTAMP  NULL
-- );
CREATE TABLE feederupdates (
    id bigint NOT NULL,
    oracle_address text NOT NULL,
    transaction_hash text NOT NULL,
    transaction_cost text NOT NULL,
    asset_key text NOT NULL,
    asset_price text NOT NULL,
    update_block bigint NOT NULL,
    update_from text NOT NULL,
    from_balance text NOT NULL,
    gas_cost text NOT NULL,
    gas_used text NOT NULL,
    creation_block bigint NOT NULL,
    chain_id text,
    update_time timestamp without time zone
);

CREATE TABLE IF NOT EXISTS feederupdatestate (
  id BIGSERIAL PRIMARY KEY,
  chain_id TEXT NOT NULL ,
  last_block BIGINT NOT NULL
);

ALTER TABLE feederupdates
ALTER COLUMN gas_used TYPE double precision USING gas_used::double precision;