package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/bmalcherek/srds_cassandra/matches"
	"github.com/bmalcherek/srds_cassandra/models"
	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx/v2"
	"github.com/scylladb/gocqlx/v2/qb"
)

const (
	goroutinesCount int = 10
)

func reserveRandomSeat(session *gocqlx.Session, id int, endChan chan (bool)) {
	for {
		var availableGames []models.Game
		q := session.Query(models.Games.SelectAll())
		if err := q.SelectRelease(&availableGames); err != nil {
			panic(err)
		}

		notFullGames := []models.Game{}
		for _, g := range availableGames {
			if !g.FullCapacity {
				notFullGames = append(notFullGames, g)
			}
		}

		if len(notFullGames) == 0 {
			endChan <- true
			return
		}

		game := notFullGames[rand.Intn(len(notFullGames))]

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
			game.FullCapacity = true
			q = session.Query(models.Games.Insert()).BindStruct(game)
			if err := q.ExecRelease(); err != nil {
				panic(err)
			}
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
			fmt.Printf("ERROR!! Expected value: %d, real value: %s\n", id, checkSeat.SeatOwner)
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

			time.Sleep(500 * time.Millisecond)
		}
	}
}

func initTables(session *gocqlx.Session) {
	err := session.ExecStmt("DROP TABLE tickets.games;")
	if err != nil {
		panic(err)
	}

	err = session.ExecStmt("DROP TABLE tickets.games_by_stadiums;")
	if err != nil {
		panic(err)
	}

	err = session.ExecStmt("DROP TABLE tickets.game_reservations;")
	if err != nil {
		panic(err)
	}

	err = session.ExecStmt(`CREATE KEYSPACE IF NOT EXISTS tickets
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
		full_capacity boolean,
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
		full_capacity boolean,
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
}

func main() {
	cluster := gocql.NewCluster("127.0.0.1")

	session, err := gocqlx.WrapSession(cluster.CreateSession())
	if err != nil {
		panic(err)
	}

	initTables(&session)
	session.Close()

	cluster.Keyspace = "tickets"
	session, err = gocqlx.WrapSession(cluster.CreateSession())
	if err != nil {
		panic(err)
	}
	defer session.Close()

	endChan := make(chan (bool))

	matches.CreateMatches(&session)

	for i := 0; i < goroutinesCount; i++ {
		go reserveRandomSeat(&session, i, endChan)
	}
	go logger(&session)

	for i := 0; i < goroutinesCount; i++ {
		<-endChan
	}
}
