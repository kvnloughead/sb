package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config holds the configuration for the sb CLI tool.
type Config struct {
	Repo  string            `yaml:"repo"`
	Slugs map[string]string `yaml:"slugs"`
}

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "Usage: go run gen_defaults.go <defaults.yaml> <output.go>\n")
		os.Exit(1)
	}
	input := os.Args[1]
	output := os.Args[2]
	data, err := os.ReadFile(input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading defaults.yaml: %v\n", err)
		os.Exit(1)
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing YAML: %v\n", err)
		os.Exit(1)
	}
	f, err := os.Create(output)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating output file: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()
	fmt.Fprintf(f, "package config\n\n// DefaultConfig is generated from defaults.yaml\nvar DefaultConfig = Config{\n\tRepo: \"%s\",\n\tSlugs: map[string]string{\n", cfg.Repo)
	for k, v := range cfg.Slugs {
		fmt.Fprintf(f, "\t\t\"%s\": \"%s\",\n", k, v)
	}
	fmt.Fprintf(f, "\t},\n}\n")
}
