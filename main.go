package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

type cliCommand struct { 
	name 						string
	description			string
	callback 				func()error
	*config
}

type location struct { 
	Name 				string
	URL 				string				
}

type config struct { 
	BaseUrl 			string
	Next					string			`json:"next"`
	Previous			string			`json:"previous"`
	Locations 		[]location 	`json:"results"`
}

var commands map[string]cliCommand

func main () { 
	commands = getCommands()
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		usrInput := scanner.Text()
		cleanedInput := strings.ToLower(strings.TrimSpace(usrInput))
		firstCmd := strings.Split(cleanedInput, " ")[0]
		if cmd, exists := commands[firstCmd]; !exists { 
			fmt.Println("Unknown command")
		} else {	
			err := cmd.callback()
			if err != nil {
				fmt.Println(err)
			}		
			}
			
	}
}

func cmdMap() error { 
	cmd,ok := commands["map"]
	if !ok { 
		return fmt.Errorf("no command with name map exists")
	}
	var url string
	if cmd.Next != "" || len(cmd.Next) > len(cmd.BaseUrl) {
		url = cmd.Next	
	} else { 
		url = cmd.BaseUrl
	}
	res, err := http.Get(url)
	if err != nil { 
		return err
	}
	decoder := json.NewDecoder(res.Body)
	var conf config	
	err = decoder.Decode(&conf)
	if err != nil { 
		return err
	}
	cmd.Next = conf.Next
	cmd.Previous = conf.Previous
	cmd.Locations = conf.Locations
	for _, location := range cmd.Locations { 
		fmt.Println(location.Name)
	}
	return nil
}

func cmdBackMap() error { 
	cmd,ok := commands["map"]
	if !ok { 
		return fmt.Errorf("no command with name map exists")
	}
	if len(cmd.Previous) < len(cmd.BaseUrl) { 
		return fmt.Errorf("no previous locations available")
	}
	var url string
	if cmd.Previous != "" || len(cmd.Previous) > len(cmd.BaseUrl) {
		url = cmd.Previous	
	} else { 
		url = cmd.BaseUrl
	}
	res, err := http.Get(url)
	if err != nil { 
		return err
	}

	decoder := json.NewDecoder(res.Body)
	var conf config	
	err = decoder.Decode(&conf)
	if err != nil { 
		return err
	}
	cmd.Next = conf.Next
	cmd.Previous = conf.Previous
	cmd.Locations = conf.Locations
	for _, location := range cmd.Locations { 
		fmt.Println(location.Name)
	}
	return nil
}

func cmdExit() error { 
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func cmdHelp() error { 
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()
	fmt.Println()	 		
	for _, v := range getCommands() { 
			fmt.Printf("%s: %s\n", v.name, v.description)
		}
	return nil
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

func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    cmdHelp,
			config: &config{
				BaseUrl: "https://pokeapi.co/api/v2/location-area/",
			},
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    cmdExit,
			config: &config{
				BaseUrl: "https://pokeapi.co/api/v2/location-area/",
			},
		},
		"map": {
			name:        "map",
			description: "Display Pokedex Next 20 Locations",
			callback:    cmdMap,
			config: &config{
				BaseUrl: "https://pokeapi.co/api/v2/location-area/",
			},
		},
		"mapb": {
			name:        "mapb",
			description: "Display Pokedex Previous 20 Locations",
			callback:    cmdBackMap,
			config: &config{
				BaseUrl: "https://pokeapi.co/api/v2/location-area/",
			},
		},
	}
}