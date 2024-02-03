package main

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

type kubeMetadata struct {
	name        string            `yaml:"name"`
	namespace   string            `yaml:"namespace,omitempty"`
	labels      map[string]string `yaml:"labels,omitempty"`
	annotations map[string]string `yaml:"annotations,omitempty"`
}

type kubeYAMLFile struct {
	apiVersion string       `yaml:"apiVersion"`
	kind       string       `yaml:"kind"`
	metadata   kubeMetadata `yaml:"metadata"`
}

func main() {
	files := os.Getenv("REPO_FILES")
	filesInRepo := strings.Split(files, "\n")
	fmt.Printf("List of all files: %s", filesInRepo)
	for _, file := range filesInRepo {
		content, err := os.ReadFile(file)
		if err != nil {
			fmt.Printf("error reading %s: %v", file, err)
			return
		}
		fileYAML := &kubeYAMLFile{}
		_ = yaml.Unmarshal(content, fileYAML)
		if fileYAML.apiVersion != "" {
			fmt.Printf("kube labels = %v", fileYAML.metadata.labels)
		} else {
			fmt.Printf("Unmarshalled file: %s, but it's not a kubernetes manifest file", file)
		}
	}
}
