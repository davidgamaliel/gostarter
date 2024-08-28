package program

import (
	"log/slog"
	"os"
	"time"

	"github.com/bitzero/gostarter/internal/database"
	"github.com/bitzero/gostarter/pkg/filestructure"
	"github.com/bitzero/gostarter/pkg/gomodule"
	"gopkg.in/yaml.v2"
)

type Project struct {
	Name        string
	Basepath    string
	DBDriverMap map[string]database.DBDriver
	DirStruct   []byte
	DBDiver     string
}

func CreateProject(name, path, dirStruct string) (Project, error) {
	dirFile, err := os.ReadFile(dirStruct)
	if err != nil {
		slog.Error("Failed to read directory structure file", err)
		return Project{}, err
	}
	return Project{
		Name:        name,
		Basepath:    path,
		DBDriverMap: database.GetDBDriverMap(),
		DirStruct:   dirFile,
	}, nil
}

func (p *Project) SetDBDriver(driver string) {
	db := database.Driver(driver)
	p.DBDiver = db.String()
}

func (p *Project) CreateMainFile() error {
	err := p.GeneraterDirStructure()
	if err != nil {
		return err
	}

	err = gomodule.GenerateGoMod(p.Basepath, p.Name)
	if err != nil {
		return err
	}

	err = gomodule.GoGetPackage(p.Basepath, p.DBDriverMap[p.DBDiver].PackageName)
	if err != nil {
		return err
	}

	err = filestructure.WrapToZip(p.Basepath)
	if err != nil {
		return err
	}

	go func() {
		time.Sleep(10 * time.Second)
		filestructure.CleanUp(p.Basepath)
	}()
	return nil
}

func (p *Project) GeneraterDirStructure() error {
	var maps map[interface{}]interface{}
	err := yaml.Unmarshal(p.DirStruct, &maps)
	if err != nil {
		slog.Error("Failed to unmarshal yaml file", err)
		return err
	}

	filestructure.RecursiveStructYaml(p.Basepath, maps)
	return nil
}
