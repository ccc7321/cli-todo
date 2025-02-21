package main

import (
	interfaces "cli-todo/pkg/interface"
	"errors"
	"fmt"
	"github.com/aquasecurity/table"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Todo struct {
	Title     string
	Completed bool
	CreatedAt time.Time
	//CompletedAt is a pointer because it can be nil, i.e. when it is not completed
	CompletedAt *time.Time
	Priority    int
	Tags        []string
}

// CoreFunctionOperator implements this contract or concrete implementation of the interface
type CoreFunctionOperator struct {
	todos   *Todos // Your todo list
	storage interfaces.StorageInterface[Todos]
}

func NewCoreFunctionOperator(filePath string) *CoreFunctionOperator {
	c := &CoreFunctionOperator{
		todos:   &Todos{},
		storage: NewStorage[Todos](filePath),
	}
	c.storage.Load(c.todos)
	return c
}

// Implement TodoOperator interface methods:
// the add function is a pointer receiver function as it receives CoreFunctionOperator
func (c *CoreFunctionOperator) AddTodo(title string) error {
	fmt.Printf("Checking if title is being passed: %s\n", title)
	fmt.Printf("Checking what the todo struct looks like here in CoreFunctionOperator before: %v\n", *c.todos)
	todo := Todo{
		Title:       title,
		Completed:   false,
		CreatedAt:   time.Now(),
		CompletedAt: nil,
	}
	*c.todos = append(*c.todos, todo)
	fmt.Printf("Checking what the todo struct looks like in CoreFunctionOperator after: %v\n", *c.todos)

	err := c.storage.Save(*c.todos)
	if err != nil {
		return fmt.Errorf("failed to save: %w", err)
	}

	// 4. Return any errors
	return nil
}

func (c *CoreFunctionOperator) DeleteTodo(index int) error {
	// 1. Validate id
	// 2. Delete from todos
	err := c.todos.validateIndex(index)
	if err != nil {
		return err
	}

	todosList := *c.todos

	*c.todos = append(todosList[:index], todosList[index+1:]...)
	// 3. Save changes
	err = c.storage.Save(*c.todos)
	if err != nil {
		return fmt.Errorf("failed to save: %w", err)
	}
	// 4. Return any errors
	return nil
}

func (c *CoreFunctionOperator) EditTodo(index int, title string) error {
	// 1. Validate id and title
	// 2. Edit todo
	t := *c.todos
	err := t.validateIndex(index)
	if err != nil {
		return err
	}

	t[index].Title = title
	// 3. Save changes
	err = c.storage.Save(*c.todos)
	if err != nil {
		return fmt.Errorf("failed to save: %w", err)
	}
	// 4. Return any errors
	return nil
}

func (c *CoreFunctionOperator) ToggleTodo(id int) error {
	fmt.Printf("id is %d\n", id)
	fmt.Printf("todos is %v\n", *c.todos)
	err := c.storage.Load(c.todos)
	fmt.Printf("todos is %v\n", *c.todos)
	if err != nil {
		return fmt.Errorf("Something wrong")
	}

	err = c.todos.validateIndex(id)
	if err != nil {
		return err
	}
	fmt.Printf("First %v", (*c.todos)[id].Completed)
	isCompleted := (*c.todos)[id].Completed
	if !isCompleted {
		completionTime := time.Now()
		(*c.todos)[id].CompletedAt = &completionTime
		fmt.Printf("2nd %v", (*c.todos)[id].Completed)
	}

	(*c.todos)[id].Completed = !(*c.todos)[id].Completed
	fmt.Printf("3rd %v", (*c.todos)[id].Completed)

	c.storage.Save(*c.todos)
	if err != nil {
		return fmt.Errorf("failed to save: %w", err)
	}
	return nil
}

func (c *CoreFunctionOperator) Print() {
	c.todos.print()
}

func (c *CoreFunctionOperator) SetPriority(id int, input int) error {
	c.todos.setPriority(id, input)
	err := c.storage.Save(*c.todos)
	if err != nil {
		return fmt.Errorf("failed to save: %w", err)
	}
	return nil
}

func (c *CoreFunctionOperator) Sort(option string) error {
	c.todos.sort(option)
	err := c.storage.Save(*c.todos)
	if err != nil {
		return fmt.Errorf("failed to save: %w", err)
	}
	return nil
}

func (c *CoreFunctionOperator) FilterByPriority(id int) error {
	c.todos.filterByPriority(id)
	err := c.storage.Save(*c.todos)
	if err != nil {
		return fmt.Errorf("failed to save: %w", err)
	}
	return nil
}

func (c *CoreFunctionOperator) SetTags(id int, tags string) error {
	c.todos.setTags(id, tags)
	err := c.storage.Save(*c.todos)
	if err != nil {
		return fmt.Errorf("failed to save: %w", err)
	}
	return nil
}

func (c *CoreFunctionOperator) DelTags(id int, tags string) error {
	c.todos.delTags(id, tags)
	err := c.storage.Save(*c.todos)
	if err != nil {
		return fmt.Errorf("failed to save: %w", err)
	}
	return nil
}

type Todos []Todo

func (todos *Todos) validateIndex(index int) error {
	if index < 0 || index >= len(*todos) {
		err := fmt.Errorf("invalid index, actual lebgth: %d", len(*todos))
		fmt.Println(err)
		return err
	}
	return nil
}

func (todos *Todos) print() {
	table := table.New(os.Stdout)
	table.SetRowLines(false)
	table.SetHeaders("#", "Title", "Completed", "Created At", "Completed At", "Priority", "Tags")
	for index, t := range *todos {
		completed := "Ongoing"
		completedAt := ""

		if t.Completed {
			completed = "Completed"
			if t.CompletedAt != nil {
				completedAt = t.CompletedAt.Format(time.RFC1123)
			}
		}
		table.AddRow(strconv.Itoa(index), t.Title, completed, t.CreatedAt.Format(time.RFC1123), completedAt,
			strconv.Itoa(t.Priority), strings.Join(t.Tags, ", "))
	}
	table.Render()
}

func (todos *Todos) setPriority(index int, input int) error {
	t := *todos
	err := t.validateIndex(index)
	if err != nil {
		return err
	}

	if 1 < input && input > 5 {
		return errors.New("Priority level is between 1 and 5")
	}

	t[index].Priority = input
	return nil
}

func (todos *Todos) sort(option string) error {
	t := *todos
	//sort by time
	if option == "time" {
		sort.Slice(t, func(i, j int) bool {
			return t[i].CreatedAt.Before(t[j].CreatedAt)
		})
	} else if option == "priority" {
		sort.Slice(t, func(i, j int) bool { return t[i].Priority < t[j].Priority })
		fmt.Println("Sorted by priority:")
		return nil
	}
	return errors.New("invalid option")
}

func (todos *Todos) filterByPriority(priority int) error {
	t := *todos
	filteredt := t
	filterHolder := Todos{}
	//filter by date
	if priority > 5 {
		return errors.New("invalid priority")
	}
	for _, t := range filteredt {
		if t.Priority == priority {
			filterHolder = append(filterHolder, t)
		}
	}

	fmt.Println("Filtered by priority:")
	filterHolder.print()
	return nil
}

func (todos *Todos) setTags(index int, Tags string) error {
	t := *todos
	err := t.validateIndex(index)
	if err != nil {
		return err
	}

	//identify the commas between using regex and put it into slices
	commaSplitter := regexp.MustCompile(`\s*,\s*`)
	parts := commaSplitter.Split(Tags, -1)
	t[index].Tags = parts
	return nil
}

func (todos *Todos) delTags(index int, Tags string) error {
	t := *todos
	err := t.validateIndex(index)
	tagMatched := false
	if err != nil {
		return err
	}

	for i, tag := range t[index].Tags {
		if tag == Tags {
			// Modify t[index].Tags directly instead of working on a copy
			t[index].Tags = append(t[index].Tags[:i], t[index].Tags[i+1:]...)
			tagMatched = true
			break // Add break here since the slice is now modified
		}
	}

	if tagMatched {
		return nil
	}
	return fmt.Errorf("no matching tags found")
}
