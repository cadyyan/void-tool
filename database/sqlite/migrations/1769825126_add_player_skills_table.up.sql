PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS player_skills (
    player_id VARCHAR NOT NULL,
    name VARCHAR NOT NULL,
    day VARCHAR NOT NULL CHECK (day IS date(day)),
    level INT NOT NULL CHECK (level >= 1),
    experience REAL NOT NULL CHECK (experience >= 0),

    PRIMARY KEY (player_id, name, day),
    FOREIGN KEY (player_id) REFERENCES players (id)
);

CREATE INDEX IF NOT EXISTS idx__player_skills__lookup ON player_skills (
    player_id,
    name
);
