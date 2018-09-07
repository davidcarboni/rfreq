package main

import (
	"fmt"
	"log"
	"os"
	"text/scanner"
)

func parse(filename string, corpus map[string]int) map[string]int {
	s := scanner.Scanner{}
	file, err := os.Open(filename)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	s.Init(file)
	words := make(map[string]int)
	for token := s.Scan(); token != scanner.EOF; token = s.Scan() {
		words[s.TokenText()] = words[s.TokenText()] + 1
		corpus[s.TokenText()] = corpus[s.TokenText()] + 1
	}

	//fmt.Printf("%s: %v\n", filename, words)
	return words
}

func frequencies(words map[string]int) map[string]float64 {

	// Get the total number of words
	var total float64
	for _, count := range words {
		total += float64(count)
	}

	// Compute the frequency of each word
	frequencies := make(map[string]float64)
	var percent float64
	for word, count := range words {
		frequencies[word] = float64(count) / total
		percent += float64(count) / total
	}

	fmt.Println()
	fmt.Println(percent)
	fmt.Println()
	return frequencies
}

func score(word string, words map[string]float64, corpus map[string]float64) float64 {
	wordFrequency := words[word]
	corpusFrequency := corpus[word]
	return wordFrequency - corpusFrequency
}

func top10(words map[string]float64, corpus map[string]float64) []string {
	var top10words []string
	var top10scores []float64

	for word := range words {
		score := score(word, words, corpus)

		// Work out the ranking:
		index := 0
		for _, top10score := range top10scores {
			if score < top10score {
				index++
			}
		}

		// Insert into the top 10
		if index < 10 {
			top10words = append(top10words, "")
			copy(top10words[index+1:], top10words[index:])
			top10words[index] = word
			top10scores = append(top10scores, 0)
			copy(top10scores[index+1:], top10scores[index:])
			top10scores[index] = score
		}
	}

	return top10words[:10]
}

func main() {
	corpus := make(map[string]int)
	parse("text/agents.txt", corpus)
	w := parse("text/culture.txt", corpus)
	parse("text/lazy.txt", corpus)
	//parse("text/test.txt", corpus)
	fmt.Println(corpus)

	corpusFreq := frequencies(corpus)
	agentsFreq := frequencies(w)

	for word := range w {
		fmt.Printf("%s: %f", word, score(word, agentsFreq, corpusFreq))
	}

	for _, topWord := range top10(agentsFreq, corpusFreq) {
		fmt.Println(topWord)
	}

}
