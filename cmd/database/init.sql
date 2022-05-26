CREATE TABLE IF NOT EXISTS players (
    id uuid,
    name varchar(250) NOT NULL,
    password varchar(250) NOT NULL,
    state ENUM('offline', 'online-availale', 'online-playing') NOT NULL,
    connected_ip STRING NULL,
    connected_port STRING NULL,
    points INT NOT NULL,
    PRIMARY KEY (id)
);
