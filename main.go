package main

import (
	"log"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

const (
	ownerLabel = "owner"
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
	filesInRepo := strings.Split(files, " ")
	log.Printf("List of files to check: %v\n", filesInRepo)
	for _, file := range filesInRepo {
		filePath := "/github/workspace/" + file
		buf, err := os.ReadFile(filePath)
		log.Printf("Checking %s", filePath)
		if err != nil {
			log.Fatalf("error reading %s: %v", file, err)
			return
		}
		fileYAML := &kubeYAMLFile{}
		_ = yaml.Unmarshal(buf, fileYAML)
		if fileYAML.ApiVersion != "" {
			log.Printf("k8s manifest labels = %v", fileYAML.Metadata.Labels)
			if _, ok := fileYAML.Metadata.Labels[ownerLabel]; !ok {
				log.Printf("adding 'owner' label")
				fileYAML.Metadata.Labels[ownerLabel] = "platform"
				buf, err = yaml.Marshal(fileYAML)
				if err != nil {
					log.Fatalf("Failed to marshal YAML: %v", err)
				}
				err := os.WriteFile(filePath, buf, 0644)
				if err != nil {
					log.Fatalf("Failed to write file: %v", err)
				}
			}
		} else {
			log.Printf("Unmarshalled file: %s, but it's not a kubernetes manifest file", file)
		}
	}
}
