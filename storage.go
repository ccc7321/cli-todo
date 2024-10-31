package main

import (
	"encoding/json"
	"os"
)

type FileStorage[T any] struct {
	Filename string
}

func NewStorage[T any](filename string) *FileStorage[T] {
	return &FileStorage[T]{Filename: filename}
}

func (s *FileStorage[T]) Save(data T) error {
	filedata, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return err
	}
	//permission, everyone can read and write to the file
	return os.WriteFile(s.Filename, filedata, 0644)
}

func (s *FileStorage[T]) Load(data *T) error {
	filedata, err := os.ReadFile(s.Filename)
	if err != nil {
		return err
	}
	return json.Unmarshal(filedata, data)
}
