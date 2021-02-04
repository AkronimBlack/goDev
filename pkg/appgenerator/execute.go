package appgenerator

import (
	"fmt"
	"html/template"
	"log"
	"os"
	"strings"

	"github.com/AkronimBlack/dev-tools/common"
	"github.com/AkronimBlack/dev-tools/pkg/templates"
)

const (
	apiDr             = "api"
	applicationDir    = "application"
	cmdDir            = "cmd"
	dockerDir         = "docker"
	domainDir         = "domain"
	infrastructureDir = "infrastructure"

	generatorVersion = 0.1
)

func getMap() map[string]interface{} {
	return map[string]interface{}{
		"api":         []string{"openapi", "proto"},
		"application": nil,
		"cmd": map[string]interface{}{
			namePlaceholder: []string{"main.go", "main_test.go"},
		},
		"docker": []file{
			{name: "Dockerfile"}, {name: "Dockerfile.dev"},
		},
		"domain": nil,
		"infrastructure": map[string]interface{}{
			"transport":    []string{"http", "grpc", "amqp"},
			"repositories": nil,
		},
		"docker-compose.yml": nil,
		"go.mod":             nil,
		"go.sum":             nil,
	}
}

//NewOptions opts constructor
func NewOptions(name string) *Options {
	return &Options{
		Name: name,
	}
}

//Options opts to use when running execute func for creating scaffolding
type Options struct {
	Name string `json:"name"`
}

/*Execute run the app scaffolding builder */
func Execute(opts *Options) {
	pwd, err := os.Getwd()
	if err != nil {
		log.Println(err.Error())
		return
	}

	log.Println("Running appgenerator package")
	log.Printf("VERSION: %f", generatorVersion)
	log.Printf("Directory: %s", pwd)

	fullName := opts.Name
	var maintainer string
	parts := strings.Split(opts.Name, "/")
	if len(parts) > 1 {
		opts.Name = parts[len(parts)-1]
	}
	if len(parts) == 3 {
		maintainer = parts[1]
	}

	if strings.Contains(opts.Name, " ") {
		log.Println("Whitespace detected. Cleaning up app name ...")
		components := strings.Split(opts.Name, " ")
		for i, v := range components {
			if i == 0 {
				v = strings.ToLower(v)
				continue
			}
			v = strings.Title(strings.ToLower(v))
		}
		opts.Name = strings.Join(components, "")
		log.Printf("App name: %s", opts.Name)
	}
	log.Println("Creating directory structure....")
	rootDir := fmt.Sprintf("%s/%s", pwd, opts.Name)
	config := NewConfig(opts.Name, rootDir, pwd, fullName, maintainer)

	iterateAndBuild(*config, getMap())
}

type file struct {
	name string
}

var namePlaceholder = "{app_name}"

//NewConfig config constructor
func NewConfig(name, rootDir, pwd, fullName, maintainer string) *Config {
	return &Config{
		Name:       name,
		RootDir:    rootDir,
		Pwd:        pwd,
		FullName:   fullName,
		Maintainer: maintainer,
	}
}

//Config config struct for generators
type Config struct {
	Name       string `json:"name"`
	RootDir    string `json:"root_dir"`
	Pwd        string `json:"call_location"`
	FullName   string `json:"full_name"`
	Maintainer string `json:"maintainer"`
}

func iterateAndBuild(opts Config, dirMap interface{}) error {
	var err error
	if dirMap == nil {
		return nil
	}
	switch dirMap.(type) {
	case string:
		if dirMap.(string) == namePlaceholder {
			dirMap = opts.Name
		}
		common.PanicOnError(generate(opts.RootDir, dirMap.(string), "", opts))
		break
	case []string:
		for _, v := range dirMap.([]string) {
			if v == namePlaceholder {
				v = opts.Name
			}
			common.PanicOnError(generate(opts.RootDir, v, "", opts))
		}
		break
	case map[string]interface{}:
		for i, v := range dirMap.(map[string]interface{}) {
			if i == namePlaceholder {
				i = opts.Name
			}
			common.PanicOnError(generate(opts.RootDir, i, "", opts))
			newOpts := opts
			newOpts.RootDir = fmt.Sprintf("%s/%s", opts.RootDir, i)
			iterateAndBuild(newOpts, v)
		}
		break
	case []file:
		for _, v := range dirMap.([]file) {
			if v.name == namePlaceholder {
				v.name = opts.Name
			}
			common.PanicOnError(generate(opts.RootDir, v.name, "file", opts))
		}
	case file:
		dm := dirMap.(file)
		if dm.name == namePlaceholder {
			dm.name = opts.Name
		}
		common.PanicOnError(generate(opts.RootDir, dm.name, "file", opts))
		break
	default:
		log.Println("Something is fishy here", dirMap)
	}
	return err
}

func generate(location, value, valueType string, opts Config) error {
	log.Printf("Generating %s %s", location, value)
	switch valueType {
	case "file":
		parts := strings.Split(value, ".")
		if len(parts) > 1 {
			return generateFile(location, value, opts)
		}
		return generateFile(location, value, opts)
	case "dir":
		return os.MkdirAll(fmt.Sprintf("%s/%s", location, value), os.ModePerm)
	}

	parts := strings.Split(value, ".")
	if len(parts) > 1 {
		return generateFile(location, value, opts)
	}

	return os.MkdirAll(fmt.Sprintf("%s/%s", location, value), os.ModePerm)
}

func generateFile(location, name string, opts Config) error {
	filename := name
	mainFile := fmt.Sprintf("%s/%s", location, filename)
	templateFunc := getTemplate(name)
	log.Printf("Creating: %s", mainFile)
	f, err := os.Create(mainFile)
	if err != nil {
		log.Panic(err.Error())
	}
	if templateFunc != nil {
		mainTemplate := template.Must(template.New("main").Parse(string(templateFunc())))
		err = mainTemplate.Execute(f, opts)
	}
	if err != nil {
		log.Panic(err)
	}
	return f.Close()
}

var templatesMap = map[string]func() []byte{
	"main.go":            templates.MainTemplate,
	"main_test.go":       templates.MainTestTemplate,
	"docker-compose.yml": templates.DockerComposeTemplate,
	"go.mod":             templates.GoModTemplate,
	"Dockerfile":         templates.DockerfileTemplate,
	"Dockerfile.dev":     templates.DockerfileDevTemplate,
}

func getTemplate(name string) func() []byte {
	if f, ok := templatesMap[name]; ok {
		return f
	}
	return nil
}
