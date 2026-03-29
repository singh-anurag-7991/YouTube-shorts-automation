package script

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
)

type Script struct {
	ID       int    `json:"id"`
	Text     string `json:"text"`
	Category string `json:"category"`
	Used     bool   `json:"used"`
}

type Manager struct {
	FilePath string
	Scripts  []Script
}

func NewManager(filePath string) *Manager {
	return &Manager{
		FilePath: filePath,
	}
}

func (m *Manager) Load() error {
	file, err := os.Open(m.FilePath)
	if err != nil {
		return fmt.Errorf("could not open scripts file: %w", err)
	}
	defer file.Close()

	byteValue, err := ioutil.ReadAll(file)
	if err != nil {
		return fmt.Errorf("could not read scripts file: %w", err)
	}

	if err := json.Unmarshal(byteValue, &m.Scripts); err != nil {
		return fmt.Errorf("could not unmarshal scripts: %w", err)
	}

	return nil
}

func (m *Manager) GetNext() (*Script, error) {
	for i := range m.Scripts {
		if !m.Scripts[i].Used {
			return &m.Scripts[i], nil
		}
	}
	return nil, errors.New("no unused scripts found")
}

func (m *Manager) MarkAsUsed(id int) error {
	found := false
	for i := range m.Scripts {
		if m.Scripts[i].ID == id {
			m.Scripts[i].Used = true
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("script with ID %d not found", id)
	}

	// Save the updated scripts back to the file
	data, err := json.MarshalIndent(m.Scripts, "", "  ")
	if err != nil {
		return fmt.Errorf("could not marshal scripts: %w", err)
	}

	if err := ioutil.WriteFile(m.FilePath, data, 0644); err != nil {
		return fmt.Errorf("could not write scripts file: %w", err)
	}

	return nil
}
