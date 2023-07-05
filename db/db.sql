CREATE DATABASE IF NOT EXISTS usermanagement;

USE usermanagement;

DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS role;

CREATE TABLE role (
    id           BIGINT(20) PRIMARY KEY AUTO_INCREMENT,
    name         VARCHAR(128) NOT NULL
);

INSERT INTO role (name) VALUES
    ("user"),
    ("admin");

CREATE TABLE users (
    id           BIGINT(20) PRIMARY KEY AUTO_INCREMENT,
    name         VARCHAR(128) NOT NULL,
    email        VARCHAR(128) NOT NULL UNIQUE,
    password     LONGTEXT NOT NULL,
    role_id      BIGINT(20) NOT NULL,
    created_at   DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
    last_login   DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
    FOREIGN KEY (role_id) REFERENCES role(id)
);
