package main

import (
	"log"
	"os"
	"strings"

	v1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/cli-runtime/pkg/printers"
	"k8s.io/client-go/kubernetes/scheme"
)

const (
	ownerLabel   = "owner"
	ghVolumePath = "/github/workspace/"
)

func main() {
	noFilesAddedModified := false
	noFilesRenamed := false
	filesAddedModifiedSplit := []string{}
	filesRenamedSplit := []string{}

	filesAddedModified := os.Getenv("ADDED_MODIFIED_FILES")
	log.Printf("ADDED_MODIFIED_FILES = '%v'", filesAddedModified)
	if filesAddedModified == "" {
		noFilesAddedModified = true
	} else {
		filesAddedModifiedSplit = strings.Split(filesAddedModified, " ")
	}
	filesRenamed := os.Getenv("RENAMED_FILES")
	log.Printf("RENAMED_FILES = '%v'", filesRenamed)
	if filesRenamed == "" {
		noFilesRenamed = true
	} else {
		filesRenamedSplit = strings.Split(filesRenamed, " ")
	}
	if noFilesAddedModified && noFilesRenamed {
		log.Printf("Check complete, good process.")
		return
	}
	files := []string{}
	for _, f := range filesAddedModifiedSplit {
		files = append(files, f)
	}
	for _, f := range filesRenamedSplit {
		files = append(files, f)
	}
	for _, file := range files {
		filePath := ghVolumePath + file
		log.Printf("Checking %s", filePath)
		buf, err := os.ReadFile(filePath)
		if err != nil {
			log.Fatalf("error reading %s: %v", file, err)
			return
		}
		decode := scheme.Codecs.UniversalDeserializer().Decode
		obj, gvk, err := decode(buf, nil, nil)
		if err != nil {
			log.Fatalf("Error while decoding YAML object. Err was: %s", err)
		}
		switch gvk.Kind {
		case "":
			log.Printf("Unmarshalled file: %s, but it's not a kubernetes manifest file", file)
		case "Deployment":
			log.Printf("Checking deploy manifest")
			deploy := obj.(*v1.Deployment)
			deploy.ObjectMeta.Labels = checkLabels(deploy.ObjectMeta.Labels)
			if err := writeManifest(deploy, filePath); err != nil {
				log.Printf("error writing file: %v", err)
			}
		case "StatefulSet":
			log.Printf("Checking statefulset manifest")
			sts := obj.(*v1.StatefulSet)
			sts.ObjectMeta.Labels = checkLabels(sts.ObjectMeta.Labels)
			if err := writeManifest(sts, filePath); err != nil {
				log.Printf("error writing file: %v", err)
			}
		default:
			log.Printf("Unrecognised object type %s", gvk.Kind)
		}
	}
}

func checkLabels(labels map[string]string) map[string]string {
	if _, ok := labels[ownerLabel]; !ok {
		log.Printf("adding 'owner' label")
		labels[ownerLabel] = "platform"
	}
	return labels
}

func writeManifest(manifest runtime.Object, filePath string) error {
	path, err := os.Create(filePath)
	y := printers.YAMLPrinter{}
	err = y.PrintObj(manifest, path)
	if err != nil {
		log.Fatalf("Failed to write file: %v", err)
	}
	return err
}
