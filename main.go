package main

func main() {
	todos := Todos{}
	storage := NewStorage[Todos]("todos.json")
	operator := NewCoreFunctionOperator("todos.json")
	storage.Load(&todos)
	cmdFlags := NewCmdFlags()
	cmdFlags.Execute(operator, &todos)
	//storage.Save(todos)
}
