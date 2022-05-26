CREATE TABLE IF NOT EXISTS players (
    id uuid,
    name varchar(250) NOT NULL,
    password varchar(250) NOT NULL,
    PRIMARY KEY (id)
);
