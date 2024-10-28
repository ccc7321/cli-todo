package main

import (
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

type Todos []Todo

// the add function is a pointer receiver function as it is only accessible with Todos pointer
func (todos *Todos) add(title string) {
	todo := Todo{
		Title:       title,
		Completed:   false,
		CreatedAt:   time.Now(),
		CompletedAt: nil,
	}
	*todos = append(*todos, todo)
}

func (todos *Todos) validateIndex(index int) error {
	if index < 0 || index >= len(*todos) {
		err := errors.New("invalid index")
		fmt.Println(err)
		return err
	}
	return nil
}

// the delete function that receives a pointer - *Todos
func (todos *Todos) delete(index int) error {
	err := todos.validateIndex(index)
	if err != nil {
		return err
	}

	todosList := *todos

	*todos = append(todosList[:index], todosList[index+1:]...)

	return nil
}

func (todos *Todos) toggle(index int) error {
	t := *todos
	err := t.validateIndex(index)
	if err != nil {
		return err
	}

	isCompleted := t[index].Completed
	if !isCompleted {
		completionTime := time.Now()
		t[index].CompletedAt = &completionTime
	}

	t[index].Completed = !t[index].Completed

	return nil
}

func (todos *Todos) edit(index int, title string) error {
	t := *todos
	err := t.validateIndex(index)
	if err != nil {
		return err
	}

	t[index].Title = title

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
