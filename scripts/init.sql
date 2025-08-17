CREATE ROLE myuser WITH LOGIN PASSWORD 'mypassword';
CREATE DATABASE my_db OWNER myuser;

\c my_db;

CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username TEXT UNIQUE,
    registeration_date TIMESTAMP NOT NULL , 
    password TEXT NOT NULL
);
CREATE TABLE IF NOT EXISTS otps (
    id SERIAL PRIMARY KEY , 
    username TEXT NOT NULL ,  
    phone_number TEXT NOT NULL ,
    otp     VARCHAR(12) NOT NULL ,  
    created_at TIMESTAMP
);

GRANT ALL PRIVILEGES ON TABLE users TO myuser;
GRANT USAGE, SELECT ON SEQUENCE users_id_seq TO myuser;

GRANT ALL PRIVILEGES ON TABLE otps TO myuser;
GRANT USAGE, SELECT ON SEQUENCE otps_id_seq TO myuser;