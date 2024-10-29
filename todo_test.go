package main

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
	"time"
)

// In mock/mock_operator.go
type MockTodoOperator struct {
	// Test tracking fields
	AddWasCalled    bool
	AddedTitles     []string
	ShouldError     bool
	DelWasCalled    bool
	DeletedIndex    int
	EditWasCalled   bool
	EditedTitles    []string
	ToggleWasCalled bool
	mockTodos       Todos
}

var mockTodos = []Todo{
	{
		Title:     "buy milk",
		Completed: false,
		CreatedAt: time.Now(),
		Priority:  1,
		Tags:      []string{"shopping"},
	},
	{
		Title:       "buy video games",
		Completed:   true,
		CreatedAt:   time.Now(),
		CompletedAt: &time.Time{}, // or nil if not completed
		Priority:    2,
		Tags:        []string{"shopping", "entertainment"},
	},
}

// These method implementations make it satisfy the TodoOperator interface
func (m *MockTodoOperator) AddTodo(title string) error {
	m.AddWasCalled = true
	if m.ShouldError {
		return errors.New("Error")
	}

	m.AddedTitles = append(m.AddedTitles, title)
	return nil

}

func (m *MockTodoOperator) DeleteTodo(id int) error {
	m.DelWasCalled = true
	m.DeletedIndex = id
	if m.ShouldError {
		return errors.New("Error")
	}
	m.AddWasCalled = true
	return nil
}
func (m *MockTodoOperator) EditTodo(id int, title string) error {

	m.EditWasCalled = true
	if m.ShouldError {
		return errors.New("Error")
	}
	for i := 0; i < id+2; i++ {
		m.EditedTitles = append(m.EditedTitles, strconv.Itoa(i))
	}
	fmt.Println(m.EditedTitles)
	m.EditedTitles[id] = title
	fmt.Println(m.EditedTitles)
	return nil
}
func (m *MockTodoOperator) ToggleTodo(id int) error {
	m.ToggleWasCalled = true
	if m.ShouldError {
		return errors.New("Error")
	}
	m.mockTodos = mockTodos
	m.mockTodos.toggle(id)
	return nil
}

// If you prefer separate test functions:
func TestCommandFlags_Execute_Add(t *testing.T) {
	mock := &MockTodoOperator{}
	todos := &Todos{}

	flags := CommandFlags{
		Add:    "buy milk",
		Filter: -1,
		Del:    -1,
		Toggle: -1,
	}

	flags.Execute(mock, todos)

	assert.True(t, mock.AddWasCalled)
	assert.Equal(t, 1, len(mock.AddedTitles))
	assert.Equal(t, "buy milk", mock.AddedTitles[0])
}

func TestCommandFlags_Execute_Delete(t *testing.T) {
	mock := &MockTodoOperator{}
	todos := Todos{
		{Title: "buy milk", Completed: false},
		{Title: "buy video games", Completed: true},
	}

	flags := CommandFlags{
		Add:    "",
		Filter: -1,
		Del:    0,
		Toggle: -1,
	}

	flags.Execute(mock, &todos)

	assert.True(t, mock.DelWasCalled)
	assert.Equal(t, 0, mock.DeletedIndex)
}

func TestCommandFlags_Execute_Edit(t *testing.T) {
	mock := &MockTodoOperator{}
	todos := Todos{
		{Title: "buy milk", Completed: false},
		{Title: "buy video games", Completed: true},
	}

	flags := CommandFlags{
		Filter: -1,
		Del:    -1,
		Toggle: -1,
		Edit:   "0:buy chocolate milk",
	}
	fmt.Printf("%s\n", flags.Edit)
	flags.Execute(mock, &todos)

	assert.True(t, mock.EditWasCalled)
	assert.Equal(t, "buy chocolate milk", mock.EditedTitles[0])
}

func TestCommandFlags_Execute_Toggle(t *testing.T) {
	mock := &MockTodoOperator{}
	todos := Todos{
		{Title: "buy milk", Completed: false},
		{Title: "buy video games", Completed: true},
	}

	flags := CommandFlags{
		Filter: -1,
		Del:    -1,
		Toggle: 0,
	}
	fmt.Printf("%s\n", flags.Edit)
	flags.Execute(mock, &todos)

	assert.True(t, mock.ToggleWasCalled)
	assert.Equal(t, true, mock.mockTodos[0].Completed)
}
