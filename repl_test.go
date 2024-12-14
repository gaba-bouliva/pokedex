package main

import (
	"fmt"
	"testing"
)

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "  hello  world  ",
			expected: []string{"hello", "world"},
		},
		{
			input: 		"Charmander Bulbasaur PIKACHU",
			expected: []string{"Charmander", "Bulbasaur", "PIKACHU"},
		},
	}
	for _, c := range cases {
		actual := cleanInput(c.input)
		fmt.Println("length of actual:", len(actual))
		fmt.Println("length of expected:", len(c.expected))
		if len(actual) != len(c.expected) {
			t.Errorf("Test failed expected %v but got %v", c.expected, actual)
			t.Fail()
		}
		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]
			if word != expectedWord { 
				t.Errorf("Test failed expected %v but got %v", c.expected, actual)
				t.Fail()
			}
		}
	}
}