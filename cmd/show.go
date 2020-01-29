/*
Copyright © 2020 Christopher J. Maahs <cmaahs@gmail.com>

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
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/google/go-github/github"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
)

var limit bool
var organization string

// GitRepository - A subset of the []*github.Repository fields that we care about.
type GitRepository struct {
	Name      string `json:"name"`
	URL       string `json:"url"`
	CloneSSH  string `json:"cloneSsh"`
	CloneHTTP string `json:"cloneHttp"`
	Private   bool   `json:"private"`
}

// showCmd represents the show command
var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Get a list of repositories for your Organization",
	Long: `Optionally specify a "search" term to match repository names
for example:

#> git-code show "spoon" --organization "myorg"`,
	Aliases: []string{"s", "list", "l"},
	Run: func(cmd *cobra.Command, args []string) {

		var nameFilter string

		if len(args) > 0 {
			nameFilter = args[0]
		} else {
			nameFilter = ""
		}

		getRepos, _ := getRepositories(nameFilter)

		buff, _ := json.MarshalIndent(&getRepos, "", "  ")
		fmt.Println(string(buff))

	},
}

func init() {

	rootCmd.AddCommand(showCmd)

}

func getRepositories(nameFilter string) ([]GitRepository, error) {

	organization = viper.GetString("organization")
	homeDir, _ := homedir.Dir()
	homeDirPath, _ := homedir.Expand(homeDir)
	data, _ := ioutil.ReadFile(homeDirPath + "/.gittoken")
	userToken := strings.TrimSpace(string(data))
	if len(userToken) == 0 {
		fmt.Printf("Token file not found at %v.\n", homeDirPath+"/.gittoken")
		os.Exit(1)
	}
	// Establish authentication context with the token.
	context := context.Background()
	tokenService := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: userToken},
	)
	tokenClient := oauth2.NewClient(context, tokenService)

	g := github.NewClient(tokenClient)

	// Check rate limits for authenticated user. If this fails,
	// we assume that your authentication has failed.
	_, _, err := g.RateLimits(context)
	if err != nil {
		fmt.Printf("Problem in getting rate limit information %v\n", err)
		os.Exit(1)
	}

	opt := &github.RepositoryListByOrgOptions{}

	var allRepos []GitRepository
	for {
		// iterate until we have all pages.
		repos, resp, err := g.Repositories.ListByOrg(context, organization, opt)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// loop through our page's returned repos, and add to our collection
		for _, repo := range repos {
			gitRepo := GitRepository{}
			gitRepo.Name = repo.GetName()
			gitRepo.URL = repo.GetHTMLURL()
			gitRepo.CloneSSH = repo.GetSSHURL()
			gitRepo.CloneHTTP = repo.GetCloneURL()
			gitRepo.Private = repo.GetPrivate()
			if len(nameFilter) > 0 {
				if strings.Contains(gitRepo.Name, nameFilter) {
					allRepos = append(allRepos, gitRepo)
				}
			} else {
				allRepos = append(allRepos, gitRepo)
			}
		}
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
	return allRepos, nil
}
