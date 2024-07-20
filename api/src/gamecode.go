package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
)

func NewGameCode() string {
	var output strings.Builder
	charSet := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	length := 4

	for i := 0; i < length; i++ {
		random := rand.Intn(len(charSet))
		randomChar := charSet[random]
		output.WriteString(string(randomChar))
	}

	return output.String()
}

func readIntoLines(filePath string, prefix string) []string {
	lines := []string{}

	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text := scanner.Text()

		if strings.HasPrefix(strings.ToUpper(text), strings.ToUpper(prefix)) {
			lines = append(lines, text)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}

	return lines
}

func capitalizeFirst(s string) string {
	if len(s) == 0 {
		return s
	}

	return strings.ToUpper(string(s[0])) + string(s[1:])
}

func pickRandom(words []string, deflt string) string {
	if len(words) == 0 {
		return deflt
	}

	return words[rand.Intn(len(words))]
}

func MakeGamePneumonic(gameCode string) string {
	adjectives := readIntoLines("static/words.txt", string(gameCode[0]))
	nounsA := readIntoLines("static/words.txt", string(gameCode[1]))
	verbs := readIntoLines("static/words.txt", string(gameCode[2]))
	nounsB := readIntoLines("static/words.txt", string(gameCode[3]))

	adjective := pickRandom(adjectives, string(gameCode[0]))
	nounA := pickRandom(nounsA, string(gameCode[1]))
	verb := pickRandom(verbs, string(gameCode[2]))
	nounB := pickRandom(nounsB, string(gameCode[3]))

	return fmt.Sprintf(
		"%s %s %s %s",
		capitalizeFirst(adjective),
		capitalizeFirst(nounA),
		capitalizeFirst(verb),
		capitalizeFirst(nounB),
	)
}
