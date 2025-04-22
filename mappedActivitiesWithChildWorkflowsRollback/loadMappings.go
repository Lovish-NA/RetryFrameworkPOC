package app

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

// ActivityMappings holds the mapping of activities to their respective fallback activities
type ActivityMappings struct {
	Mappings map[string]string `json:"mappings"`
}

// LoadMappings loads the activity mappings from a JSON file
func LoadMappings(filename string) (ActivityMappings, error) {
	// Print the current working directory for debugging
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get current working directory: %v", err)
	}
	fmt.Println("Current working directory:", cwd)

	fmt.Println("Loading mappings from " + filename)

	// Use an absolute path to ensure the file is found
	absPath, err := filepath.Abs(filename)
	fmt.Println("Loading mappings from absPath " + absPath)
	var mappings ActivityMappings

	data, err := os.ReadFile(absPath)
	fmt.Println("Data read from file:", string(data))
	if err != nil {
		fmt.Println("Error reading " + filename + ": " + err.Error())
		return mappings, err
	}
	err = json.Unmarshal(data, &mappings)
	fmt.Println("Mapping data:", mappings)
	if err != nil {
		fmt.Println("Error parsing " + filename + ": " + err.Error())
		return mappings, err
	}
	return mappings, nil
}
