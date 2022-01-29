package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/google/go-github/github"
	_ "github.com/joho/godotenv/autoload"
	"golang.org/x/oauth2"
)

func main() {

	var (
		prNum           int
		owner           string
		repo            string
		commentIncludes string
		ghAccessToken   string
		message         string
	)

	ghAccessToken = os.Getenv("GH_ACCESS_TOKEN")

	if ghAccessToken == "" {
		fmt.Println("Must supply env var GH_ACCESS_TOKEN")
		os.Exit(1)
	}

	flag.IntVar(&prNum, "pr-id", 0, "The ID of the PR")
	flag.StringVar(&owner, "owner", "", "The owner of the GH repo")
	flag.StringVar(&repo, "repo", "", "The GH repository")
	flag.StringVar(&commentIncludes, "includes", "<!-- Created by one-comment -->", "The string to look for in the comment")
	flag.StringVar(&message, "message", "", "The comment message")

	flag.Parse()

	if prNum == 0 {
		fmt.Println("Please specify the pull request ID using the -pr-id argument")
		os.Exit(1)
	}
	if owner == "" {
		fmt.Println("Please specify the repo owner with the -owner argument")
		os.Exit(1)
	}

	if repo == "" {
		fmt.Println("Please specify the repo with the -repo flag")
		os.Exit(1)
	}

	if message == "" {
		fmt.Println("Please supply a message with the -message flag")
		os.Exit(1)
	}

	err := findOrCreateIssueComment(ghAccessToken, owner, repo, commentIncludes, message, prNum)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}

func findOrCreateIssueComment(ghAccessToken, owner, repo, commentIncludes, message string, prNum int) error {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: ghAccessToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	comments, _, err := client.Issues.ListComments(ctx, owner, repo, prNum, &github.IssueListCommentsOptions{})

	if err != nil {
		return fmt.Errorf("Error retrieving existing comments: %w", err)
	}

	var foundComment *github.IssueComment
	for _, comment := range comments {
		if strings.Contains(*comment.Body, commentIncludes) {
			foundComment = comment
			break
		}
	}
	message = message + commentIncludes

	if foundComment != nil {
		_, _, err = client.Issues.EditComment(ctx, owner, repo, *foundComment.ID, &github.IssueComment{Body: &message})
		if err != nil {
			return fmt.Errorf("Error editing comment: %w", err)
		}
		fmt.Println("Comment found and edited")
	} else {
		_, _, err = client.Issues.CreateComment(ctx, owner, repo, prNum, &github.IssueComment{Body: &message})
		if err != nil {
			return fmt.Errorf("Error creating comment: %w", err)
		}
		fmt.Println("Comment created")
	}

	return nil
}
