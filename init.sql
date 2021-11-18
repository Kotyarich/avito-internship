CREATE TABLE IF NOT EXISTS balances(
  id SERIAL PRIMARY KEY,
  amount NUMERIC(1000, 2) NOT NULL DEFAULT 0
);

CREATE TYPE transaction_type AS ENUM ('product', 'transfer', 'fill');

CREATE TABLE IF NOT EXISTS transactions(
  id SERIAL PRIMARY KEY,
  user_id INTEGER NOT NULL REFERENCES balances(id),
  amount NUMERIC(1000, 2) NOT NULL,
  target_id INTEGER NOT NULL,
  type transaction_type NOT NULL,
  date TIMESTAMP DEFAULT NOW()
);