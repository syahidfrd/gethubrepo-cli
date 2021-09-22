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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/spf13/viper"
)

var cfgFile string
var githubUsername string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gethubrepo",
	Short: "Get Github repository",
	Long:  "gethubrepo-cli is a tool that gives you github repository data from some user",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		repos, err := fetchGithubRepo(githubUsername)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		for _, v := range repos {
			fmt.Println("Name: " + v.Name)
			fmt.Println("Description: " + v.Description)
			fmt.Println("Repository: " + v.HtmlURL)
			fmt.Println("========================================")
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.gethubrepo-cli.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	rootCmd.PersistentFlags().StringVarP(&githubUsername, "username", "u", "", "github username")
	rootCmd.MarkPersistentFlagRequired("username")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".gethubrepo-cli" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".gethubrepo-cli")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}

type GithubRepo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	HtmlURL     string `json:"html_url"`
}

func fetchGithubRepo(username string) ([]*GithubRepo, error) {
	url := fmt.Sprintf("https://api.github.com/users/%s/repos", username)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("error request creation failed: %s", err.Error())
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("User-Agent", "gethubrepo-cli(github.com/syahidfrd/gethubrepo-cli)")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error when request via http client, cannot send request with error: %s", err.Error())
	}

	defer res.Body.Close()

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("cannot read response body with error: %s", err.Error())
	}

	if res.StatusCode >= 400 {
		return nil, fmt.Errorf("api is returning error. http status code: %s  api response: %s", strconv.Itoa(res.StatusCode), string(resBody))
	}

	var githubRepos []*GithubRepo
	if err := json.Unmarshal(resBody, &githubRepos); err != nil {
		return nil, fmt.Errorf("invalid body response, parse error during api request with message: %s", err.Error())
	}

	return githubRepos, nil

}
