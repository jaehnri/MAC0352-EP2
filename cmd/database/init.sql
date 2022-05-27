CREATE TABLE IF NOT EXISTS players (
    id UUID,
    name varchar(250) NOT NULL,
    password varchar(250) NOT NULL,
    state varchar(20) NOT NULL,
    connected_ip varchar(15) NULL,
    connected_port INT NULL,
    points INT NOT NULL,
    PRIMARY KEY (id),
    UNIQUE (name)
);