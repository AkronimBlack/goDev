package wizard

import (
	"fmt"
	"log"
	"os"
	"strings"
	"text/template"

	"github.com/AkronimBlack/stock/common"
	"github.com/AkronimBlack/stock/pkg/templates"
)

const (
	TypeFile = "file"
	TypeDir  = "dir"
	//NamePlaceholder is a string the builder will look for and replace with name of full-name of the project
	NamePlaceholder = "{app_name}"

	ginFramework     = "gin"
	gorillaFramework = "gorilla"
	noFramework      = "none"
)

var templateMap = map[string]func() []byte{
	"main.go":            templates.GinTemplate,
	"main_test.go":       templates.MainTestTemplate,
	"Dockerfile":         templates.DockerfileTemplate,
	"Dockerfile.dev":     templates.DockerfileDevTemplate,
	"docker-compose.yml": templates.DockerComposeTemplate,
	"go.mod":             goModTemplate,
	".env":               templates.EnvTemplate,
	".env.example":       templates.EnvTemplate,
	".gitignore":         templates.GitIgnoreTemplate,
}

var executeOptions *Options

//HTTPFrameworks list of available framework templates
func HTTPFrameworks() []string {
	return []string{ginFramework, gorillaFramework, noFramework}
}

//DefaultHTTPFramework list of available framework templates
func DefaultHTTPFramework() string {
	return ginFramework
}

func Execute(opts *Options) {
	executeOptions = opts
	config := NewConfig(opts.Name, opts.FullName, opts.Maintainer, opts)
	common.LogJson(config)
	for _, x := range getObjectMap() {
		x.Build(config)
	}
}

func NewOptions(name, Maintainer, framework, fullName string) *Options {
	return &Options{
		Name:       name,
		Maintainer: Maintainer,
		Framework:  framework,
		FullName:   fullName,
	}
}

type Options struct {
	Name       string `json:"project_name"`
	Maintainer string `json:"maintainer"`
	Framework  string `json:"framework"`
	FullName   string `json:"full_name"`
}

func getTemplate(file string) func() []byte {
	var x func() []byte
	var ok bool
	if x, ok = templateMap[file]; !ok {
		log.Panicf("Template for file %s does not exist", file)
	}
	return x
}

func mainTemplate() []byte {
	if executeOptions.Framework != "" {

		switch executeOptions.Framework {
		case ginFramework:
			return templates.GinTemplate()
		}
	}
	return templates.MainTemplate()
}

func goModTemplate() []byte {
	base := []byte(`module {{.FullName}}

go 1.15
require (
	github.com/joho/godotenv v1.3.0
`)

	if executeOptions.Framework != "" {

		base = append(base, templates.GetDependency(executeOptions.Framework)...)
	}

	return append(base, []byte(`    
	)`)...)

}

func getObjectMap() []*Object {
	return []*Object{
		{
			Name: NamePlaceholder,
			Type: TypeDir,
			SubObjects: []*Object{
				{
					Name: "api",
					Type: TypeDir,
					SubObjects: []*Object{
						{
							Name: "openapi",
							Type: TypeDir,
						},
						{
							Name: "proto",
							Type: TypeDir,
						},
					},
				},
				{
					Name: "application",
					Type: TypeDir,
				},
				{
					Name: "cmd",
					Type: TypeDir,
					SubObjects: []*Object{
						{
							Name: NamePlaceholder,
							Type: TypeDir,
							SubObjects: []*Object{
								{
									Name:     "main.go",
									Type:     TypeFile,
									Template: getTemplate("main.go"),
								},
								{
									Name:     "main_test.go",
									Type:     TypeFile,
									Template: getTemplate("main_test.go"),
								},
							},
						},
					},
				},
				{
					Name: "docker",
					Type: TypeDir,
					SubObjects: []*Object{
						{
							Name:     "Dockerfile",
							Type:     TypeFile,
							Template: getTemplate("Dockerfile"),
						},
					},
				},
				{
					Name: "domain",
					Type: TypeDir,
				},
				{
					Name: "infrastructure",
					Type: TypeDir,
					SubObjects: []*Object{
						{
							Name: "transport",
							Type: TypeDir,
							SubObjects: []*Object{
								{
									Name: "http",
									Type: TypeDir,
								},
								{
									Name: "grpc",
									Type: TypeDir,
								},
								{
									Name: "amqp",
									Type: TypeDir,
								},
							},
						},
						{
							Name: "repositories",
							Type: TypeDir,
						},
					},
				},
				{
					Name: "logs",
					Type: TypeDir,
				},
				{
					Name:     "docker-compose.yml",
					Type:     TypeFile,
					Template: getTemplate("docker-compose.yml"),
				},
				{
					Name:     "Dockerfile",
					Type:     TypeFile,
					Template: getTemplate("Dockerfile.dev"),
				},
				{
					Name:     "go.mod",
					Type:     TypeFile,
					Template: getTemplate("go.mod"),
				},
				{
					Name:     ".env",
					Type:     TypeFile,
					Template: getTemplate(".env"),
				},
				{
					Name:     ".env.example",
					Type:     TypeFile,
					Template: getTemplate(".env"),
				},
				{
					Name:     ".gitignore",
					Type:     TypeFile,
					Template: getTemplate(".gitignore"),
				},
			},
		},
	}
}

//NewConfig config constructor
func NewConfig(name, fullName, maintainer string, templateData interface{}) *Config {
	return &Config{
		Name:         name,
		FullName:     fullName,
		Maintainer:   maintainer,
		TemplateData: templateData,
	}
}

//Config config struct for generators
type Config struct {
	Name         string      `json:"name"`
	FullName     string      `json:"full_name"`
	Maintainer   string      `json:"maintainer"`
	TemplateData interface{} `json:"template_data"`
}

type Builder interface {
	Build() error
}

type Object struct {
	Name       string
	Type       string
	SubObjects []*Object
	Template   func() []byte
}

//Build start building
func (o *Object) Build(config *Config) error {
	o.replacePlaceholder(config)
	log.Printf("Generating %s ", o.Name)
	var err error
	switch o.Type {
	case TypeFile:
		err = o.generateFile(config)
	case TypeDir:
		log.Println("Making directory", o.Name)
		err = os.MkdirAll(o.Name, os.ModePerm)
	}
	if err != nil {
		return err
	}
	if o.SubObjects == nil {
		return nil
	}
	for _, x := range o.SubObjects {
		x.Name = fmt.Sprintf("%s/%s", o.Name, x.Name)
		err = x.Build(config)
		if err != nil {
			return err
		}
	}
	return nil
}

func (o *Object) replacePlaceholder(config *Config) {
	var name string
	parts := strings.Split(o.Name, "/")
	if len(parts) > 1 {
		name = parts[len(parts)-1]
	} else {
		name = o.Name
	}
	if name == NamePlaceholder {
		var replacement string
		log.Println("name is placeholder. replacing with ", config.FullName)
		if config.Name != "" {
			replacement = config.Name
		} else {
			replacement = config.FullName
		}

		if len(parts) > 1 {
			parts[len(parts)-1] = replacement
			o.Name = strings.Join(parts, "/")
			return
		}
		o.Name = replacement
	}
}

func (o *Object) generateFile(config *Config) error {
	log.Printf("Creating: %s", o.Name)
	f, err := os.Create(o.Name)
	if err != nil {
		log.Panic(err.Error())
	}
	if o.Template != nil {
		mainTemplate := template.Must(template.New("main").Parse(string(o.Template())))
		err = mainTemplate.Execute(f, config.TemplateData)
	}
	if err != nil {
		log.Panic(err)
	}
	return f.Close()
}
