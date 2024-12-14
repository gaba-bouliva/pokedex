package main

import (
	"fmt"
	"strings"
)

func main () { 
  fmt.Println("Hello, World!")
}

func cleanInput(text string) []string { 
	words := strings.Split(text, " ")
	cleanWords := []string{}
	for _, word := range words {
		currentWord := strings.TrimSpace(word)
		if  currentWord != "" {
			cleanWords = append(cleanWords, currentWord)
		}
	}
	
	return cleanWords
}