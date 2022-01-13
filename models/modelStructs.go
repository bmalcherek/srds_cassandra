package models

import "github.com/gocql/gocql"

type Game struct {
	GameId      gocql.UUID
	GameDate    int64
	GameTeam1   string
	GameTeam2   string
	StadiumName string
	Capacity    int
}

type Stadium struct {
	StadiumName string
	MaxCapacity int
	City        string
}

type GameReservation struct {
	GameId       gocql.UUID
	SeatId       string
	SeatOwner    string
	SeatPrice    int
	SeatDiscount string
}
