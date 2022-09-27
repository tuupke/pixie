CREATE TABLE `external_data`
(
    `guid`         text NOT NULL,
    `user_id`      integer,
    `team_id`      integer,
    `username`     text,
    `teamname`     text,

    `loc_x`        real NULL,
    `loc_y`        real NULL,
    `loc_rotation` real NULL,
    `host_id`      text NULL,

    PRIMARY KEY ('guid')
);

CREATE UNIQUE INDEX idx_dj_unique_team ON external_data (`team_id`);
CREATE INDEX idx_dj_host ON external_data (`host_id`);

CREATE table `hosts`
(
    `guid`        text NOT NULL,
    `data`        text NOT NULL,
    `hostname`    text NOT NULL,
    `primary_ip`  text NOT NULL,
    `primary_mac` text NOT NULL,
    `last_seen`   datetime,

    PRIMARY KEY (`guid`)
);

CREATE UNIQUE INDEX idx_guid_unique ON hosts (`guid`);
-- CREATE INDEX idx_hosts ON external_data (`host_id`);

CREATE table `settings`
(
    `key`   text not null,
    `value` text,

    PRIMARY KEY (`key`)
);

CREATE table `problems`
(
    `id`    text NOT NULL,
    `rgb`   text NULL,
    `loc_x` real null,
    `loc_y` real null,

    PRIMARY KEY (`id`)
);

insert into settings
values
    ('contest', 'nwerc18'),
    ('domjudge', 'https://www.domjudge.org/demoweb/api/v4/')
       ;
