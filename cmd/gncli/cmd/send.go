package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var sendCmd = &cobra.Command{
	Use:   "send",
	Short: "Send message to a group",
	Long:  "This subcommand sends a message to the given group",
	// Aliases: []string{"ls"},
	Run: send,
}

var group string

func send(cmd *cobra.Command, args []string) {
	data := strings.Join(args, " ")

	conf := getConfig()
	client := getClient(conf.BaseURL, conf.Token)

	err := client.Send(data, group)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	fmt.Println("Message sent successfully")
}

func init() {
	sendCmd.Flags().StringVarP(&group, "group", "g", "", "Name of group to send message")
	rootCmd.AddCommand(sendCmd)
}
