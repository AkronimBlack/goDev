/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"

	"github.com/AkronimBlack/stock/common"
	"github.com/AkronimBlack/stock/pkg/wizard"
	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
)

// wizCmd represents the wiz command
var wizCmd = &cobra.Command{
	Use:   "wiz",
	Short: "Scaffold generation wizard",
	Long:  `Add-on to the scaffold generator that will let you customize your generated project.`,
	Run: func(cmd *cobra.Command, args []string) {
		answers := struct {
			ProjectName   string   `json:"project_name"`
			HTTPFramework string   `survey:"http_framework" json:"http_framework"`
			Database      []string `survey:"database"`
		}{}

		// perform the questions
		err := survey.Ask(qs, &answers)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		common.LogJson(answers)
		nameData := common.ExtractNameData(common.SanitizeName(answers.ProjectName))
		wizard.Execute(wizard.NewOptions(nameData.ProjectName, nameData.Maintainer, answers.HTTPFramework, answers.ProjectName, answers.Database))
	},
}

// the questions to ask
var qs = []*survey.Question{
	{
		Name:     "projectName",
		Prompt:   &survey.Input{Message: "Your project name?"},
		Validate: survey.Required,
	},
	{
		Name: "http_framework",
		Prompt: &survey.Select{
			Message: "Choose a http framework:",
			Options: wizard.HTTPFrameworks(),
			Default: wizard.DefaultHTTPFramework(),
		},
	},
	{
		Name: "database",
		Prompt: &survey.MultiSelect{
			Message: `Choose a database adapater (you can always add more later):
NOTE: gorm is used as an ORM and coosing one or more adapters will pull gorm as well`,
			Options: wizard.DatabaseAdapters(),
			Default: wizard.DefaultHTTPFramework(),
		},
	},
}

func init() {
	rootCmd.AddCommand(wizCmd)
}
