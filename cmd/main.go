package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/zonder12120/quotebook-autofill/pkg/env"
)

const (
	pathEnv            = ".env"
	defaultQuotesCount = 100
)

type Quote struct {
	Author string `json:"author"`
	Quote  string `json:"quote"`
}

var firstNames = []string{
	"Alan", "Grace", "Ada", "Linus", "Steve", "Marie", "Carl", "Niels", "Isaac", "Katherine",
}

var lastNames = []string{
	"Turing", "Hopper", "Lovelace", "Torvalds", "Jobs", "Curie", "Sagan", "Bohr", "Newton", "Johnson",
}

var verbs = []string{
	"creates", "destroys", "reveals", "hides", "changes", "simplifies", "complicates", "defines", "questions", "guides",
}

var adjectives = []string{
	"simple", "complex", "beautiful", "painful", "strange", "powerful", "weak", "hidden", "true", "false",
}

var nouns = []string{
	"life", "truth", "mind", "soul", "path", "world", "light", "darkness", "silence", "freedom",
}

func generateRandomQuote(r *rand.Rand) string {
	patterns := []string{
		"The %s %s the %s.",
		"%s is a %s %s.",
		"In the %s, the %s %s.",
		"%s %s %s and %s.",
		"%s cannot %s the %s.",
	}

	p := patterns[r.Intn(len(patterns))]

	words := []any{}
	for i := 0; i < strings.Count(p, "%s"); i++ {
		switch i % 3 {
		case 0:
			words = append(words, random(r, adjectives))
		case 1:
			words = append(words, random(r, verbs))
		case 2:
			words = append(words, random(r, nouns))
		}
	}

	return fmt.Sprintf(p, words...)
}

func generateRandomAuthor(r *rand.Rand) string {
	return fmt.Sprintf("%s %s", random(r, firstNames), random(r, lastNames))
}

func random(r *rand.Rand, arr []string) string {
	return arr[r.Intn(len(arr))]
}

func main() {
	err := env.LoadEnv(pathEnv)
	if err != nil {
		fmt.Println("Error loading .env file")
		panic(err)
		return
	}
	quotesCount, err := env.GetIntFromEnv("QUOTES_COUNT")
	if err != nil {
		fmt.Printf("[WARNING] Error parse from QUOTES_COUNT in .env file, expected int, got \"%v\", using default count: %d\n\n", quotesCount, defaultQuotesCount)
		quotesCount = defaultQuotesCount
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	url := "http://localhost:8080/quotes"

	for i := 1; i <= quotesCount; i++ {
		q := Quote{
			Author: generateRandomAuthor(r),
			Quote:  generateRandomQuote(r),
		}

		data, err := json.Marshal(q)
		if err != nil {
			fmt.Printf("[%d] JSON error: %v\n", i, err)
			continue
		}

		resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
		if err != nil {
			fmt.Printf("[%d] POST error: %v\n", i, err)
			continue
		}
		err = resp.Body.Close()
		if err != nil {
			fmt.Printf("[%d] response body close: %v\n", i, err)
			return
		}

		fmt.Printf("[%d] Sent: \"%s\" — %s\n", i, q.Quote, q.Author)
	}
}
