CREATE TABLE IF NOT EXISTS crons (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    name        text NOT NULL,
    expression  text NOT NULL,
    script      text NOT NULL
);

CREATE TABLE IF NOT EXISTS financial_logs (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    amount      INTEGER NOT NULL,
    currency    text NOT NULL,
    investiment text NOT NULL,
    created_at   DATE DEFAULT CURRENT_DATE
);
