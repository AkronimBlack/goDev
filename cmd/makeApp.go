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
	"github.com/AkronimBlack/stock/common"
	"github.com/AkronimBlack/stock/pkg/appgenerator"
	"github.com/spf13/cobra"
)

// makeAppCmd represents the make:app command
var makeAppCmd = &cobra.Command{
	Use:   "make:app",
	Short: "Create basic scaffolding for a new service",
	Long: `Will create a project with the follwing directory structure:
{app_name}	
	|-api
	|   |- openapi
	|	|- proto
	|-application
 	|-cmd
	|	|-{app_name}
	|		|- main.go
	|		|- main_test.go
	|-docker
	|    |- Dockerfile
	|	 |- Dockerfile.dev
	|-domain
	|-infrastructure
	|    |-transport
	|    |   |- http  
	|	 |   |- grpc
	|	 |   |- amqp
	|	 |-repositories  
	|-docker-compose.yml
	|-README.md
	|-.env
	|-.env.test
	`,
	Run: func(cmd *cobra.Command, args []string) {
		appName, err := cmd.Flags().GetString("app-name")
		common.PanicOnError(err)
		appgenerator.Execute(appgenerator.NewOptions(appName))
	},
}

func init() {

	makeAppCmd.Flags().StringP("app-name", "n", "", `The name of the app you are making. Take note that the convention is to use the github path.
Example : github.com\AkronimBlack\exampleProject
In this case the full name will be used to generate the go.mod file and the exampleProject will be use for dir and file generation as well as template data`)
	makeAppCmd.MarkFlagRequired("app-name")
	rootCmd.AddCommand(makeAppCmd)
}
