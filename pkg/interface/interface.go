package interfaces

// anything that wants to be a ToDoOperator must have these functions
type TodoOperator interface {
	AddTodo(title string) error
	DeleteTodo(id int) error
	EditTodo(id int, title string) error
	ToggleTodo(id int) error
	SetPriority(index int, input int) error
	Sort(option string) error
	FilterByPriority(priority int) error
	SetTags(index int, Tags string) error
	DelTags(index int, Tags string) error
	Print()
}
