CREATE TABLE IF NOT EXISTS users
(
    id serial,
    email TEXT NOT NULL UNIQUE,
    pass_hash bytea NOT NULL,
    PRIMARY KEY (id, email)  -- Включаем все столбцы разделения в PRIMARY KEY
) PARTITION BY HASH (email);

CREATE INDEX IF NOT EXISTS idx_email ON users (email);

CREATE TABLE users_p1 PARTITION OF users
FOR VALUES WITH (MODULUS 4, REMAINDER 0);

CREATE TABLE users_p2 PARTITION OF users
    FOR VALUES WITH (MODULUS 4, REMAINDER 1);

CREATE TABLE users_p3 PARTITION OF users
    FOR VALUES WITH (MODULUS 4, REMAINDER 2);

CREATE TABLE users_p4 PARTITION OF users
    FOR VALUES WITH (MODULUS 4, REMAINDER 3);
