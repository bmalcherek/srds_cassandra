// Code generated by "gocqlx/cmd/schemagen"; DO NOT EDIT.

package models

import "github.com/scylladb/gocqlx/v2/table"

// Table models.
var (
	GameReservations = table.New(table.Metadata{
		Name: "game_reservations",
		Columns: []string{
			"game_id",
			"seat_discount",
			"seat_id",
			"seat_owner",
			"seat_price",
		},
		PartKey: []string{
			"game_id",
		},
		SortKey: []string{
			"seat_id",
		},
	})

	Games = table.New(table.Metadata{
		Name: "games",
		Columns: []string{
			"capacity",
			"full_capacity",
			"game_date",
			"game_id",
			"game_team1",
			"game_team2",
			"stadium_name",
		},
		PartKey: []string{
			"game_id",
		},
		SortKey: []string{},
	})

	GamesByStadiums = table.New(table.Metadata{
		Name: "games_by_stadiums",
		Columns: []string{
			"capacity",
			"full_capacity",
			"game_date",
			"game_id",
			"game_team1",
			"game_team2",
			"stadium_name",
		},
		PartKey: []string{
			"stadium_name",
		},
		SortKey: []string{
			"game_id",
		},
	})
)
