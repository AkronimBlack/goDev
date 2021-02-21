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
	//NamePlaceholder is a string the builder will look for and replace with name of full-name of the project
	NamePlaceholder = "{app_name}"
)

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
	nameData := common.ExtractNameData(common.SanitizeName(opts.Name))
	config := NewConfig(nameData.ProjectName, opts.Name, nameData.Maintainer, NewConfig(nameData.ProjectName, opts.Name, nameData.Maintainer, nil))
	for _, x := range getObjectMap() {
		x.Build(config)
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

func getObjectMap() []*Object {
	return []*Object{
		{
			Name: NamePlaceholder,
			Type: "dir",
			SubObjects: []*Object{
				{
					Name: "api",
					Type: "dir",
					SubObjects: []*Object{
						{
							Name: "openapi",
							Type: "dir",
						},
						{
							Name: "proto",
							Type: "dir",
						},
					},
				},
				{
					Name: "application",
					Type: "dir",
				},
				{
					Name: "cmd",
					Type: "dir",
					SubObjects: []*Object{
						{
							Name: NamePlaceholder,
							Type: "dir",
							SubObjects: []*Object{
								{
									Name:     "main.go",
									Type:     "file",
									Template: templates.MainTemplate,
								},
								{
									Name:     "main_test.go",
									Type:     "file",
									Template: templates.MainTestTemplate,
								},
							},
						},
					},
				},
				{
					Name: "docker",
					Type: "dir",
					SubObjects: []*Object{
						{
							Name:     "Dockerfile",
							Type:     "file",
							Template: templates.DockerfileTemplate,
						},
						{
							Name:     "Dockerfile.dev",
							Type:     "file",
							Template: templates.DockerfileDevTemplate,
						},
					},
				},
				{
					Name: "domain",
					Type: "dir",
				},
				{
					Name: "infrastructure",
					Type: "dir",
					SubObjects: []*Object{
						{
							Name: "transport",
							Type: "dir",
							SubObjects: []*Object{
								{
									Name: "http",
									Type: "dir",
								},
								{
									Name: "grpc",
									Type: "dir",
								},
								{
									Name: "amqp",
									Type: "dir",
								},
							},
						},
						{
							Name: "repositories",
							Type: "dir",
						},
					},
				},
				{
					Name:     "docker-compose.yml",
					Type:     "file",
					Template: templates.DockerComposeTemplate,
				},
				{
					Name:     "go.mod",
					Type:     "file",
					Template: templates.GoModTemplate,
				},
			},
		},
	}
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
	case "file":
		err = o.generateFile(config)
	case "dir":
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
