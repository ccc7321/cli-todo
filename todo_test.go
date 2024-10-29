package main

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

// In mock/mock_operator.go
type MockTodoOperator struct {
	// Test tracking fields
	AddWasCalled bool
	AddedTitles  []string
	ShouldError  bool
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

func TestCommandFlags_Execute_Add(t *testing.T) {
	tests := []struct {
		name              string
		title             string
		shouldAddBeCalled bool
		wantError         bool
		wantAdded         string
	}{
		{
			name:              "successful add todo",
			title:             "buy milk",
			shouldAddBeCalled: true,
			wantError:         false,
			wantAdded:         "buy milk",
		},
		// ...
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &MockTodoOperator{
				AddWasCalled: false,
				AddedTitles:  []string{},
				ShouldError:  tt.wantError,
			}

			mockTodos := &Todos{}

			flags := &CommandFlags{
				Add:    tt.title,
				Filter: -1,
				Del:    -1,
				Toggle: -1,
			}

			flags.Execute(mock, mockTodos)
			fmt.Printf("Flags state before Execute: Add='%s', Filter=%d\n", flags.Add, flags.Filter) // Debug print

			// First just check if the method was called
			assert.Equal(t, tt.shouldAddBeCalled, mock.AddWasCalled, "AddTodo called status incorrect")

			// Then, ONLY if we don't expect an error and the method was called:
			if !tt.wantError && mock.AddWasCalled {
				assert.Equal(t, 1, len(mock.AddedTitles), "should have added exactly one todo")
				// Only check the added title if we actually have one
				if len(mock.AddedTitles) > 0 {
					assert.Equal(t, tt.wantAdded, mock.AddedTitles[0], "added todo doesn't match expected")
				}
			}
		})
	}
}
