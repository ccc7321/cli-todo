package main

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

// In mock/mock_operator.go
type MockTodoOperator struct {
	// Test tracking fields
	AddWasCalled bool
	AddedTitles  []string
	ShouldError  bool
	DelWasCalled bool
	DeletedIndex int
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
	if m.ShouldError {
		return errors.New("Error")
	}
	m.AddWasCalled = true
	return nil
}
func (m *MockTodoOperator) ToggleTodo(id int) error {
	if m.ShouldError {
		return errors.New("Error")
	}
	m.AddWasCalled = true
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
