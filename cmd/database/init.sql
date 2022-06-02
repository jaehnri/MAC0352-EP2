CREATE TABLE IF NOT EXISTS players (
    id UUID,
    name varchar(250) NOT NULL,
    password varchar(250) NOT NULL,
    state varchar(20) NOT NULL,
    address varchar(30) NULL,
    points INT NOT NULL,
    PRIMARY KEY (id),
    UNIQUE (name)
);
