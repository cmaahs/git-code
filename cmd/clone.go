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
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
)

var directory string

// cloneCmd represents the clone command
var cloneCmd = &cobra.Command{
	Use:   "clone",
	Short: "Clone a repository",
	Long: `Specify a partial repository name, multiple matches will output
all of the matches, only when a single match is made that the clone will happen
example:

#> git-code clone "spoon"

May output:
Multiple Repository match, please be more specific
	miro-windows-spoon
	global-mute-spoon

#> git-code clone "mute-spoon"

Will match on the latter and clone it into directory: global-mute-spoon

#> git-code clone "mute-spoon" --directory "JIRA-2150"

Will match on the latter and clone it into a directory named: JIRA-2150`,
	Aliases: []string{"c", "cl"},
	Run: func(cmd *cobra.Command, args []string) {

		var repoName string
		directory, _ = cmd.Flags().GetString("directory")
		if len(args) > 0 {
			repoName = args[0]
		} else {
			fmt.Println("Must specify a repository name as an argument.")
			fmt.Println("#> git-code clone 'reponame'")
			os.Exit(1)
		}
		homeDir, _ := homedir.Dir()
		homeDirPath, _ := homedir.Expand(homeDir)
		data, _ := ioutil.ReadFile(homeDirPath + "/.gittoken")
		userToken := strings.TrimSpace(string(data))
		if len(userToken) == 0 {
			// this implies an error so we need to bail.
			fmt.Printf("Token file not found at %v.\n", homeDirPath+"/.gittoken")
			os.Exit(1)
		}
		getRepos, _ := getRepositories(repoName)
		if len(getRepos) > 1 {
			fmt.Println("Multiple Repository match, please be more specific:")
			for _, repo := range getRepos {
				fmt.Println(fmt.Sprintf("\t%s", repo.Name))
			}
			os.Exit(1)
		}
		if len(directory) == 0 {
			directory = fmt.Sprintf("./%s", getRepos[0].Name)
		}
		fmt.Println(fmt.Sprintf("Cloning into %s", directory))
		url := getRepos[0].CloneHTTP
		_, err := git.PlainClone(directory, false, &git.CloneOptions{
			Auth: &http.BasicAuth{
				Username: "gittoken", // yes, this can be anything except an empty string
				Password: userToken,
			},
			URL:               url,
			Progress:          os.Stdout,
			RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
		})
		if err != nil {
			fmt.Printf("Failed to clone the repository")
			os.Exit(1)
		}
	},
}

func init() {

	rootCmd.AddCommand(cloneCmd)
	cloneCmd.Flags().String("directory", "", "Target directory to clone into, existing directory must be empty.")

}
