# One Comment
A small utility to keep a comment up to date on a PR when building continuously. 

```
$ onecomment -h
 -includes string
        The string to look for in the comment (default "<!-- Created by one-comment -->")
  -message string
        The comment message
  -owner string
        The owner of the GH repo
  -pr-id int
        The ID of the PR
  -repo string
        The GH repository
```

You also need to supply your GH token by setting the `GH_ACCESS_TOKEN` env var.

Inspired by: https://github.com/actions-cool/maintain-one-comment
