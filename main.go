package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	pokecache "github.com/gaba-bouliva/pokedex-cli/internal/pokecache"
)

type cliCommand struct {
	name        string
	description string
	callback    func() error
	*config
}

type location struct {
	Name string				`json:"name"`
	URL  string				`json:"url"`
}

type PokemonEncounter struct { 
	Pokemon 					location 		`json:"pokemon"`
	VersionDetails		[]any				`json:"version_details"`
}

type locationArea struct { 
	ID       						int										`json:"id"`
	Name								string								`json:"name"`
	PokemonEncounters		[]PokemonEncounter		`json:"pokemon_encounters"`	
}

type config struct {
	BaseUrl   string
	Next      string     `json:"next"`
	Previous  string     `json:"previous"`
	Locations []location `json:"results"`
}

var cacheEntries = pokecache.NewCache(20)

var baseUrl = "https://pokeapi.co/api/v2"

var commands map[string]cliCommand

func main() {
	commands = getCommands()
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		usrInput := scanner.Text()
		cleanedInput := strings.ToLower(strings.TrimSpace(usrInput))
		cmdArgs := strings.Split(cleanedInput, " ")
		firstCmd := cmdArgs[0]
		if cmd, exists := commands[firstCmd]; !exists {
			fmt.Println("Unknown command")
		} else {
			if firstCmd == "explore" {
				if len(cmdArgs) < 2 {
					fmt.Println("Invalid use of command explore")
					fmt.Println("[Usage:] explore <location name>")
				} else {
					locationName := cmdArgs[1]
					cmd.callback = getExploreCmdCallback(locationName)
				}
			}
			if cmd.callback != nil {
				err := cmd.callback()
				if err != nil {
					fmt.Println(err)
				}
			}

		}

	}
}

func getExploreCmdCallback(name string) func() error {
	return func() error {
		_, ok := commands["explore"]
		if !ok {
			return fmt.Errorf("no command with name explore exists")
		}
		locaionArea := locationArea{}
		resource := "location-area"
		url := fmt.Sprintf("%s/%s/%s/", baseUrl, resource, name)
		for k := range cacheEntries.Entries { 
				fmt.Println(k)
		}
		jsonRes, exists := cacheEntries.Get(url)
		if exists {
			err := json.Unmarshal(jsonRes, &locaionArea)
			if err != nil {
				return err
			}
		} else {
			res, err := http.Get(url)
			if err != nil {
				return err
			}
			defer res.Body.Close()
			decoder := json.NewDecoder(res.Body)
			err = decoder.Decode(&locaionArea)
			if err != nil {
				fmt.Printf("error decoding json response: got error: %v\n", err)
				if res.StatusCode == 404 {
					return fmt.Errorf("location area %s not found", name)
				}
				return err
			}
			jsonVal, err := json.Marshal(&locaionArea)
			if err != nil {
				return err
			}
			cacheEntries.Add(url, jsonVal)
		}
	
		for _, encounter := range locaionArea.PokemonEncounters {
			fmt.Println(encounter.Pokemon.Name)
		}

		return nil
	}

}

func cmdMap() error {
	cmd, ok := commands["map"]
	if !ok {
		return fmt.Errorf("no command with name map exists")
	}

	var url string
	var conf config

	if cmd.Next != "" || len(cmd.Next) > len(cmd.BaseUrl) {
		url = cmd.Next
	} else {
		url = cmd.BaseUrl
	}

	if cachedRes, ok := cacheEntries.Get(url); ok {
		err := json.Unmarshal(cachedRes, &conf)
		if err != nil {
			return err
		}
	} else {
		res, err := http.Get(url)
		if err != nil {
			return err
		}
		decoder := json.NewDecoder(res.Body)
		err = decoder.Decode(&conf)
		if err != nil {
			return err
		}
		jsonVal, err := json.Marshal(&conf)
		if err != nil {
			return err
		}
		cacheEntries.Add(url, jsonVal)
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
	cmd, ok := commands["map"]
	if !ok {
		return fmt.Errorf("no command with name map exists")
	}
	if len(cmd.Previous) < len(cmd.BaseUrl) {
		return fmt.Errorf("no previous locations available")
	}
	var url string
	var conf config

	if cmd.Previous != "" || len(cmd.Previous) > len(cmd.BaseUrl) {
		url = cmd.Previous
	} else {
		url = cmd.BaseUrl
	}

	if cachedRes, ok := cacheEntries.Get(url); ok {
		fmt.Println("cached response for url: ", url)
		err := json.Unmarshal(cachedRes, &conf)
		if err != nil {
			return err
		}
	} else {
		res, err := http.Get(url)
		if err != nil {
			return err
		}

		decoder := json.NewDecoder(res.Body)
		err = decoder.Decode(&conf)
		if err != nil {
			return err
		}
		jsonVal, err := json.Marshal(&conf)
		if err != nil {
			return err
		}
		cacheEntries.Add(url, jsonVal)
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
		if currentWord != "" {
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
		"explore": {
			name:        "explore",
			description: "Explore pokemons in a location",
			config:      &config{},
		},
	}
}
