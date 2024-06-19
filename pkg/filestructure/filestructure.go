package filestructure

import (
	"archive/zip"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func CreateFolderStructure(tempFolderPath string) error {
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

		fmt.Println(">>> Path:", path)
		basePath := strings.TrimPrefix(tempFolderPath, "./")
		sanitizePath := strings.TrimPrefix(path, basePath+"/")
		fmt.Println(">>> Sanitized Path:", sanitizePath)
		if strings.Contains(path, "generated.zip") {
			return nil
		}

		if sanitizePath == path {
			return nil
		}

		if info.IsDir() {
			fmt.Println(">>> Creating folder entry:", path)
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

	// err = os.Remove("generated.zip")
	// if err != nil {
	// 	return err
	// }

	return nil
}
