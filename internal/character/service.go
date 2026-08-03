package character

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
)

var (
	ErrInvalidCategory = errors.New("invalid category")
	ErrInvalidLevel    = errors.New("invalid level")
	ErrInvalidAmount   = errors.New("amount must be greater than 0")
	ErrNotEnoughChars  = errors.New("not enough characters in the database for this filter")
)

var validCategories = map[string]bool{
	"hiragana": true,
	"katakana": true,
	"kanji":    true,
	"random":   true,
}

var validLevels = map[string]bool{
	"A1": true,
	"A2": true,
	"B1": true,
	"B2": true,
}

type CharacterService interface {
	GenerateChallenge(ctx context.Context, req ChallengeRequest) (ChallengeResponse, error)
	AddCharacter(ctx context.Context, c Character) (Character, error)
}

type characterService struct {
	repo CharacterRepository
}

func NewCharacterService(repo CharacterRepository) CharacterService {
	return &characterService{repo: repo}
}

func (s *characterService) AddCharacter(ctx context.Context, c Character) (Character, error) {
	if c.Character == "" || c.Reading == "" || c.Meaning == "" {
		return Character{}, fmt.Errorf("character, reading and meaning are required")
	}
	if !validCategories[c.Category] || c.Category == "random" {
		return Character{}, ErrInvalidCategory
	}
	if !validLevels[c.Level] {
		return Character{}, ErrInvalidLevel
	}
	return s.repo.Create(ctx, c)
}

func (s *characterService) GenerateChallenge(ctx context.Context, req ChallengeRequest) (ChallengeResponse, error) {
	// validate request
	if !validCategories[req.Category] {
		return ChallengeResponse{}, ErrInvalidCategory
	}
	if !validLevels[req.Level] {
		return ChallengeResponse{}, ErrInvalidLevel
	}
	if req.Amount <= 0 {
		return ChallengeResponse{}, ErrInvalidAmount
	}

	// fetch the question characters
	chars, err := s.repo.GetRandom(ctx, req.Category, req.Level, req.Amount)
	if err != nil {
		return ChallengeResponse{}, fmt.Errorf("fetching characters: %w", err)
	}
	if len(chars) < req.Amount {
		return ChallengeResponse{}, ErrNotEnoughChars
	}

	// collect IDs to exclude from distractors
	excludeIDs := make([]int64, len(chars))
	for i, c := range chars {
		excludeIDs[i] = c.ID
	}

	// fetch 2 distractors per question — we pick from this pool per question
	// 2 wrong options + 1 correct = 3 total choices
	distractors, err := s.repo.GetDistractors(ctx, excludeIDs, req.Amount*2)
	if err != nil {
		return ChallengeResponse{}, fmt.Errorf("fetching distractors: %w", err)
	}

	// build questions
	questions := make([]ChallengeQuestion, 0, len(chars))
	for i, c := range chars {
		question, err := buildQuestion(c, distractors, i)
		if err != nil {
			return ChallengeResponse{}, err
		}
		questions = append(questions, question)
	}

	return ChallengeResponse{Questions: questions}, nil
}

// buildQuestion builds a single question with 3 shuffled options.
// The correct answer is randomly either the reading or the meaning.
func buildQuestion(c Character, distractors []Character, index int) (ChallengeQuestion, error) {
	// randomly decide if the correct answer is the reading or the meaning
	var answer string
	answer = getRandomAnswerType(c)

	// from the index of the char list, get the index position of the distractor list
	// index = 0 then distractor will get the 0 and 1 index position of distractor list
	d1Index := index * 2
	d2Index := index*2 + 1

	// if index goes over the len of the distractor list then throw error
	if d2Index >= len(distractors) {
		return ChallengeQuestion{}, ErrNotEnoughChars
	}

	// wrong options — also randomly reading or meaning from distractor characters
	wrong1 := getRandomAnswerType(distractors[d1Index])
	wrong2 := getRandomAnswerType(distractors[d2Index])

	// shuffle the 3 options so the correct answer isn't always in the same position
	options := []string{answer, wrong1, wrong2}
	rand.Shuffle(len(options), func(i, j int) {
		options[i], options[j] = options[j], options[i]
	})

	return ChallengeQuestion{
		Character: c.Character,
		Options:   options,
		Answer:    answer,
	}, nil
}

func getRandomAnswerType(c Character) (answer string) {
	if rand.Intn(2) == 0 {
		answer = c.Reading
	} else {
		answer = c.Meaning
	}

	return answer
}
