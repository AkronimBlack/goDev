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
	sender "github.com/AkronimBlack/dev-tools/pkg/amqpSender"
	"github.com/spf13/cobra"
)

// sendCmd represents the send command
var sendCmd = &cobra.Command{
	Use:   "send:amqp",
	Short: "Send a message to a defined topics",
	Long: `Read from a defined .json file, compose message and send it to event bus.
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
		"payload": {
			"email_subject": "{{faker.Sentence}}",
			"email_body": "{{faker.Paragraph}}",
			"test":{
				"some key":"{{faker.Name}}",
				"some key2":{
					"once more":"{{faker.Name}}"
				}
			}
		},
		"properties": {
			"message-name": "trigger-process",
			"sampleKey": "{{faker.UUID}}"
		}
	}
	`,
	Run: func(cmd *cobra.Command, args []string) {
		hostname, err := cmd.Flags().GetString("hostname")
		port, err := cmd.Flags().GetString("port")
		username, err := cmd.Flags().GetString("username")
		password, err := cmd.Flags().GetString("password")
		topics, err := cmd.Flags().GetStringSlice("topics")
		filename, err := cmd.Flags().GetString("message")
		num, err := cmd.Flags().GetInt("num")
		common.PanicOnError(err)
		sender.SendMessages(hostname, port, filename, username, password, topics, num)
	},
}

func init() {
	rootCmd.AddCommand(sendCmd)
	sendCmd.PersistentFlags().StringP("username", "u", "admin", "username")
	sendCmd.PersistentFlags().StringP("password", "p", "admin", "password")
	sendCmd.PersistentFlags().StringP("hostname", "z", "127.0.0.1", "host to connect to")
	sendCmd.PersistentFlags().StringP("port", "r", "5672", "port to connect to")
	sendCmd.PersistentFlags().StringSliceP("topics", "t", nil, "List of topics to send to")
	sendCmd.PersistentFlags().StringP("message", "f", "", "Use this message file to build and send message")
	sendCmd.PersistentFlags().IntP("num", "n", 1, "Number of massages to be sent")

	sendCmd.MarkFlagRequired("topics")
	sendCmd.MarkFlagRequired("message")
}
