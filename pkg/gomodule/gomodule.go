package gomodule

import (
	"os"
)

func GenerateGoMod(tempFolderPath string) error {
	content := []byte("module github.com/bitzero/teststarter")
	content = append(content, []byte("\n")...)
	content = append(content, []byte("\n")...)
	content = append(content, []byte("go 1.22")...)
	content = append(content, []byte("\n")...)
	content = append(content, []byte("\n")...)
	content = append(content, []byte("require (")...)
	content = append(content, []byte("\n")...)
	content = append(content, []byte("\tgithub.com/go-sql-driver/mysql")...)
	content = append(content, []byte("\n")...)
	content = append(content, []byte(")")...)

	err := os.WriteFile(tempFolderPath+"/go.mod", content, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

// module github.com/bitzero/teststarter

// 		go 1.20

// 		require (
// 			github.com/go-sql-driver/mysql
// 		)
