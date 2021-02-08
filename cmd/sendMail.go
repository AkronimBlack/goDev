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
	"github.com/AkronimBlack/dev-tools/common"
	sender "github.com/AkronimBlack/dev-tools/pkg/mailSender"
	"github.com/spf13/cobra"
)

// sendMailCmd represents the sendMail command
var sendMailCmd = &cobra.Command{
	Use:   "send:mail",
	Short: "Send an email",
	Long: `Service to send an x number of email to y number of contacts z number of times. Offers faking data as well.
	Read from a defined .json file, compose message and send it to event bus.
	Available faker options for making messages. 	

	Email           string
	PhoneNumber     string
	URL             string
	UserName        string
	TitleMale       string
	TitleFemale     string
	FirstName       string
	FirstNameMale   string
	FirstNameFemale string
	LastName        string
	Name            string
	Date            string
	Time            string
	MonthName       string
	Year            string
	DayOfWeek       string
	DayOfMonth      string
	Timestamp       string
	Century         string
	TimeZone        string
	TimePeriod      string
	Word            string
	Sentence        string
	Paragraph       string
	Currency        string
	UUID            string

	How to use:
	Payload section should be in map[string]interface{} format and you can use faker wherever as log as it is a value and not a key
	Example:
	{
		"subject": "{{faker.Sentence}}",
		"body":"{{faker.Paragraph}}",
		"to":["someEmail@email.com","someEmail2@email.com"]
	}
	`,
	Run: func(cmd *cobra.Command, args []string) {
		message, err := cmd.Flags().GetString("message")
		num, err := cmd.Flags().GetInt("num")
		from, err := cmd.Flags().GetString("from")
		auth, err := cmd.Flags().GetBool("auth")
		concurrent, err := cmd.Flags().GetBool("concurrent")
		hostname, err := cmd.Flags().GetString("hostname")
		username, err := cmd.Flags().GetString("hostname")
		password, err := cmd.Flags().GetString("password")
		common.PanicOnError(err)
		sender.Execute(message, from, username, password, hostname, auth, concurrent, num)
	},
}

func init() {
	rootCmd.AddCommand(sendMailCmd)
	sendMailCmd.Flags().StringP("from", "q", "", "Set mail sender")
	sendMailCmd.Flags().BoolP("auth", "a", false, "Auth required")
	sendMailCmd.Flags().BoolP("concurrent", "x", true, "Run in concurrent mod")

	sendMailCmd.Flags().StringP("username", "u", "admin", "username")
	sendMailCmd.Flags().StringP("password", "p", "admin", "password")

	sendMailCmd.Flags().StringP("hostname", "z", "127.0.0.1", "host to connect to")

	sendMailCmd.Flags().StringP("message", "f", "", "Use this config.json file to build and send message")
	sendMailCmd.Flags().IntP("num", "n", 1, "Number of massages to be sent")
	sendMailCmd.MarkFlagRequired("from")
	sendMailCmd.MarkFlagRequired("hostname")
	sendMailCmd.MarkFlagRequired("message")
}
