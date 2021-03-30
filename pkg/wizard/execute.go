package wizard

import (
	"log"

	"github.com/AkronimBlack/stock/common"
	"github.com/AkronimBlack/stock/pkg/appgenerator"
	"github.com/AkronimBlack/stock/pkg/templates"
)

const (
	ginFramework     = "gin"
	gorillaFramework = "gorilla"
	noFramework      = "none"

	mysqlAdapter    = "mysql"
	postgresAdapter = "postgres"
	sqliteAdapter   = "sqlite"
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

//DatabaseAdapters list of available database adapters
func DatabaseAdapters() []string {
	return []string{mysqlAdapter, postgresAdapter, sqliteAdapter}
}

//DefaultHTTPFramework list of available framework templates
func DefaultHTTPFramework() string {
	return ginFramework
}

func DefaultDatabaseAdapter() string {
	return mysqlAdapter
}

func Execute(opts *Options) {
	executeOptions = opts
	config := appgenerator.NewConfig(opts.Name, opts.FullName, opts.Maintainer, opts)
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

func getObjectMap() []*appgenerator.Object {
	return []*appgenerator.Object{
		{
			Name: appgenerator.NamePlaceholder,
			Type: appgenerator.TypeDir,
			SubObjects: []*appgenerator.Object{
				{
					Name: "api",
					Type: appgenerator.TypeDir,
					SubObjects: []*appgenerator.Object{
						{
							Name: "openapi",
							Type: appgenerator.TypeDir,
						},
						{
							Name: "proto",
							Type: appgenerator.TypeDir,
						},
					},
				},
				{
					Name: "application",
					Type: appgenerator.TypeDir,
				},
				{
					Name: "cmd",
					Type: appgenerator.TypeDir,
					SubObjects: []*appgenerator.Object{
						{
							Name: appgenerator.NamePlaceholder,
							Type: appgenerator.TypeDir,
							SubObjects: []*appgenerator.Object{
								{
									Name:     "main.go",
									Type:     appgenerator.TypeFile,
									Template: getTemplate("main.go"),
								},
								{
									Name:     "main_test.go",
									Type:     appgenerator.TypeFile,
									Template: getTemplate("main_test.go"),
								},
							},
						},
					},
				},
				{
					Name: "docker",
					Type: appgenerator.TypeDir,
					SubObjects: []*appgenerator.Object{
						{
							Name:     "Dockerfile",
							Type:     appgenerator.TypeFile,
							Template: getTemplate("Dockerfile"),
						},
					},
				},
				{
					Name: "domain",
					Type: appgenerator.TypeDir,
				},
				{
					Name: "infrastructure",
					Type: appgenerator.TypeDir,
					SubObjects: []*appgenerator.Object{
						{
							Name: "transport",
							Type: appgenerator.TypeDir,
							SubObjects: []*appgenerator.Object{
								{
									Name: "http",
									Type: appgenerator.TypeDir,
								},
								{
									Name: "grpc",
									Type: appgenerator.TypeDir,
								},
								{
									Name: "amqp",
									Type: appgenerator.TypeDir,
								},
							},
						},
						{
							Name: "repositories",
							Type: appgenerator.TypeDir,
						},
					},
				},
				{
					Name: "logs",
					Type: appgenerator.TypeDir,
				},
				{
					Name:     "docker-compose.yml",
					Type:     appgenerator.TypeFile,
					Template: getTemplate("docker-compose.yml"),
				},
				{
					Name:     "Dockerfile",
					Type:     appgenerator.TypeFile,
					Template: getTemplate("Dockerfile.dev"),
				},
				{
					Name:     "go.mod",
					Type:     appgenerator.TypeFile,
					Template: getTemplate("go.mod"),
				},
				{
					Name:     ".env",
					Type:     appgenerator.TypeFile,
					Template: getTemplate(".env"),
				},
				{
					Name:     ".env.example",
					Type:     appgenerator.TypeFile,
					Template: getTemplate(".env"),
				},
				{
					Name:     ".gitignore",
					Type:     appgenerator.TypeFile,
					Template: getTemplate(".gitignore"),
				},
			},
		},
	}
}
