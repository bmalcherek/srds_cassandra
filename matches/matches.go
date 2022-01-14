package matches

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/bmalcherek/srds_cassandra/models"
	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx/v2"
)

const (
	sectorCapacity int = 5000
	sectorRowCount int = 50
	rowSeatCount   int = 100
)

func createGame(session *gocqlx.Session, capacity int, stadiumName, team1, team2 string) {
	gameId, err := gocql.RandomUUID()
	if err != nil {
		panic(err)
	}
	g := models.Game{
		GameId:      gameId,
		GameDate:    time.Now().Unix() * 1000,
		GameTeam1:   team1,
		GameTeam2:   team2,
		StadiumName: stadiumName,
		Capacity:    capacity,
	}

	q := session.Query(models.Games.Insert()).BindStruct(g)
	if err := q.ExecRelease(); err != nil {
		panic(err)
	}

	q = session.Query(models.GamesByStadiums.Insert()).BindStruct(g)
	if err := q.ExecRelease(); err != nil {
		panic(err)
	}

	selectG := models.Game{
		GameId: gameId,
	}
	qu := session.Query(models.Games.Get()).BindStruct(selectG)
	if err := qu.GetRelease(&selectG); err != nil {
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

	// var gameReservations []models.GameReservation

	// q1 := session.Query(models.GameReservations.Select()).BindMap(qb.M{"game_id": gameId})
	// if err = q1.SelectRelease(&gameReservations); err != nil {
	// 	panic(err)
	// }

	// fmt.Println(len(gameReservations))
}

func groupA(session *gocqlx.Session) {
	createGame(session, 780, "Luzhniki Stadium", "Russia", "Saudi Arabia")
	createGame(session, 270, "Central Stadium", "Egypt", "Uruguay")
	createGame(session, 644, "Krestovsky Stadium", "Russia", "Egypt")
	createGame(session, 422, "Rostov Arena", "Uruguay", "Saudi Arabia")
	createGame(session, 419, "Cosmos Arena", "Uruguay", "Russia")
	createGame(session, 368, "Volgograd Arena", "Saudi Arabia", "Egypt")
	fmt.Println("Group A created...")
}
func groupB(session *gocqlx.Session) {
	createGame(session, 625, "Krestovsky Stadium", "Morocco", "Iran")
	createGame(session, 438, "Fisht Olympic Stadium", "Portugal", "Spain")
	createGame(session, 780, "Luzhniki Stadium", "Portugal", "Morocco")
	createGame(session, 427, "Kazan Arena", "Iran", "Spain")
	createGame(session, 416, "Mordovia Arena", "Iran", "Portugal")
	createGame(session, 339, "Kaliningrad Stadium", "Spain", "Morocco")
	fmt.Println("Group B created...")
}
func groupC(session *gocqlx.Session) {
	createGame(session, 412, "Kazan Arena", "France", "Australia")
	createGame(session, 405, "Mordovia Arena", "Peru", "Denmark")
	createGame(session, 407, "Cosmos Arena", "Denmark", "Australia")
	createGame(session, 327, "Central Stadium", "France", "Peru")
	createGame(session, 780, "Luzhniki Stadium", "Denmark", "France")
	createGame(session, 440, "Fisht Olympic Stadium", "Australia", "Peru")
	fmt.Println("Group C created...")
}

func groupD(session *gocqlx.Session) {
	createGame(session, 441, "Otkritie Arena", "Argentina", "Iceland")
	createGame(session, 311, "Kaliningrad Stadium", "Croatia", "Nigeria")
	createGame(session, 433, "Nizhny Novgorod Stadium", "Argentina", "Croatia")
	createGame(session, 409, "Volgograd Arena", "Nigeria", "Iceland")
	createGame(session, 644, "Krestovsky Stadium", "Nigeria", "Argentina")
	createGame(session, 433, "Rostov Arena", "Iceland", "Croatia")
	fmt.Println("Group D created...")
}

func groupE(session *gocqlx.Session) {
	createGame(session, 414, "Cosmos Arena", "Costa Rica", "Serbia")
	createGame(session, 431, "Rostov Arena", "Brazil", "Switzerland")
	createGame(session, 644, "Krestovsky Stadium", "Brazil", "Costa Rica")
	createGame(session, 331, "Kaliningrad Stadium", "Serbia", "Switzerland")
	createGame(session, 441, "Otkritie Arena", "Serbia", "Brazil")
	createGame(session, 433, "Nizhny Novgorod Stadium", "Switzerland", "Costa Rica")
	fmt.Println("Group E created...")
}
func groupF(session *gocqlx.Session) {
	createGame(session, 780, "Luzhniki Stadium", "Germany", "Mexico")
	createGame(session, 423, "Nizhny Novgorod Stadium", "Sweden", "South Korea")
	createGame(session, 434, "Rostov Arena", "South Korea", "Mexico")
	createGame(session, 442, "Fisht Olympic Stadium", "Germany", "Sweden")
	createGame(session, 418, "Kazan Arena", "South Korea", "Germany")
	createGame(session, 330, "Central Stadium", "Mexico", "Sweden")
	fmt.Println("Group F created...")
}
func groupG(session *gocqlx.Session) {
	createGame(session, 432, "Fisht Olympic Stadium", "Belgium", "Panama")
	createGame(session, 410, "Volgograd Arena", "Tunisia", "England")
	createGame(session, 441, "Otkritie Arena", "Belgium", "Tunisia")
	createGame(session, 433, "Nizhny Novgorod Stadium", "England", "Panama")
	createGame(session, 339, "Kaliningrad Stadium", "England", "Belgium")
	createGame(session, 371, "Mordovia Arena", "Panama", "Tunisia")
	fmt.Println("Group G created...")
}

func groupH(session *gocqlx.Session) {
	createGame(session, 404, "Mordovia Arena", "Colombia", "Japan")
	createGame(session, 441, "Otkritie Arena", "Poland", "Senegal")
	createGame(session, 325, "Central Stadium", "Japan", "Senegal")
	createGame(session, 428, "Kazan Arena", "Poland", "Colombia")
	createGame(session, 421, "Volgograd Arena", "Japan", "Poland")
	createGame(session, 419, "Cosmos Arena", "Senegal", "Colombia")
	fmt.Println("Group H created...")
}

func CreateMatches(session *gocqlx.Session) {
	groupA(session)
	groupB(session)
	groupC(session)
	groupD(session)
	groupE(session)
	groupF(session)
	groupG(session)
	groupH(session)
}
