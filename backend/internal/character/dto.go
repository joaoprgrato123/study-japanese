package character

// --- requests --- //
type AddCharacterRequest struct {
	Character string `json:"character"`
	Reading   string `json:"reading"`
	Meaning   string `json:"meaning"`
	Category  string `json:"category"`
	Level     string `json:"level"`
}

type ChallengeRequest struct {
	Category string `json:"category"`
	Level    string `json:"level"`
	Amount   int    `json:"amount"`
}

// --- responses --- //
type CharacterResponse struct {
	ID        int64  `json:"id"`
	Character string `json:"character"`
	Reading   string `json:"reading"`
	Meaning   string `json:"meaning"`
	Category  string `json:"category"`
	Level     string `json:"level"`
}

type QuestionResponse struct {
	Character string   `json:"character"`
	Options   []string `json:"options"`
	Answer    string   `json:"answer"`
}

type ChallengeResponse struct {
	Questions []QuestionResponse `json:"questions"`
}

// --- mappers --- //

func toCharacter(req AddCharacterRequest) Character {
	return Character{
		Character: req.Character,
		Reading:   req.Reading,
		Meaning:   req.Meaning,
		Category:  req.Category,
		Level:     req.Level,
	}
}

func toCharacterResponse(c Character) CharacterResponse {
	return CharacterResponse{
		ID:        c.ID,
		Character: c.Character,
		Reading:   c.Reading,
		Meaning:   c.Meaning,
		Category:  c.Category,
		Level:     c.Level,
	}
}

func toChallengeParams(req ChallengeRequest) ChallengeParams {
	return ChallengeParams{
		Category: req.Category,
		Level:    req.Level,
		Amount:   req.Amount,
	}
}

func toChallengeResponse(questions []Question) ChallengeResponse {
	out := make([]QuestionResponse, 0, len(questions))
	for _, q := range questions {
		out = append(out, QuestionResponse{
			Character: q.Character,
			Options:   q.Options,
			Answer:    q.Answer,
		})
	}
	return ChallengeResponse{Questions: out}
}
