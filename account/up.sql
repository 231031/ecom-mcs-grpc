CREATE TABLE IF NOT EXISTS sellers (
    id VARCHAR(27) PRIMARY KEY,
    email VARCHAR(64) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    store_name VARCHAR(64) NOT NULL,
    first_name VARCHAR(64) NOT NULL,
    last_name VARCHAR(64) NOT NULL,
    phone VARCHAR(64) NOT NULL,
    address TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS buyers (
    id VARCHAR(27) PRIMARY KEY,
    email VARCHAR(64) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    first_name VARCHAR(64) NOT NULL,
    last_name VARCHAR(64) NOT NULL,
    phone VARCHAR(64) NOT NULL,
    address TEXT NOT NULL
);