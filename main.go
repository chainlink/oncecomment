package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/google/go-github/github"
	_ "github.com/joho/godotenv/autoload"
	"golang.org/x/oauth2"
)

func main() {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GH_ACCESS_TOKEN")},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)
	owner := ""
	repo := ""
	includes := "<!-- Created by one-comment -->"
	issueNum := 00
	comments, _, err := client.Issues.ListComments(ctx, owner, repo, issueNum, &github.IssueListCommentsOptions{})

	if err != nil {
		panic(err)
	}

	var cid int64
	for _, comment := range comments {
		if strings.Contains(*comment.Body, includes) {
			fmt.Printf("%+v\n", comment)
			cid = *comment.ID
			break
		}
	}

	fmt.Printf("comment ID:%d\n", cid)
	newBody := "blah"
	_, _, err = client.Issues.EditComment(ctx, owner, repo, cid, &github.IssueComment{Body: &newBody})

	if err != nil {
		panic(err)
	}

}
