package main

import (
	"context"
	"fmt"
	"log"
  "sync"

	"github.com/google/go-github/v50/github"
	"github.com/spf13/cobra"
)

var user string
var threads int

func main() {
	rootCmd := &cobra.Command{
		Use:   "github",
		Short: "A command-line tool to retrieve GitHub user and repository details",
	}

	userCmd := &cobra.Command{
		Use:   "user",
		Short: "Retrieve user details",
		RunE:  getUser,
	}

	userCmd.Flags().StringVarP(&user, "user", "u", "", "GitHub username")
	userCmd.Flags().IntVarP(&threads, "threads", "t", 1, "Number of threads to use")

	repoCmd := &cobra.Command{
		Use:   "repo",
		Short: "Retrieve repository details",
		RunE:  getRepos,
	}

	repoCmd.Flags().StringVarP(&user, "user", "u", "", "GitHub username")
	repoCmd.Flags().IntVarP(&threads, "threads", "t", 1, "Number of threads to use")

	orgCmd := &cobra.Command{
		Use:   "org",
		Short: "Retrieve organization details",
		RunE:  getOrgs,
	}

	orgCmd.Flags().StringVarP(&user, "user", "u", "", "GitHub username")
	orgCmd.Flags().IntVarP(&threads, "threads", "t", 1, "Number of threads to use")


	rootCmd.AddCommand(userCmd, repoCmd, orgCmd)

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func getClient() *github.Client {
	return github.NewClient(nil)
}

func getUser(cmd *cobra.Command, args []string) error {
	if user == "" {
		return fmt.Errorf("user flag is required")
	}

	client := getClient()

	u, _, err := client.Users.Get(context.Background(), user)
	if err != nil {
		return err
	}

	fmt.Println("Name:", u.GetName())
	fmt.Println("Email:", u.GetEmail())
	fmt.Println("Company:", u.GetCompany())
	fmt.Println("Location:", u.GetLocation())
	fmt.Println("Bio:", u.GetBio())
  fmt.Println("Followers:", u.GetFollowers())
  fmt.Println("Following:", u.GetFollowing())
  fmt.Println("Created:", u.GetCreatedAt())
  fmt.Println("Updated:", u.GetUpdatedAt())
	return nil
}

func getRepos(cmd *cobra.Command, args []string) error {
	if user == "" {
		return fmt.Errorf("user flag is required")
	}

	client := getClient()

	opt := &github.RepositoryListOptions{Type: "owner", Sort: "updated", Direction: "desc"}
	repos, _, err := client.Repositories.List(context.Background(), user, opt)
	if err != nil {
		return err
	}

	wg := &sync.WaitGroup{}
	for i := 0; i < threads; i++ {
		wg.Add(1)
		go func(start, end int) {
			defer wg.Done()
			for j := start; j < end; j++ {
				r := repos[j]
				fmt.Println("Name:", r.GetName())
				fmt.Println("Description:", r.GetDescription())
				fmt.Println("Language:", r.GetLanguage())
				fmt.Println("Stars:", r.GetStargazersCount())
				fmt.Println("Forks:", r.GetForksCount())
				fmt.Println("Updated:", r.GetUpdatedAt())
				fmt.Println()
			}
		}(i*len(repos)/threads, (i+1)*len(repos)/threads)
	}
	wg.Wait()

	return nil
}

func getOrgs(cmd *cobra.Command, args []string) error {

  if user == "" {
    return fmt.Errorf("user flag is required")
  }

  client := getClient()

  opt := &github.ListOptions{PerPage: 10}
  orgs, _, err := client.Organizations.List(context.Background(), user, opt)
  if err != nil {
    return err
  }

  wg := &sync.WaitGroup{}
  for i := 0; i < threads; i++ {
    wg.Add(1)
    go func(start, end int) {
      defer wg.Done()
      for j := start; j < end; j++ {
        o := orgs[j]
        	fmt.Println("Name:", o.GetName())
      		fmt.Println("Description:", o.GetDescription())
      		fmt.Println("Location:", o.GetLocation())
      		fmt.Println()
      }
    }(i*len(orgs)/threads, (i+1)*len(orgs)/threads)
  }
  wg.Wait()


  return nil
}
