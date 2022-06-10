CREATE TABLE IF NOT EXISTS printers
(
    printer_id INTEGER PRIMARY KEY,
    name       TEXT NOT NULL,
    uri        TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS 

CREATE TABLE IF NOT EXISTS rooms
(
    room_id         INTEGER PRIMARY KEY,
    name            text,
    default_printer INTEGER NOT NULL REFERENCES printers (printer_id),
    shape           text
);

CREATE TABLE IF NOT EXISTS teams
(
    team_id          INTEGER PRIMARY KEY,
    room_id          INTEGER REFERENCES rooms (room_id),
    override_printer INTEGER NULL REFERENCES printers (printer_id),
    domjudge_id      integer NOT NULL,
    user_id          integer NULL,

    name             text,
    x                INTEGER,
    y                INTEGER,
    orientation      INTEGER
);

CREATE TABLE IF NOT EXISTS hosts
(
    host_id    INTEGER PRIMARY KEY,
    team_id    INTEGER NULL REFERENCES teams (team_id),
    identifier text,
    hostname   text
);

CREATE UNIQUE INDEX IF NOT EXISTS hosts_team_id ON hosts (identifier);

CREATE TABLE IF NOT EXISTS ips
(
    host_id INTEGER REFERENCES hosts (host_id),
    ip      text,
    PRIMARY KEY (host_id, ip ASC)
);