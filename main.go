package main

import (
	"fmt"
	"math/rand"

	"github.com/bmalcherek/srds_cassandra/models"
	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx/v2"
	"github.com/scylladb/gocqlx/v2/qb"
)

func reserveRandomSeat(session *gocqlx.Session, id int) {
	for {
		var availableGames []models.Game
		q := session.Query(models.Games.SelectAll())
		if err := q.SelectRelease(&availableGames); err != nil {
			panic(err)
		}

		game := availableGames[rand.Intn(len(availableGames))]

		var seats []models.GameReservation
		q = session.Query(models.GameReservations.Select()).BindMap(qb.M{"game_id": game.GameId})
		if err := q.SelectRelease(&seats); err != nil {
			panic(err)
		}

		var emptySeats []models.GameReservation
		for _, seat := range seats {
			if seat.SeatOwner == "" {
				emptySeats = append(emptySeats, seat)
			}
		}

		if len(emptySeats) == 0 {
			fmt.Println("No empty seats")
			continue
		}
		seatToReserve := emptySeats[rand.Intn(len(emptySeats))]
		seatToReserve.SeatOwner = fmt.Sprintf("%d", id)
		q = session.Query(models.GameReservations.Insert()).BindStruct(seatToReserve)
		if err := q.ExecRelease(); err != nil {
			panic(err)
		}

		checkSeat := models.GameReservation{
			GameId: seatToReserve.GameId,
			SeatId: seatToReserve.SeatId,
		}
		q = session.Query(models.GameReservations.Get()).BindStruct(checkSeat)
		if err := q.GetRelease(&checkSeat); err != nil {
			panic(err)
		}

		if checkSeat.SeatOwner != fmt.Sprintf("%d", id) {
			fmt.Println("ERROROROROROROR")
		}

		// fmt.Println(game, len(seats), len(emptySeats))
	}
}

func logger(session *gocqlx.Session) {

	var availableGames []models.Game
	q := session.Query(models.Games.SelectAll())
	if err := q.SelectRelease(&availableGames); err != nil {
		panic(err)
	}
	for {
		for _, game := range availableGames {
			var seats []models.GameReservation
			q = session.Query(models.GameReservations.Select()).BindMap(qb.M{"game_id": game.GameId})
			if err := q.SelectRelease(&seats); err != nil {
				panic(err)
			}

			var emptySeats []models.GameReservation
			for _, seat := range seats {
				if seat.SeatOwner == "" {
					emptySeats = append(emptySeats, seat)
				}
			}

			ratio := float32(len(emptySeats)) / float32(game.Capacity)
			fmt.Printf("%v, %s, %s, tickets left: %f\n", game.GameId, game.GameTeam1, game.GameTeam2, ratio)
		}
	}
}

func initTables(session *gocqlx.Session) {
	err := session.ExecStmt(`CREATE KEYSPACE IF NOT EXISTS tickets
		WITH REPLICATION = {
			'class': 'SimpleStrategy',
			'replication_factor': 2
	};`)
	if err != nil {
		panic(err)
	}

	err = session.ExecStmt(`CREATE TABLE IF NOT EXISTS tickets.games_by_stadiums (
		game_id uuid,
		game_date timestamp,
		game_team1 text,
		game_team2 text,
		stadium_name text,
		capacity int,
		PRIMARY KEY (stadium_name, game_id)
	);`)
	if err != nil {
		panic(err)
	}

	err = session.ExecStmt(`CREATE TABLE IF NOT EXISTS tickets.games (
		game_id uuid,
		game_date timestamp,
		game_team1 text,
		game_team2 text,
		stadium_name text,
		capacity int,
		PRIMARY KEY (game_id)
	);`)
	if err != nil {
		panic(err)
	}

	err = session.ExecStmt(`CREATE TABLE IF NOT EXISTS tickets.game_reservations (
		game_id uuid,
		seat_id text,
		seat_owner text,
		seat_price int,
		seat_discount text,
		PRIMARY KEY (game_id, seat_id)
	);`)
	if err != nil {
		panic(err)
	}

	err = session.ExecStmt("TRUNCATE tickets.games;")
	if err != nil {
		panic(err)
	}

	err = session.ExecStmt("TRUNCATE tickets.games_by_stadiums;")
	if err != nil {
		panic(err)
	}

	err = session.ExecStmt("TRUNCATE tickets.game_reservations;")
	if err != nil {
		panic(err)
	}
}

func main() {
	cluster := gocql.NewCluster("127.0.0.1")

	session, err := gocqlx.WrapSession(cluster.CreateSession())
	if err != nil {
		panic(err)
	}

	// initTables(&session)
	session.Close()

	cluster.Keyspace = "tickets"
	session, err = gocqlx.WrapSession(cluster.CreateSession())
	if err != nil {
		panic(err)
	}
	defer session.Close()

	endChan := make(chan (bool))

	// matches.CreateMatches(&session)

	for i := 0; i < 10; i++ {
		go reserveRandomSeat(&session, i)
	}
	go logger(&session)

	<-endChan
}
