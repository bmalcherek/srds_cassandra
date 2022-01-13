package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/bmalcherek/srds_cassandra/models"
	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx/v2"
	"github.com/scylladb/gocqlx/v2/qb"
	"github.com/scylladb/gocqlx/v2/table"
)

const (
	sectorCapacity int = 5000
	sectorRowCount int = 50
	rowSeatCount   int = 100
)

var gameMetadata = table.Metadata{
	Name:    "game",
	Columns: []string{"game_id", "game_date", "game_team1", "game_team2", "stadium_name", "capacity"},
	PartKey: []string{"game_id"},
}

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

		batch.Query(stmt, gameId, fmt.Sprintf("%02d-%02d-%02d", sector, row, seat), rand.Intn(200))
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

	// err := session.Query()
}

func main() {
	cluster := gocql.NewCluster("127.0.0.1")
	cluster.Keyspace = "test"

	// session
	// session, err := cluster.CreateSession()
	session, err := gocqlx.WrapSession(cluster.CreateSession())
	if err != nil {
		panic(err)
	}
	defer session.Close()

	createGame(&session, time.Now().Unix()*1000, 700, "Lusail Iconic Stadium", "Nigeria", "Germany")

	// ctx := context.Background()

	// gameMetadata := table.Metadata{
	// 	Name:    "games",
	// 	Columns: []string{"game_id", "game_date", "game_team1", "game_team2", "stadium_name", "capacity"},
	// 	PartKey: []string{"game_id"},
	// }

	// gameTable := table.New(gameMetadata)

	// type Game struct {
	// 	GameId      string
	// 	GameDate    int64
	// 	GameTeam1   string
	// 	GameTeam2   string
	// 	StadiumName string
	// 	Capacity    int
	// }

	// g := Game{
	// 	GameId:      "cd97ff90-7191-11ec-8d7e-5b0fd7190d80",
	// 	GameDate:    1641933312,
	// 	GameTeam1:   "Nigeria",
	// 	GameTeam2:   "Germany",
	// 	StadiumName: "Lusail",
	// 	Capacity:    80000,
	// }

	// q := session.Query(gameTable.Insert()).BindStruct(g)
	// if err := q.ExecRelease(); err != nil {
	// 	panic(err)
	// }

}
