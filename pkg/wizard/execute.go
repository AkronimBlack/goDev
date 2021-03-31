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
	"gin":                templates.GinTemplate,
	"main_test.go":       templates.MainTestTemplate,
	"Dockerfile":         templates.DockerfileTemplate,
	"Dockerfile.dev":     templates.DockerfileDevTemplate,
	"docker-compose.yml": templates.DockerComposeTemplate,
	"go.mod":             goModTemplate,
	".env":               templates.EnvTemplate,
	".env.example":       templates.EnvTemplate,
	".gitignore":         templates.GitIgnoreTemplate,
	"migrate.go":         templates.MigrateTemplate,
	"connection.go":      templates.ConnectionTemplate,
	"bootOptions.go":     templates.BootOptionsTemplate,
	"constants.go":       templates.ConstantsTemplate,
	"helpers.go":         templates.HelpersTemplate,
}

var executeOptions *Options

//HTTPFrameworks list of available framework templates
func HTTPFrameworks() []string {
	return []string{ginFramework, gorillaFramework, noFramework}
}

//DatabaseAdapters list of available database adapters
func DatabaseAdapters() []string {
	return []string{mysqlAdapter}
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
	for _, x := range getObjectMap(opts) {
		x.Build(config)
	}
}

func NewOptions(name, Maintainer, framework, fullName string, databaseAdapters []string, logrus bool) *Options {
	return &Options{
		Name:             name,
		Maintainer:       Maintainer,
		Framework:        framework,
		FullName:         fullName,
		DatabaseAdapters: databaseAdapters,
		Logrus:           logrus,
	}
}

type Options struct {
	Name             string   `json:"project_name"`
	Maintainer       string   `json:"maintainer"`
	Framework        string   `json:"framework"`
	FullName         string   `json:"full_name"`
	DatabaseAdapters []string `json:"database_adapters"`
	Logrus           bool     `json:"logrus"`
}

func getTemplate(file string) func() []byte {
	var x func() []byte
	var ok bool
	if x, ok = templateMap[file]; !ok {
		log.Panicf("Template for file %s does not exist", file)
	}
	return x
}

func goModTemplate() []byte {
	base := []byte(`module {{.FullName}}

go 1.15
require (
	github.com/joho/godotenv v1.3.0
`)

	if executeOptions.Framework != "" {

		base = append(base, templates.GetDependency(executeOptions.Framework)...)
		base = append(base, []byte("\n")...)
	}

	if executeOptions.DatabaseAdapters != nil {

		for _, v := range executeOptions.DatabaseAdapters {
			base = append(base, templates.GetDependency(v)...)
		}
		base = append(base, []byte("\n")...)
	}

	if executeOptions.Logrus {
		base = append(base, templates.GetDependency("logrus")...)
		base = append(base, []byte("\n")...)
	}

	return append(base, []byte(`    
	)`)...)

}

