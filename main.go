package main

import (
	"log"
	"os"
	"strings"

	v1 "k8s.io/api/apps/v1"
	"k8s.io/cli-runtime/pkg/printers"
	"k8s.io/client-go/kubernetes/scheme"
)

const (
	ownerLabel = "owner"
)

func main() {
	files := os.Getenv("REPO_FILES")
	filesInRepo := strings.Split(files, " ")
	for _, file := range filesInRepo {
		filePath := "/github/workspace/" + file
		buf, err := os.ReadFile(filePath)
		log.Printf("Checking %s", filePath)
		if err != nil {
			log.Fatalf("error reading %s: %v", file, err)
			return
		}
		decode := scheme.Codecs.UniversalDeserializer().Decode
		obj, gvk, err := decode(buf, nil, nil)
		if gvk.Version == "" {
			log.Printf("Unmarshalled file: %s, but it's not a kubernetes manifest file", file)
		}
		fileYAML := obj.(*v1.Deployment)
		if err != nil {
			log.Fatalf("Error while decoding YAML object. Err was: %s", err)
		}
		log.Printf("k8s manifest labels = %v", fileYAML.ObjectMeta.Labels)
		if _, ok := fileYAML.ObjectMeta.Labels[ownerLabel]; !ok {
			log.Printf("adding 'owner' label")
			fileYAML.ObjectMeta.Labels[ownerLabel] = "platform"
			path, err := os.Create(filePath)
			y := printers.YAMLPrinter{}
			err = y.PrintObj(fileYAML, path)
			if err != nil {
				log.Fatalf("Failed to write file: %v", err)
			}
		}
	}
}
