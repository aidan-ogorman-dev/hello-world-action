package main

import (
	"log"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
	v1 "k8s.io/api/apps/v1"
)

const (
	ownerLabel = "owner"
)

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
		fileYAML := &v1.Deployment{}
		_ = yaml.Unmarshal(buf, fileYAML)
		if fileYAML.TypeMeta.APIVersion != "" {
			log.Printf("k8s manifest labels = %v", fileYAML.ObjectMeta.Labels)
			if _, ok := fileYAML.ObjectMeta.Labels[ownerLabel]; !ok {
				log.Printf("adding 'owner' label")
				fileYAML.ObjectMeta.Labels[ownerLabel] = "platform"
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
