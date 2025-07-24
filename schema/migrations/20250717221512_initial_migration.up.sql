CREATE TABLE wallets (
    id UUID PRIMARY KEY,
    balance INTEGER NOT NULL DEFAULT 0 CHECK (balance >= 0)
);
