package cmd

import (
	interfaces "cli-todo/pkg/interface"
	"errors"
	"fmt"
	"strings"
)

//var rootCmd *flag.FlagSet

type CommandRouter struct {
	commands map[string]Command
}

var standardCommands = map[string]func(interfaces.TodoOperator) Command{
	"add": NewAddCommand,
	// easy to add more commands
}

func NewAddCommand(operator interfaces.TodoOperator) Command {
	return &AddCommand{
		BaseCommand: BaseCommand{
			Name:      "add",
			UsageText: "Usage: todo add <text>",
			operator:  operator,
		},
	}
}

func NewCommandRouter(operator interfaces.TodoOperator) *CommandRouter {
	router := &CommandRouter{
		commands: make(map[string]Command),
	}

	// Register all standard commands
	router.registerCommands(operator)
	return router
}

func (r *CommandRouter) registerCommands(operator interfaces.TodoOperator) {
	for name, createCmd := range standardCommands {
		r.commands[name] = createCmd(operator)
	}
}

type Command interface {
	Run() error
	HandleArgs(args []string) error
}

// BaseCommand struct to provide some common fields
type BaseCommand struct {
	Name      string
	UsageText string
	operator  interfaces.TodoOperator
}

type AddCommand struct {
	BaseCommand
	todoText string
}

func (a *AddCommand) Run() error {
	fmt.Println("Run() is called")
	if a.todoText != "" {
		err := a.operator.AddTodo(a.todoText)
		if err != nil {
			return err
		}
		return nil
	}
	err := fmt.Errorf("empty string")
	return err
}

// In CommandRouter package
func (r *CommandRouter) HandleArgs(args []string) error {
	fmt.Println(args)
	if len(args) < 2 {
		return errors.New("not enough arguments")
	}

	cmdName := args[1]
	fmt.Println(cmdName)
	cmd, exists := r.commands[cmdName]
	if !exists {
		return errors.New("unknown command")
	}

	if err := cmd.HandleArgs(args[2:]); err != nil {
		return err
	} // Pass remaining args to command
	return cmd.Run()
}

// In AddCommand
func (a *AddCommand) HandleArgs(args []string) error {
	// Now only handles command-specific args
	// No need to check "todo" or command name
	fmt.Println("Add command got called")
	if len(args) == 0 {
		return errors.New("no text provided")
	}
	a.todoText = strings.Join(args, " ")
	fmt.Println(a.todoText)
	return nil
}

//func newFlag() {
//	rootCmd = flag.NewFlagSet("todo", flag.ExitOnError)
//	// String defines a string flag with specified name, default value, and usage string.
//	// The return value is the address of a string variable that stores the value of the flag.
//	todoDelete := rootCmd.Int("del", -1, "del")
//	todoEdit := rootCmd.String("edit", "", "edit")
//	todoToggle := rootCmd.Int("toggle", -1, "toggle")
//
//}

//func todoAddCmd(args []string) error {
//	var textAdded string
//	flagSet := flag.NewFlagSet("add", flag.ExitOnError)
//	flagSet.StringVar(&textAdded, "add", "", "Add a new entry")
//	return nil
//}

//step 1 create a root command that can have subcommands
//each current flag operation should becme its own subcommand
//each subcommand can have its own flags for additional options
