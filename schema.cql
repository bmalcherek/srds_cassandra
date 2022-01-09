CREATE TABLE game (
    game_id uuid,
    game_date date,
    game_team1 text,
    game_team2 text,
    stadium uuid,
    stadium_capacity int,
    PRIMARY KEY (game_id)
)