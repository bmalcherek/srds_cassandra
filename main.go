package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/bmalcherek/srds_cassandra/models"
	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx/v2"
	"github.com/scylladb/gocqlx/v2/qb"
)

const (
	sectorCapacity int = 5000
	sectorRowCount int = 50
	rowSeatCount   int = 100
)

func createGame(session *gocqlx.Session, date int64, capacity int, stadiumName, team1, team2 string) {
	gameId, err := gocql.RandomUUID()
	if err != nil {
		panic(err)
	}
	g := models.Game{
		GameId:      gameId,
		GameDate:    date,
		GameTeam1:   team1,
		GameTeam2:   team2,
		StadiumName: stadiumName,
		Capacity:    capacity,
	}

	q := session.Query(models.Games.Insert()).BindStruct(g)
	if err := q.ExecRelease(); err != nil {
		panic(err)
	}

	// stmt, _ := qb.Select("gocqlx_test.bench_person").Columns(benchPersonCols...).Where(qb.Eq("id")).Limit(1).ToCql()
	selectG := models.Game{
		GameId: gameId,
	}
	qu := session.Query(models.Games.Get()).BindStruct(selectG)
	if err := qu.GetRelease(&selectG); err != nil {
		panic(err)
	}
	fmt.Println(selectG)

	batch := session.NewBatch(gocql.LoggedBatch)
	stmt := `INSERT INTO game_reservations (game_id, seat_id, seat_price) VALUES (?, ?, ?)`

	sector := 0
	row := 0
	seat := 0
	for i := 0; i < capacity; i++ {
		if i%sectorCapacity == 0 {
			sector += 1
			row = 0
			seat = 0
		}
		if i%rowSeatCount == 0 && i%sectorCapacity != 0 {
			row += 1
			seat = 0
		}

		batch.Query(stmt, gameId, fmt.Sprintf("%04d-%02d-%02d", sector, row, seat), rand.Intn(200))
		if i%100 == 0 {
			err = session.ExecuteBatch(batch)
			if err != nil {
				panic(err)
			}
			batch = session.NewBatch(gocql.LoggedBatch)
		}

		seat += 1
	}
	err = session.ExecuteBatch(batch)
	if err != nil {
		panic(err)
	}

	var gameReservations []models.GameReservation

	q1 := session.Query(models.GameReservations.Select()).BindMap(qb.M{"game_id": gameId})
	if err = q1.SelectRelease(&gameReservations); err != nil {
		panic(err)
	}

	fmt.Println(len(gameReservations))
}

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

	initTables(&session)
	session.Close()

	cluster.Keyspace = "tickets"
	session, err = gocqlx.WrapSession(cluster.CreateSession())
	if err != nil {
		panic(err)
	}
	defer session.Close()

	endChan := make(chan (bool))

	createGame(&session, time.Now().Unix()*1000, 100, "Lusail Iconic Stadium", "Nigeria", "Germany")

	for i := 0; i < 5; i++ {
		go reserveRandomSeat(&session, i)
	}

	<-endChan
}
