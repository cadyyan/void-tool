-- name: GetAllPlayers :many
SELECT
    id,
    username,
    created_on
FROM players;

-- name: GetPlayerByName :one
SELECT
    id,
    username,
    created_on
FROM players
WHERE
    username = ?;

-- name: CreatePlayer :one
INSERT INTO players (
    id,
    username,
    created_on
) VALUES (
    ?,
    ?,
    ?
) RETURNING *;

-- name: CreatePlayerIfNotExist :exec
INSERT OR IGNORE INTO players (
    id,
    username,
    created_on
) VALUES (
    ?,
    ?,
    ?
);

-- name: GetAllPlayerSkillsByPlayerName :many
WITH ordered_skill_entries AS (
    SELECT
        player_skills.player_id,
        player_skills.name,
        player_skills.day,
        player_skills.experience,
        player_skills.level,
        ROW_NUMBER() OVER (
            PARTITION BY player_skills.name
            ORDER BY player_skills.day DESC
        ) AS row_num
    FROM player_skills
    INNER JOIN players
        ON
            player_skills.player_id = players.id
            AND
            players.username = ?
)

SELECT
    ordered_skill_entries.player_id,
    ordered_skill_entries.name,
    ordered_skill_entries.day,
    ordered_skill_entries.experience,
    ordered_skill_entries.level
FROM ordered_skill_entries
WHERE
    ordered_skill_entries.row_num = 1;

-- name: GetPlayerSkillOverTimeByPlayerName :many
SELECT
    player_skills.player_id,
    player_skills.name,
    player_skills.day,
    player_skills.experience,
    player_skills.level
FROM player_skills
INNER JOIN players
    ON
        player_skills.player_id = players.id
        AND
        players.username = ?
WHERE
    player_skills.name = ?
ORDER BY player_skills.day ASC;

-- name: GetHighscoresForSkill :many
WITH latest_row AS (
    SELECT
        player_id,
        name,
        MAX(day) AS latest_day
    FROM player_skills
    GROUP BY player_id, name
),

latest_skills_by_player AS (
    SELECT
        player_skills.player_id,
        player_skills.name,
        player_skills.experience,
        player_skills.level
    FROM player_skills
    INNER JOIN latest_row
        ON
            player_skills.player_id = latest_row.player_id
            AND
            player_skills.name = latest_row.name
            AND
            player_skills.day = latest_row.latest_day
    WHERE
        player_skills.name = ?
)

SELECT
    players.id,
    players.username,
    skills.name,
    skills.experience,
    skills.level
FROM latest_skills_by_player AS skills
INNER JOIN players
    ON
        skills.player_id = players.id
ORDER BY skills.experience DESC;

-- name: RecordPlayerSkill :exec
INSERT INTO player_skills (
    player_id,
    name,
    day,
    experience,
    level
) VALUES (
    ?,
    ?,
    ?,
    ?,
    ?
) ON CONFLICT (player_id, name, day)
DO UPDATE SET experience = excluded.experience, level = excluded.level;
