package wizard

import "github.com/AkronimBlack/dev-tools/pkg/appgenerator"

//AvailableFrameforks list of available framework templates
func AvailableFrameforks() []string {
	return []string{"gin", "gorilla"}
}

func NewOptions(projectName, maintainerName, framework string) *Options {
	return &Options{
		ProjectName:    projectName,
		MaintainerName: maintainerName,
		Framework:      framework,
	}
}

type Options struct {
	ProjectName    string `json:"project_name"`
	MaintainerName string `json:"maintainer"`
	Framework      string `json:"framework"`
}

func Execute(opts *Options) {
	appgenerator.Execute(appgenerator.NewOptions(opts.ProjectName))
}
