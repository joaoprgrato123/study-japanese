package character

import "time"

type Character struct {
	ID        int64
	Character string
	Reading   string
	Meaning   string
	Category  string
	Level     string
	CreatedAt time.Time
}

type ChallengeParams struct {
	Category string
	Level    string
	Amount   int
}

type Question struct {
	Character string
	Options   []string
	Answer    string
}
