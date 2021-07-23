BEGIN;

CREATE TABLE IF NOT EXISTS users
(
    id                  int PRIMARY KEY AUTO_INCREMENT,
    account_id          int          NOT NULL,
    nickname            varchar(255) NOT NULL,
    profile_picture_uri varchar(255) NOT NULL,
    created_at          TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at          DATETIME  DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

ALTER TABLE users
    ADD FOREIGN KEY (account_id) REFERENCES accounts (id);

CREATE INDEX users_index_0 ON users (account_id);

COMMIT;