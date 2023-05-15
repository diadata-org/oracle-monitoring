/*CREATE DATABASE monitoring;*/

CREATE TABLE IF NOT EXISTS oracles (
  contract_address TEXT NOT NULL PRIMARY KEY,
  contract_abi TEXT NOT NULL,
  chain_id TEXT NOT NULL,
  node_url TEXT NOT NULL,
  creation_block BIGINT NOT NULL
);

CREATE TABLE IF NOT EXISTS metrics (
  id BIGSERIAL PRIMARY KEY,
  oracle_address TEXT NOT NULL REFERENCES oracles (contract_address),
  transaction_hash TEXT NOT NULL UNIQUE,
  transaction_cost TEXT NOT NULL,
  asset_key TEXT NOT NULL,
  asset_price TEXT NOT NULL,
  update_block BIGINT NOT NULL,
  update_from TEXT NOT NULL,
  from_balance TEXT NOT NULL,
  gas_cost TEXT NOT NULL,
  gas_used TEXT NOT NULL,
  FOREIGN KEY (oracle_address) REFERENCES oracles (contract_address)
);
