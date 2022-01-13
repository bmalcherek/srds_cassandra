CREATE TABLE games_by_stadiums (
    game_id uuid,
    game_date timestamp,
    game_team1 text,
    game_team2 text,
    stadium_name text,
    capacity int,
    PRIMARY KEY (stadium_name, game_id)
);

CREATE TABLE games (
    game_id uuid,
    game_date timestamp,
    game_team1 text,
    game_team2 text,
    stadium_name text,
    capacity int,
    PRIMARY KEY (game_id)
);

INSERT INTO games (stadium_name, game_id, game_date, game_team1, game_team2, capacity) VALUES ('Lusail Iconic Stadium', cd97ff90-7191-11ec-8d7e-5b0fd7190d80, 1641933312, 'Nigeria', 'Germany', 80000);
INSERT INTO games (stadium_name, game_id, game_date, game_team1, game_team2, capacity) VALUES ('Lusail Iconic Stadium', 37abfec0-7194-11ec-8d7e-5b0fd7190d80, 1641933312, 'Poland', 'Ukraine', 80000);


CREATE TABLE stadiums (
    stadium_name text,
    max_capacity int,
    city text,
    PRIMARY KEY (stadium_name)
);

INSERT INTO stadiums (stadium_name, max_capacity, city) VALUES ('Lusail Iconic Stadium', 80000, 'Lusail');

CREATE TABLE game_reservations (
    game_id uuid,
    seat_id text,
    seat_owner text,
    seat_price int,
    seat_discount text,
    PRIMARY KEY (game_id, seat_id)
);

INSERT INTO game_reservations (game_id, seat_id, seat_price) VALUES (cd97ff90-7191-11ec-8d7e-5b0fd7190d80, 'A01', 69);
INSERT INTO game_reservations (game_id, seat_id, seat_price) VALUES (cd97ff90-7191-11ec-8d7e-5b0fd7190d80, 'A02', 420);
INSERT INTO game_reservations (game_id, seat_id, seat_price) VALUES (cd97ff90-7191-11ec-8d7e-5b0fd7190d80, 'A03', 2137);


INSERT INTO game_reservations (game_id, seat_id, seat_price) VALUES (37abfec0-7194-11ec-8d7e-5b0fd7190d80, 'B01', 69);

SELECT * FROM game_reservations WHERE game_id = cd97ff90-7191-11ec-8d7e-5b0fd7190d80 and seat_id = 'A01';