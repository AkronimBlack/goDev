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
	listener "github.com/AkronimBlack/stock/pkg/amqpListener"
	"github.com/spf13/cobra"
)

// listenCmd represents the listen command
var listenCmd = &cobra.Command{
	Use:   "listen:amqp",
	Short: "Run listeners on defined topics",
	Long: `
	Listens to n topics and prints out messages.
	NOTE: Text payload represents STOMP messages and BINARY for amqp messages
	`,
	Run: func(cmd *cobra.Command, args []string) {
		hostname, err := cmd.Flags().GetString("hostname")
		port, err := cmd.Flags().GetString("port")
		username, err := cmd.Flags().GetString("username")
		password, err := cmd.Flags().GetString("password")
		topics, err := cmd.Flags().GetStringSlice("topics")
		common.PanicOnError(err)
		listener.LaunchListener(hostname, port, username, password, topics)
	},
}

func init() {
	rootCmd.AddCommand(listenCmd)
	listenCmd.PersistentFlags().StringP("username", "u", "admin", "username")
	listenCmd.PersistentFlags().StringP("password", "p", "admin", "password")
	listenCmd.PersistentFlags().StringP("hostname", "z", "127.0.0.1", "host to connect to")
	listenCmd.PersistentFlags().StringP("port", "r", "5672", "port to connect to")
	listenCmd.PersistentFlags().StringSliceP("topics", "t", nil, "List of topics to send to")
	listenCmd.PersistentFlags().StringP("message", "f", "", "Use this message file to build and send message")
	listenCmd.PersistentFlags().IntP("num", "n", 1, "Number of masseges to be sent")

	listenCmd.MarkFlagRequired("topics")
}
