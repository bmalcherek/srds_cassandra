package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/gocql/gocql"
	"github.com/google/uuid"
)

const (
	sectorCapacity int = 5000
	sectorRowCount int = 50
	rowSeatCount   int = 100
)

func createGame(session *gocql.Session, ctx *context.Context, date int64, capacity int, stadiumName, team1, team2 string) {
	gameId := uuid.New().String()
	err := session.Query(`INSERT INTO games (stadium_name, game_id, game_date, game_team1, game_team2, capacity) VALUES (?, ?, ?, ?, ?, ?)`, stadiumName, gameId, date, team1, team2, capacity).WithContext(*ctx).Exec()
	if err != nil {
		panic(err)
	}
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

		batch.Query(stmt, gameId, fmt.Sprintf("%d-%d-%d", sector, row, seat), rand.Intn(200))
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

	// err := session.Query()
}

func main() {
	cluster := gocql.NewCluster("127.0.0.1")
	cluster.Keyspace = "test"

	session, err := cluster.CreateSession()
	if err != nil {
		panic(err)
	}
	defer session.Close()

	ctx := context.Background()

	// var gameId gocql.UUID
	// var gameTeam1 string

	// err = session.Query(`INSERT INTO games (game_id, seat_id) VALUES (?, ?)`, "8e654920-718e-11ec-8d7e-5b0fd7190d80", "A23").WithContext(ctx).Exec()
	// if err != nil {
	// 	panic(err)
	// }

	// scanner := session.Query(`SELECT game_id, game_team1 FROM games`).WithContext(ctx).Iter().Scanner()
	// for scanner.Next() {
	// 	err = scanner.Scan(&gameId, &gameTeam1)
	// 	if err != nil {
	// 		panic(err)
	// 	}

	// 	fmt.Println(gameId, gameTeam1)
	// }

	createGame(session, &ctx, time.Now().Unix()*1000, 80000, "Lusail Iconic Stadium", "Nigeria", "Germany")

}
