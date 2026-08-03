package character

import "time"

type Character struct {
	ID        int64     `json:"id"`
	Character string    `json:"character"`
	Reading   string    `json:"reading"`
	Meaning   string    `json:"meaning"`
	Category  string    `json:"category"`
	Level     string    `json:"level"`
	CreatedAt time.Time `json:"created_at"`
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
