package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/cvc-comanescu-catalin/pokedexcli/internal/pokeapi"
)

type config struct {
	pokeapiClient    pokeapi.Client
	nextLocationsURL *string
	prevLocationsURL *string
	caughtPokemon    map[string]pokeapi.Pokemon
}

type cliCommand struct {
	name        string
	description string
	callback    func(*config, ...string) error
}

type Cli struct {
	commands map[string]cliCommand
}

func NewCli() *Cli {
	return &Cli{
		commands: make(map[string]cliCommand),
	}
}

// Method to register commands
func (c *Cli) RegisterCommand(cmdName string, cmd cliCommand) {
	c.commands[cmdName] = cmd
}

func (c *Cli) RegisterAllCommands() {
	c.RegisterCommand(
		"explore",
		cliCommand{
			name:        "explore <location_name>",
			description: "Explore a location",
			callback:    commandExplore,
		})
	c.RegisterCommand(
		"catch",
		cliCommand{
			name:        "catch <pokemon_name>",
			description: "Attempt to catch a Pokemon",
			callback:    commandCatchPokemon,
		})
	c.RegisterCommand(
		"inspect",
		cliCommand{
			name:        "inspect <pokemon_name>",
			description: "View details about a caught Pokemon",
			callback:    commandInspect,
		})
	c.RegisterCommand(
		"map",
		cliCommand{
			name:        "map",
			description: "Get the next page of locations",
			callback:    commandMapf,
		})
	c.RegisterCommand(
		"mapb",
		cliCommand{
			name:        "mapb",
			description: "Get the previous page of locations",
			callback:    commandMapb,
		})
	c.RegisterCommand(
		"pokedex",
		cliCommand{
			name:        "pokedex",
			description: "See all the pokemon you've caught",
			callback:    commandPokedex,
		})
	c.RegisterCommand(
		"exit",
		cliCommand{
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		})
	c.CreateHelpCommand()
}

func (c *Cli) CreateHelpCommand() {
	c.RegisterCommand(
		"help",
		cliCommand{
			name:        "help",
			description: "Displays a help message",
			callback: func(*config, ...string) error {
				fmt.Print("Welcome to the Pokedex!\n")
				fmt.Print("Usage:\n\n")
				for name, command := range c.commands {
					fmt.Printf("%s: %s\n", name, command.description)
				}
				return nil
			},
		})
}

func (c *Cli) ExecuteCommand(name string, cfg *config, args ...string) {
	command, ok := c.commands[name]
	if !ok {
		fmt.Println("Unknown command")
		return
	}

	err := command.callback(cfg, args...)
	if err != nil {
		fmt.Println(err)
	}
}

func commandExit(*config, ...string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func startRepl(cfg *config) {
	cli := NewCli()
	cli.RegisterAllCommands()

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		if scanned := scanner.Scan(); !scanned {
			break
		}

		enteredWords := cleanInput(scanner.Text())
		if len(enteredWords) == 0 {
			continue
		}

		command := enteredWords[0]
		commandArgs := []string{}
		if len(enteredWords) > 1 {
			commandArgs = enteredWords[1:]
		}

		cli.ExecuteCommand(command, cfg, commandArgs...)
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
}

func cleanInput(text string) []string {
	output := strings.ToLower(text)
	words := strings.Fields(output)
	return words
}
