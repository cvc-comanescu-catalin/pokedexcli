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
}

type cliCommand struct {
	name        string
	description string
	callback    func(*config) error
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
func (c *Cli) RegisterCommand(cmd cliCommand) {
    c.commands[cmd.name] = cmd
}

func (c *Cli) RegisterAllCommands() {
	c.RegisterCommand(cliCommand{
        name:        "map",
        description: "Get the next page of locations",
        callback:    commandMapf,
    })
	c.RegisterCommand(cliCommand{
        name:        "mapb",
        description: "Get the previous page of locations",
        callback:    commandMapb,
    })
	c.RegisterCommand(cliCommand{
        name:        "exit",
        description: "Exit the Pokedex",
        callback:    commandExit,
    })
	c.CreateHelpCommand()
}

func (c *Cli) CreateHelpCommand() {
	c.RegisterCommand(cliCommand{
        name:        "help",
        description: "Displays a help message",
        callback:    func(*config) error {
			fmt.Print("Welcome to the Pokedex!\n")
			fmt.Print("Usage:\n\n")
			for name, command := range c.commands {
				fmt.Printf("%s: %s\n", name, command.description)
			}
			return nil
		},
    })
}

func (c *Cli) ExecuteCommand(name string, cfg *config) {
    command, ok := c.commands[name]
	if !ok {
		fmt.Println("Unknown command")
		return
	}

	err := command.callback(cfg)
	if err != nil {
		fmt.Println(err)
	}
}

func commandExit(*config) error {
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
			break;
		}
		
		enteredWords := cleanInput(scanner.Text())
		if (len(enteredWords) == 0) {
			continue
		}

		command := enteredWords[0]
		cli.ExecuteCommand(command, cfg)
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