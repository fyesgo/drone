package middleware

import (
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/drone/drone/remote"
	"github.com/drone/drone/remote/bitbucket"
	"github.com/drone/drone/remote/bitbucketserver"
	"github.com/drone/drone/remote/github"
	"github.com/drone/drone/remote/gitlab"
	"github.com/drone/drone/remote/gogs"
	"github.com/gin-gonic/gin"
)

// Remote is a middleware function that initializes the Remote and attaches to
// the context of every http.Request.
func Remote(c *cli.Context) gin.HandlerFunc {
	v, err := setupRemote(c)
	if err != nil {
		logrus.Fatalln(err)
	}
	return func(c *gin.Context) {
		remote.ToContext(c, v)
	}
}

// helper function to setup the remote from the CLI arguments.
func setupRemote(c *cli.Context) (remote.Remote, error) {
	switch {
	case c.Bool("github"):
		return setupGithub(c)
	case c.Bool("gitlab"):
		return setupGitlab(c)
	case c.Bool("bitbucket"):
		return setupBitbucket(c)
	case c.Bool("stash"):
		return setupStash(c)
	case c.Bool("gogs"):
		return setupGogs(c)
	default:
		return nil, fmt.Errorf("version control system not configured")
	}
}

// helper function to setup the Bitbucket remote from the CLI arguments.
func setupBitbucket(c *cli.Context) (remote.Remote, error) {
	return bitbucket.New(
		c.String("bitbucket-client"),
		c.String("bitbucket-server"),
	), nil
}

// helper function to setup the Gogs remote from the CLI arguments.
func setupGogs(c *cli.Context) (remote.Remote, error) {
	return gogs.New(gogs.Opts{
		URL:         c.String("gogs-server"),
		Username:    c.String("gogs-git-username"),
		Password:    c.String("gogs-git-password"),
		PrivateMode: c.Bool("gogs-private-mode"),
		SkipVerify:  c.Bool("gogs-skip-verify"),
	})
}

// helper function to setup the Stash remote from the CLI arguments.
func setupStash(c *cli.Context) (remote.Remote, error) {
	return bitbucketserver.New(bitbucketserver.Opts{
		URL:         c.String("stash-server"),
		Username:    c.String("stash-git-username"),
		Password:    c.String("stash-git-password"),
		ConsumerKey: c.String("stash-consumer-key"),
		ConsumerRSA: c.String("stash-consumer-rsa"),
		SkipVerify:  c.Bool("stash-skip-verify"),
	})
}

// helper function to setup the Gitlab remote from the CLI arguments.
func setupGitlab(c *cli.Context) (remote.Remote, error) {
	return gitlab.New(gitlab.Opts{
		URL:         c.String("gitlab-server"),
		Client:      c.String("gitlab-client"),
		Secret:      c.String("gitlab-sercret"),
		Username:    c.String("gitlab-git-username"),
		Password:    c.String("gitlab-git-password"),
		PrivateMode: c.Bool("gitlab-private-mode"),
		SkipVerify:  c.Bool("gitlab-skip-verify"),
	})
}

// helper function to setup the GitHub remote from the CLI arguments.
func setupGithub(c *cli.Context) (remote.Remote, error) {
	return github.New(
		c.String("github-server"),
		c.String("github-client"),
		c.String("github-sercret"),
		c.StringSlice("github-scope"),
		c.Bool("github-private-mode"),
		c.Bool("github-skip-verify"),
		c.BoolT("github-merge-ref"),
	)
}
