package main

import (
	"errors"
	"fmt"
	"time"
)

type Todo struct {
	Title     string
	Completed bool
	CreatedAt time.Time
	//CompletedAt is a pointer because it can be nil, i.e. when it is not completed
	CompletedAt *time.Time
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
