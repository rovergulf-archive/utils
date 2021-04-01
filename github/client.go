package ghx

import (
	"context"
	"github.com/google/go-github/v33/github"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
)

type Options struct {
	Logger   *zap.SugaredLogger `json:"-" yaml:"-"`
	ApiToken string             `json:"api_token" yaml:"api_token"`
}

type Client struct {
	ctx    context.Context
	logger *zap.SugaredLogger
	token  string
	*github.Client
}

func New(o *Options) *Client {
	return &Client{
		logger: o.Logger,
		token:  o.ApiToken,
	}
}

func (c *Client) Run(ctx context.Context) error {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: c.token},
	)
	tc := oauth2.NewClient(ctx, ts)

	c.Client = github.NewClient(tc)

	return nil
}

func (c *Client) ListOrganizations() ([]*github.Organization, error) {
	orgs, _, err := c.Organizations.List(context.Background(), "scDisorder", nil)
	if err != nil {
		return nil, err
	}

	return orgs, nil
}

func (c *Client) ListRepositories() ([]*github.Repository, error) {
	repos, _, err := c.Repositories.ListByOrg(context.Background(), "rovergulf", nil)
	if err != nil {
		return nil, err
	}

	return repos, nil
}