func getObjectMap(opts *Options) []*appgenerator.Object {
	return []*appgenerator.Object{
		{
			Name:    appgenerator.NamePlaceholder,
			Type:    appgenerator.TypeDir,
			Renders: true,
			SubObjects: []*appgenerator.Object{
				{
					Name:    "api",
					Type:    appgenerator.TypeDir,
					Renders: true,
					SubObjects: []*appgenerator.Object{
						{
							Name:    "openapi",
							Type:    appgenerator.TypeDir,
							Renders: true,
						},
						{
							Name:    "proto",
							Type:    appgenerator.TypeDir,
							Renders: false,
						},
					},
				},
				{
					Name:    "application",
					Type:    appgenerator.TypeDir,
					Renders: true,
				},
				{
					Name:    "cmd",
					Type:    appgenerator.TypeDir,
					Renders: true,
					SubObjects: []*appgenerator.Object{
						{
							Name:    appgenerator.NamePlaceholder,
							Type:    appgenerator.TypeDir,
							Renders: true,
							SubObjects: []*appgenerator.Object{
								{
									Name:     "main.go",
									Type:     appgenerator.TypeFile,
									Template: getTemplate("main.go"),
									Renders:  true,
									Evaluate: func(object *appgenerator.Object) {
										if opts.Framework != "" {
											object.Template = getTemplate(opts.Framework)
											object.Renders = true
										}
									},
								},
								{
									Name:     "main_test.go",
									Type:     appgenerator.TypeFile,
									Renders:  true,
									Template: getTemplate("main_test.go"),
								},
							},
						},
					},
				},
				{
					Name:    "docker",
					Type:    appgenerator.TypeDir,
					Renders: true,
					SubObjects: []*appgenerator.Object{
						{
							Name:     "Dockerfile",
							Type:     appgenerator.TypeFile,
							Renders:  true,
							Template: getTemplate("Dockerfile"),
						},
					},
				},
				{
					Name:    "domain",
					Renders: true,
					Type:    appgenerator.TypeDir,
				},
				{
					Name:    "infrastructure",
					Type:    appgenerator.TypeDir,
					Renders: true,
					SubObjects: []*appgenerator.Object{
						{
							Name:    "transport",
							Type:    appgenerator.TypeDir,
							Renders: true,
							SubObjects: []*appgenerator.Object{
								{
									Name:    "http",
									Renders: true,
									Type:    appgenerator.TypeDir,
								},
								{
									Name:    "grpc",
									Renders: true,
									Type:    appgenerator.TypeDir,
								},
								{
									Name:    "amqp",
									Renders: true,
									Type:    appgenerator.TypeDir,
								},
							},
						},
						{
							Name:    "repositories",
							Renders: true,
							Type:    appgenerator.TypeDir,
							SubObjects: []*appgenerator.Object{
								{
									Name:     "connection.go",
									Type:     appgenerator.TypeFile,
									Renders:  false,
									Template: getTemplate("connection.go"),
									Evaluate: func(object *appgenerator.Object) {
										if opts.DatabaseAdapters != nil {
											object.Renders = true
										}
									},
								},
								{
									Name:     "migrate.go",
									Type:     appgenerator.TypeFile,
									Renders:  true,
									Template: getTemplate("migrate.go"),
									Evaluate: func(object *appgenerator.Object) {
										if opts.DatabaseAdapters != nil {
											object.Renders = true
										}
									},
								},
							},
						},
					},
				},
				{
					Name:    "shared",
					Type:    appgenerator.TypeDir,
					Renders: true,
					SubObjects: []*appgenerator.Object{
						{
							Name:     "bootOptions.go",
							Type:     appgenerator.TypeFile,
							Renders:  true,
							Template: getTemplate("bootOptions.go"),
						},
						{
							Name:     "constants.go",
							Type:     appgenerator.TypeFile,
							Renders:  true,
							Template: getTemplate("constants.go"),
						},
						{
							Name:     "helpers.go",
							Type:     appgenerator.TypeFile,
							Renders:  true,
							Template: getTemplate("helpers.go"),
						},
					},
				},
				{
					Name:    "logs",
					Renders: true,
					Type:    appgenerator.TypeDir,
				},
				{
					Name:     "docker-compose.yml",
					Type:     appgenerator.TypeFile,
					Renders:  true,
					Template: getTemplate("docker-compose.yml"),
				},
				{
					Name:     "Dockerfile",
					Type:     appgenerator.TypeFile,
					Renders:  true,
					Template: getTemplate("Dockerfile.dev"),
				},
				{
					Name:     "go.mod",
					Renders:  true,
					Type:     appgenerator.TypeFile,
					Template: getTemplate("go.mod"),
				},
				{
					Name:     ".env",
					Renders:  true,
					Type:     appgenerator.TypeFile,
					Template: getTemplate(".env"),
				},
				{
					Name:     ".env.example",
					Renders:  true,
					Type:     appgenerator.TypeFile,
					Template: getTemplate(".env"),
				},
				{
					Name:     ".gitignore",
					Renders:  true,
					Type:     appgenerator.TypeFile,
					Template: getTemplate(".gitignore"),
				},
			},
		},
	}
}
