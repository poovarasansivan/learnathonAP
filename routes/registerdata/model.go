package registerdata

import "database/sql"

type CTeam struct {
	User1     string         `json:"user_1"`
	User1Name string         `json:"user1_name"`
	User2     sql.NullString `json:"user_2"`
	User2Name string         `json:"user2_name"`
	User3     sql.NullString `json:"user_3"`
	User3Name string         `json:"user3_name"`
}