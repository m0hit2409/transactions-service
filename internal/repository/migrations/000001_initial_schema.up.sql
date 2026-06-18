CREATE TABLE IF NOT EXISTS accounts (
    account_id      INTEGER  PRIMARY KEY AUTOINCREMENT,
    document_number TEXT     NOT NULL UNIQUE,
    is_active       INTEGER  NOT NULL DEFAULT 1,
    created_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS operation_types (
    operation_type_id INTEGER PRIMARY KEY,
    description       TEXT    NOT NULL,
    sign              INTEGER NOT NULL CHECK (sign IN (-1, 1)),
    is_active         INTEGER NOT NULL DEFAULT 1,
    created_at        DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

INSERT OR IGNORE INTO operation_types (operation_type_id, description, sign) VALUES
    (1, 'Normal Purchase',            -1),
    (2, 'Purchase with installments', -1),
    (3, 'Withdrawal',                 -1),
    (4, 'Credit Voucher',              1);

CREATE TABLE IF NOT EXISTS transactions (
    transaction_id    INTEGER  PRIMARY KEY AUTOINCREMENT,
    account_id        INTEGER  NOT NULL REFERENCES accounts (account_id),
    operation_type_id INTEGER  NOT NULL REFERENCES operation_types (operation_type_id),
    amount            TEXT     NOT NULL,
    event_date        DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);
