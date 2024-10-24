package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type CommandFlags struct {
	Add    string
	Edit   string
	Del    int
	Toggle int
	List   bool
}

func NewCmdFlags() *CommandFlags {
	cf := CommandFlags{}

	flag.StringVar(&cf.Add, "add", "", "Add a new entry")
	flag.StringVar(&cf.Edit, "edit", "", "Edit an existing new entry by specifying a new title id: new_title")
	flag.IntVar(&cf.Del, "del", -1, "Delete an existing an entry by specifying the index")
	flag.IntVar(&cf.Toggle, "toggle", -1, "Set an entry to completed/ongoing")
	flag.BoolVar(&cf.List, "list", false, "List all entries")

	flag.Parse()

	return &cf
}

// Execute is a receiver function that receives a pointer that points to the CommandFlags struct and takes in a
// pointer to the Todos which is a slice of Todo struct as an argument
// Because Execute returns a pointer of Todos,
// we have access to the other pointer receiver function that receives *Todos
// and it should be a pointer because we want to alter the actual thing, not a copy of it
func (cf *CommandFlags) Execute(todos *Todos) {
	switch {
	case cf.List:
		todos.print()
	case cf.Add != "":
		todos.add(cf.Add)
	case cf.Edit != "":
		parts := strings.SplitN(cf.Edit, ":", 2)
		if len(parts) != 2 {
			fmt.Println("Invalid edit. Please input id:New:title")
		}
		index, err := strconv.Atoi(parts[0])

		if err != nil {
			fmt.Println("Invalid index. Please input id:New:title")
			os.Exit(1)
		}
		todos.edit(index, parts[1])

	case cf.Del != -1:
		todos.delete(cf.Del)

	case cf.Toggle != -1:
		todos.toggle(cf.Toggle)

	default:
		fmt.Println("Invalid Command")
	}

}
