package interfaces

// anything that wants to be a ToDoOperator must have these functions
type TodoOperator interface {
	AddTodo(title string) error
	DeleteTodo(id int) error
	EditTodo(id int, title string) error
	ToggleTodo(id int) error
}
