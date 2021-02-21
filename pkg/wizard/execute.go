package wizard

import (
	"log"

	"github.com/AkronimBlack/dev-tools/common"
	"github.com/AkronimBlack/dev-tools/pkg/appgenerator"
	"github.com/AkronimBlack/dev-tools/pkg/templates"
)

//HTTPFrameforks list of available framework templates
func HTTPFrameforks() []string {
	return []string{"gin", "gorilla"}
}

//DefaultHTTPFramefork list of available framework templates
func DefaultHTTPFramefork() string {
	return "gin"
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

var templateMap = map[string]func() []byte{
	"main.go":            templates.MainTemplate,
	"main_test.go":       templates.MainTestTemplate,
	"Dockerfile":         templates.DockerfileTemplate,
	"Dockerfile.dev":     templates.DockerfileDevTemplate,
	"docker-compose.yml": templates.DockerComposeTemplate,
	"go.mod":             templates.GoModTemplate,
}

func getTemplate(file string) func() []byte {
	var x func() []byte
	var ok bool
	if x, ok = templateMap[file]; !ok {
		log.Panicf("Template for file %s does not exist", file)
	}
	return x
}

func Execute(opts *Options) {
	config := appgenerator.NewConfig(opts.Name, opts.FullName, opts.Maintainer, opts)
	common.LogJson(config)
	for _, x := range getObjectMap() {
		x.Build(config)
	}
}

func getObjectMap() []*appgenerator.Object {
	return []*appgenerator.Object{
		{
			Name: appgenerator.NamePlaceholder,
			Type: "dir",
			SubObjects: []*appgenerator.Object{
				{
					Name: "api",
					Type: "dir",
					SubObjects: []*appgenerator.Object{
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
					SubObjects: []*appgenerator.Object{
						{
							Name: appgenerator.NamePlaceholder,
							Type: "dir",
							SubObjects: []*appgenerator.Object{
								{
									Name:     "main.go",
									Type:     "file",
									Template: getTemplate("main.go"),
								},
								{
									Name:     "main_test.go",
									Type:     "file",
									Template: getTemplate("main_test.go"),
								},
							},
						},
					},
				},
				{
					Name: "docker",
					Type: "dir",
					SubObjects: []*appgenerator.Object{
						{
							Name:     "Dockerfile",
							Type:     "file",
							Template: getTemplate("Dockerfile"),
						},
						{
							Name:     "Dockerfile.dev",
							Type:     "file",
							Template: getTemplate("Dockerfile.dev"),
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
					SubObjects: []*appgenerator.Object{
						{
							Name: "transport",
							Type: "dir",
							SubObjects: []*appgenerator.Object{
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
					Template: getTemplate("docker-compose.yml"),
				},
				{
					Name:     "go.mod",
					Type:     "file",
					Template: getTemplate("go.mod"),
				},
			},
		},
	}
}
