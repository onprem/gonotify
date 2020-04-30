package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to GoNotify",
	Long:  "This subcommand lets you sign in into GoNotify",
	Run:   login,
}

var number string
var password string

func login(cmd *cobra.Command, args []string) {
	conf := getConfig()

	c := getClient(conf.BaseURL, conf.Token)

	token, err := c.Login(number, password)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	conf.Token = token
	conf.Phone = number
	conf.Password = password
	err = conf.Save()
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	fmt.Println("Login successful")
}

func init() {
	loginCmd.Flags().StringVarP(&number, "number", "n", "", "Primary phone number of your account")
	loginCmd.Flags().StringVarP(&password, "password", "p", "", "Password of your account")
	rootCmd.AddCommand(loginCmd)
}
