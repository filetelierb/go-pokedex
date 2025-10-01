package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"
)

func commandExit (configurations *config, parameters []string) error{
	fmt.Print("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp (configurations *config, parameters []string) error{
	fmt.Print("Welcome to the Pokedex!\nUsage:\n\n")
	for _,val := range cliCommands{
		fmt.Printf("%s: %s\n",val.name,val.description)
	}
	return nil
}

func commandMap (configurations *config, parameters []string) error{
	locations, err := getLocationsAreas(configurations, "next")
	if err != nil {
		fmt.Printf("There was an error fetching the locations %v\n", err)
		return err
	}
	for _,location := range locations {
		fmt.Printf("%s\n",location.Name)
	}
	return nil
}

func commandMapb (configurations *config, parameters []string) error{
	locations, err := getLocationsAreas(configurations, "previous")
	if err != nil {
		fmt.Printf("There was an error fetching the locations: %v\n", err)
		return err
	}
	for _,location := range locations {
		fmt.Printf("%s\n",location.Name)
	}
	return nil
}

func commandExplore (configurations *config, parameters []string) error{
	areaName := parameters[0]
	fmt.Printf("Exploring %s...\nFound Pokemon:\n",areaName)

	pokemons, err := getLocationPokemons(areaName)
	
	if err != nil {
		fmt.Printf("There was an error exploring the area %s: %v\n", areaName, err)
		return err
	}
	
	if len(pokemons) == 0 {
		fmt.Printf("No pokemons were found in %s\n",areaName)
	} else {
		for _,pokemon := range pokemons {
			fmt.Printf("- %s\n", pokemon.Name)
		}
	}
	return nil
}

func commandCatch (configurations *config, parameters []string) error{

	if len(parameters) == 0{
		return fmt.Errorf("please pass a valid pokemon name as parameter")
	}
	pokemonName := parameters[0]
	_, hasPokemon := capturedPokemons[pokemonName]
	if hasPokemon {
		fmt.Printf("You already captured %s\n", pokemonName)
		return nil
	}
	fmt.Printf("Throwing a Pokeball at %s...\n",pokemonName)
	pokemon, err := getPokemon(pokemonName)
	
	if err != nil {
		return err
	}
	baseExp := float64(pokemon.BaseExperience)
	baseChance := 0.5 - math.Min((baseExp / 1000),0.49)
	if success := (float64(rand.Intn(100)) / 100) < baseChance; success {
		capturedPokemons[pokemonName] = pokemon
		fmt.Printf("%s was caught!\nYou may now inspect it with the inspect command.\n", pokemonName)
	} else {
		fmt.Printf("%s escaped!\n", pokemonName)
	}
	return nil

}

func commandInspect (configurations *config, parameters []string) error{
	if len(parameters) == 0{
		return fmt.Errorf("please pass a valid pokemon name as parameter")
	}
	pokemonName := parameters[0]
	pokemon, hasPokemon := capturedPokemons[pokemonName]
	if !hasPokemon {
		fmt.Print("you have not caught that pokemon")
		return nil
	}
	pokemon.PrintStats()
	return nil
	
}

func commandPokedex (configurations *config, parameters []string) error{
	fmt.Print("Your Pokedex:\n")
	for _, pokemon := range capturedPokemons{
		fmt.Printf(" - %s\n",pokemon.Name)
	}
	return nil
}

type config struct {
	Next string
	Previous string
}

type cliCommand struct {
	name        string
	description string
	callback    func(configurations *config, parameters []string) error
}

var cliCommands map[string]cliCommand
var initialConfigs config
var capturedPokemons map[string]Pokemon

func init() {

	initialConfigs = config {
		Next: "https://pokeapi.co/api/v2/location-area/?limit=20",
		Previous: "",
	}

	capturedPokemons = make(map[string]Pokemon)

    cliCommands = map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
        "exit": {
            name:        "exit",
            description: "Exit the Pokedex",
            callback:    commandExit,
        },
		"map": {
			name: "map",
			description: "Fetch next page of locations",
			callback: commandMap,
		},
		"mapb": {
			name: "mapb",
			description: "Fetch the previous page of locations",
			callback: commandMapb,
		},
		"explore": {
			name: "explore",
			description: "It explores a location fetching pokemons",
			callback: commandExplore,
		},
		"catch": {
			name: "catch",
			description: "It attempts to catch the specified pokemon",
			callback: commandCatch,
		},
		"inspect": {
			name: "inspect",
			description: "It looks for a captured pokemon",
			callback: commandInspect,
		},
		"pokedex": {
			name: "pokedex",
			description: "Shows captured pokemons",
			callback: commandPokedex,
		},
    }
}
