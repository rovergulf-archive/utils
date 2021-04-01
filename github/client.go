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
	User     string             `json:"user" yaml:"user"`
	Org      string             `json:"org" yaml:"org"`
}

type Client struct {
	ctx    context.Context
	logger *zap.SugaredLogger
	cfg    *Options
	*github.Client
}

func New(o *Options) *Client {
	return &Client{
		logger: o.Logger,
		cfg:    o,
	}
}

func (c *Client) Run(ctx context.Context) error {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: c.cfg.ApiToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	c.Client = github.NewClient(tc)

	return nil
}

func (c *Client) ListOrganizations(ctx context.Context, opts *github.ListOptions) ([]*github.Organization, error) {
	orgs, _, err := c.Organizations.List(ctx, c.cfg.Org, opts)
	if err != nil {
		return nil, err
	}

	return orgs, nil
}

func (c *Client) ListRepositories(ctx context.Context, opts *github.RepositoryListOptions) ([]*github.Repository, error) {
	repos, _, err := c.Repositories.List(ctx, c.cfg.User, opts)
	if err != nil {
		return nil, err
	}

	return repos, nil
}

func (c *Client) ListOrgRepositories(ctx context.Context, opts *github.RepositoryListByOrgOptions) ([]*github.Repository, error) {
	repos, _, err := c.Repositories.ListByOrg(ctx, c.cfg.Org, opts)
	if err != nil {
		return nil, err
	}

	return repos, nil
}
