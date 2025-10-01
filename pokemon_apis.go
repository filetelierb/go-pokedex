package main

import (
	"encoding/json"
	"fmt"
	"go-pokedex/internal/pokecache"
	"net/http"
	"time"
	"io"
	"bytes"
)
type Locations struct{
	Id int
	Name string
}
type LocationsResponse struct {
	Next string `json:"next"`
	Previous string `json:"previous"`
	Locations  []Locations 	`json:"results"`
}

type PokemonStatType struct {
	Name string `json:"name"`
}
type PokemonStats struct {
	BaseStat int `json:"base_stat"`
	Type PokemonStatType `json:"stat"`
}




type PokemonEncounter struct {
	Pokemon Pokemon `json:"pokemon"`
}

type PokemonType struct {
	Name string `json:"name"`
}
type PokemonTypeObj struct{
	Type PokemonType `json:"type"`
}

type Pokemon struct {
	Name string `json:"name"`
	Url string `json:"url"`
	BaseExperience int `json:"base_experience"`
	Height int `json:"height"`
	Weight int `json:"weight"`
	Stats []PokemonStats `json:"stats"`
	Types []PokemonTypeObj `json:"types"`

}

type SingleLocationResponse struct {
	Id int `json:"id"`
	Name string `json:"name"`
	PokemonEncounters []PokemonEncounter `json:"pokemon_encounters"`
}

var cache = pokecache.NewCache(300 * time.Second)

//POKEMON METHODS
func (pokemon Pokemon) PrintStats(){
	fmt.Printf("Name: %s\n",pokemon.Name)
	fmt.Printf("Height: %d\n",pokemon.Height)
	fmt.Printf("Weight: %d\n",pokemon.Weight)
	fmt.Print("Stats:\n")
	for _,stat := range pokemon.Stats{
		fmt.Printf("  - %s: %d\n",stat.Type.Name,stat.BaseStat)
	}
	fmt.Print("Types:\n")
	for _,t := range pokemon.Types{
		fmt.Printf("  - %s\n",t.Type.Name)
	}
}


func getLocationsAreas(configurations *config, page string) ([]Locations, error){
	locationsPageUrl := ""
	switch page {
	case "next":
		if configurations.Next == "" {
			fmt.Printf("You are in the last page")
		}
		locationsPageUrl = configurations.Next
	case "previous":
		if configurations.Previous == "" {
			fmt.Printf("You are in the first page")
		}
		locationsPageUrl = configurations.Previous
	}
	if locationsPageUrl == "" {
		return nil, fmt.Errorf("empty url")
	}

	var toDecode io.Reader
	var response LocationsResponse
	cacheData, ok := cache.Get(locationsPageUrl)
	if ok {
		toDecode = bytes.NewReader(cacheData)
	} else {
		req, err := http.NewRequest("GET",locationsPageUrl,nil)
		req.Close = true
		if err != nil {
			return nil, err
		}

		client := &http.Client{}
		res, err:=client.Do(req)
		if err != nil {
			return nil, err
		}
		defer res.Body.Close()
		toCache, err := io.ReadAll(res.Body)
		if err != nil {
			cache.Add(locationsPageUrl, toCache)
		}
		toDecode = bytes.NewReader(toCache)

	}

	
	

	decoder := json.NewDecoder(toDecode)

	if err := decoder.Decode(&response); err != nil{
		return nil, err
	}
	
	configurations.Next, configurations.Previous = response.Next, response.Previous

	return response.Locations, nil

	


}

func getLocationPokemons(locationName string) ([]Pokemon, error){
	if locationName == "" {
		return nil, fmt.Errorf("location's name is empty")
	}
	var toDecode io.Reader
	var singleLocationResp SingleLocationResponse
	var response []Pokemon
	url := "https://pokeapi.co/api/v2/location-area/" + locationName

	cacheData, ok := cache.Get(url)
	if ok {
		toDecode = bytes.NewReader(cacheData)

	} else {
		req, err := http.NewRequest("GET",url,nil)
		req.Close = true
		if err != nil {
			return nil, err
		}

		client := &http.Client{}
		res, err:=client.Do(req)
		if err != nil {
			return nil, err
		}
		defer res.Body.Close()
		toCache, err := io.ReadAll(res.Body)
		if err != nil {
			cache.Add(url, toCache)
		}
		toDecode = bytes.NewReader(toCache)
	}

	decoder := json.NewDecoder(toDecode)
	err := decoder.Decode(&singleLocationResp)
	if err != nil {
		return nil, err
	}
	
	for _,encounter := range singleLocationResp.PokemonEncounters {
		response = append(response, encounter.Pokemon)
	}



	return response, nil

}

func getPokemon (pokemonName string) (Pokemon, error){
	if pokemonName == "" {
		return Pokemon{}, fmt.Errorf("pokemon's name is empty")
	}
	var toDecode io.Reader
	var response Pokemon
	url := "https://pokeapi.co/api/v2/pokemon/" + pokemonName

	cacheData, ok := cache.Get(url)
	if ok {
		toDecode = bytes.NewReader(cacheData)
	} else {
		req, err := http.NewRequest("GET",url,nil)
		req.Close = true
		if err != nil {
			return Pokemon{}, err
		}

		client := &http.Client{}
		res, err:=client.Do(req)
		if err != nil {
			return Pokemon{}, err
		}
		defer res.Body.Close()
		toCache, err := io.ReadAll(res.Body)
		if err != nil {
			cache.Add(url, toCache)
		}
		toDecode = bytes.NewReader(toCache)
	}

	decoder := json.NewDecoder(toDecode)
	err := decoder.Decode(&response)
	if err != nil {
		return Pokemon{}, err
	}

	return response, nil

}