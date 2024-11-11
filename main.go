package main

import (
	"cli-todo/cmd"
	"fmt"
	"os"
)

func main() {
	todos := Todos{}
	storage := NewStorage[Todos]("todos.json")
	operator := NewCoreFunctionOperator("todos.json")
	storage.Load(&todos)
	//cmdFlags := cmd.NewCmdFlags()
	//cmdFlags.Execute(operator)
	cmdRouter := cmd.NewCommandRouter(operator)
	err := cmdRouter.HandleArgs(os.Args)
	if err != nil {
		fmt.Println("cannot parse os.args")
		return
	}
	storage.Save(todos)
	operator.storage.Save(*operator.todos)

}
