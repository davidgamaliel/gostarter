package filestructure

import (
	"archive/zip"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

func CreateFolderStructure(tempFolderPath, structurePath string) error {
	if structurePath != "" {
		ProcessYamlStruct(tempFolderPath, structurePath)
		return nil
	}
	folders := []string{
		"src",
		"pkg",
		"bin",
	}

	for _, folder := range folders {
		err := os.MkdirAll(tempFolderPath+"/"+folder, os.ModePerm)
		if err != nil {
			return err
		}
	}

	return nil
}

// var recursiveStructYaml func
func RecursiveStructYaml(basePath string, ss map[interface{}]interface{}) {
	for k, v := range ss {
		k := k.(string)
		switch v.(type) {
		case map[string]interface{}:
			RecursiveStructYaml(basePath+"/"+k, v.(map[interface{}]interface{}))
		case map[interface{}]interface{}:
			RecursiveStructYaml(basePath+"/"+k, v.(map[interface{}]interface{}))
		case string:
			createFolder(basePath + "/" + k)
			createFile(basePath + "/" + k + "/" + v.(string))
		default:
			createFolder(basePath + "/" + k)
		}
	}
}

func createFolder(path string) {
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		fmt.Println("failed create folder", err)
	}
}

func createFile(filepath string) {
	err := os.WriteFile(filepath, []byte{}, os.ModePerm)
	if err != nil {
		slog.Error("failed create file", slog.String("file", filepath), err)
	}
}

func ProcessYamlStruct(basePath, structurePath string) {
	var ss map[interface{}]interface{}
	yamlFile, err := os.ReadFile(structurePath)
	if err != nil {
		fmt.Println("failed read file yaml", err)
		return
	}

	err = yaml.Unmarshal(yamlFile, &ss)
	if err != nil {
		slog.Error("failed unmarshall yaml file", err)
		return
	}

	RecursiveStructYaml(basePath, ss)

	slog.Info("success read yaml file", slog.Any("result", ss))
}

func WrapToZip(tempFolderPath string) error {
	zipFile, err := os.Create(tempFolderPath + "/generated.zip")
	if err != nil {
		return err
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	err = filepath.Walk(tempFolderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// fmt.Println(">>> Path:", path)
		basePath := strings.TrimPrefix(tempFolderPath, "./")
		sanitizePath := strings.TrimPrefix(path, basePath+"/")
		// fmt.Println(">>> Sanitized Path:", sanitizePath)
		if strings.Contains(path, "generated.zip") {
			return nil
		}

		if sanitizePath == path {
			return nil
		}

		if info.IsDir() {
			// fmt.Println(">>> Creating folder entry:", path)
			_, err = zipWriter.Create(sanitizePath + "/")
			if err != nil {
				fmt.Println("Failed to create zip entry for folder:", err)
				return err
			}
			return nil
		}

		zipEntry, err := zipWriter.Create(sanitizePath)
		if err != nil {
			fmt.Println("Failed to create zip entry:", err)
			return err
		}

		fileContent, err := os.ReadFile(path)
		if err != nil {
			fmt.Println("Failed to read file content:", err)
			return err
		}

		_, err = zipEntry.Write(fileContent)
		if err != nil {
			fmt.Println("Failed to write file content to zip entry:", err)
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func CleanUp(tempFolderPath string) error {
	err := os.RemoveAll(tempFolderPath)
	if err != nil {
		return err
	}

	return nil
}
