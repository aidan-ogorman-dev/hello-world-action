package main

import (
	"log"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

type kubeMetadata struct {
	Name        string            `yaml:"name"`
	Namespace   string            `yaml:"namespace,omitempty"`
	Labels      map[string]string `yaml:"labels,omitempty"`
	Annotations map[string]string `yaml:"annotations,omitempty"`
}

type kubeYAMLFile struct {
	ApiVersion string       `yaml:"apiVersion"`
	Kind       string       `yaml:"kind"`
	Metadata   kubeMetadata `yaml:"metadata"`
}

func main() {
	files := os.Getenv("REPO_FILES")
	filesInRepo := strings.Split(files, "\n")
	log.Printf("List of all files: %s", filesInRepo)
	for _, file := range filesInRepo {
		content, err := os.ReadFile("/github/workspace/" + file)
		if err != nil {
			log.Fatalf("error reading %s: %v", file, err)
			return
		}
		fileYAML := &kubeYAMLFile{}
		_ = yaml.Unmarshal(content, fileYAML)
		log.Printf("file yaml contents = %#v\n", fileYAML)
		if fileYAML.ApiVersion != "" {
			log.Printf("k8s manifest labels = %v", fileYAML.Metadata.Labels)
			// check for missing label, add it, commit back to GH
		} else {
			log.Printf("Unmarshalled file: %s, but it's not a kubernetes manifest file", file)
		}
	}
}
