package main

import (
	"time"

	"github.com/cvc-comanescu-catalin/pokedexcli/internal/pokeapi"
)

func main() {
	pokeClient := pokeapi.NewClient(5 * time.Second)
	cfg := &config{
		pokeapiClient: pokeClient,
	}
	
	startRepl(cfg)
}