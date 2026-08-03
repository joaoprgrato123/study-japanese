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

type ChallengeRequest struct {
	Category string `json:"category"`
	Level    string `json:"level"`
	Amount   int    `json:"amount"`
}

type ChallengeQuestion struct {
	Character string   `json:"character"`
	Options   []string `json:"options"`
	Answer    string   `json:"answer"`
}

type ChallengeResponse struct {
	Questions []ChallengeQuestion `json:"questions"`
}
