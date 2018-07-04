package avest

import (
	"os"
	"github.com/peyman-abdi/avalanche/app/interfaces/services"
	"github.com/hjson/hjson-go"
	"strings"
	"path"
)

func CreateTemplateFiles(app services.Application, templates map[string]string) {
	err := os.MkdirAll(app.ResourcesPath("views/templates"), 0700)
	if err != nil {
		panic(err)
	}

	for filename, data := range templates {
		folders := strings.Split(filename, ".")
		if l := len(folders); l >= 3 {
			foldersPath := strings.Join(folders[:l-2], "/")
			filename = path.Join(foldersPath, folders[l-2] + "." + folders[l-1])
			os.MkdirAll(app.ConfigPath(foldersPath), 0700)
		}

		CreateFile(app.TemplatesPath(filename), []byte(data))
	}
}

func CreateConfigFiles(app services.Application, configs map[string]interface{}) {
	err := os.MkdirAll(app.ConfigPath(""), 0700)
	if err != nil {
		panic(err)
	}

	for filename, data := range configs {
		folders := strings.Split(filename, ".")
		if l := len(folders); l >= 3 {
			foldersPath := strings.Join(folders[:l-2], "/")
			filename = path.Join(foldersPath, folders[l-2] + "." + folders[l-1])
			os.MkdirAll(app.ConfigPath(foldersPath), 0700)
		}
		content, err := hjson.Marshal(data)
		if err != nil {
			panic(err)
		}

		CreateFile(app.ConfigPath(filename), content)
	}
}

func CreateFile(path string, content []byte) {
	file, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	file.Write(content)
	file.Close()
}

func StringRef(str string) *string {
	return &str
}
func IntRef(i int) *int {
	return &i
}
func StringRefOrNil(v interface{}) *string {
	if v == nil {
		return nil
	}

	return StringRef(v.(string))
}

func IntRefOrNil(v interface{}) *int {
	if v == nil {
		return nil
	}

	return IntRef(int(v.(float64)))
}

