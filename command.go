package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type CommandFlags struct {
	Add      string
	Edit     string
	Del      int
	Toggle   int
	List     bool
	Priority string
	Sort     string
	Filter   int
}

func NewCmdFlags() *CommandFlags {
	cf := CommandFlags{}

	flag.StringVar(&cf.Add, "add", "", "Add a new entry")
	flag.StringVar(&cf.Edit, "edit", "", "Edit an existing new entry by specifying a new title id: new_title")
	flag.IntVar(&cf.Del, "del", -1, "Delete an existing an entry by specifying the index")
	flag.IntVar(&cf.Toggle, "toggle", -1, "Set an entry to completed/ongoing")
	flag.BoolVar(&cf.List, "list", false, "List all entries")
	flag.StringVar(&cf.Sort, "sort", "", "sort all entries by either time or priority, expect input between 1 - 5")
	flag.IntVar(&cf.Filter, "filter", -1, "Filter by priority")
	flag.StringVar(&cf.Priority, "priority", "", "Set priority for the task at hand, "+
		"{first integer for the line index:second integer for priority level} {1:5}")

	flag.Parse()

	return &cf
}

// Execute is a receiver function that receives a pointer that points to the CommandFlags struct and takes in a
// pointer to the Todos which is a slice of Todo struct as an argument
// Because Execute returns a pointer of Todos,
// we have access to the other pointer receiver function that receives *Todos
// and it should be a pointer because we want to alter the actual thing, not a copy of it
// I know the cf.Sort refers to the command line input aka go run ./ - sort abc -> and it takes in abc but why
// flag.StringVar store the input abc into the pointer of &cf.sort?
func (cf *CommandFlags) Execute(todos *Todos) {
	switch {
	case cf.List:
		todos.print()
	case cf.Filter != -1:
		todos.filterByPriority(cf.Filter)
	case cf.Sort != "":
		todos.sort(cf.Sort)
	case cf.Add != "":
		todos.add(cf.Add)
		fmt.Printf("Added %s\n", cf.Add)
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
		fmt.Printf("Editted %s for index %s\n", parts[1], parts[0])

	case cf.Del != -1:
		todos.delete(cf.Del)
		fmt.Printf("Deleted %s\n", cf.Del)

	case cf.Toggle != -1:
		todos.toggle(cf.Toggle)

	case cf.Priority != "":
		parts := strings.SplitN(cf.Priority, ":", 2)
		if len(parts) != 2 {
			fmt.Println("Invalid priority. Please input [Int 1]:[Int 2]")
		}
		int1, _ := strconv.Atoi(parts[0])
		int2, _ := strconv.Atoi(parts[1])
		todos.setPriority(int1, int2)
		fmt.Printf("Set priority for index %d to %d\n", int1, int2)
	default:
		fmt.Println("Invalid Command")
	}

}
