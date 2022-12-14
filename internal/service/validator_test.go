package service

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/Namchee/konfigured/internal/entity"
	"github.com/Namchee/konfigured/mocks/mock_client"
	"github.com/golang/mock/gomock"
	"github.com/google/go-github/v48/github"
	"github.com/stretchr/testify/assert"
)

func TestNewConfigurationValidator(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := mock_client.NewMockGithubClient(ctrl)

	assert.NotPanics(t, func() {
		NewConfigurationValidator(&entity.Configuration{}, client)
	})
}

func TestConfigurationValidator_ValidateConfigurationFiles(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type response struct {
		content *github.RepositoryContent
		err     error
	}

	files := map[string]response{
		"foobar.json": {
			content: &github.RepositoryContent{
				Content: github.String("{\n"),
			},
			err: nil,
		},
		"sample.toml": {
			content: &github.RepositoryContent{
				Content: github.String(`key = "value"
`),
			},
			err: nil,
		},
		"config.yaml": {
			content: &github.RepositoryContent{
				Content: github.String("key: value\n"),
			},
			err: nil,
		},
		"configuration.hcl": {
			content: &github.RepositoryContent{
				Content: github.String(`example {
  foo = "bar"
}
`),
			},
			err: nil,
		},
		"invalid.hcl": {
			content: &github.RepositoryContent{
				Content: github.String(`example {
  foo = "
}`),
			},
			err: nil,
		},
		"no-newline.yaml": {
			content: &github.RepositoryContent{
				Content: github.String("key: value"),
			},
			err: nil,
		},
		"nested/config.yaml": {
			content: &github.RepositoryContent{
				Content: github.String(""),
			},
			err: errors.New("fail"),
		},
		"encoding.ini": {
			content: &github.RepositoryContent{
				Content:  github.String(""),
				Encoding: github.String("magic"),
			},
			err: nil,
		},
	}

	args := []*github.CommitFile{
		{
			Filename: github.String("foobar.json"),
		},
		{
			Filename: github.String("sample.toml"),
		},
		{
			Filename: github.String("config.yaml"),
		},
		{
			Filename: github.String("configuration.hcl"),
		},
		{
			Filename: github.String("invalid.hcl"),
		},
		{
			Filename: github.String("no-newline.yaml"),
		},
		{
			Filename: github.String("nested/config.yaml"),
		},
		{
			Filename: github.String("not-included.yml"),
		},
		{
			Filename: github.String("encoding.ini"),
		},
		{
			Filename: github.String("picture.png"),
		},
		{
			Filename: github.String("README"),
		},
	}

	client := mock_client.NewMockGithubClient(ctrl)

	for filename, resp := range files {
		client.EXPECT().GetFileContent(gomock.Any(), filename).
			Return(resp.content, resp.err)
	}

	validator := &ConfigurationValidator{
		cfg: &entity.Configuration{
			Newline: true,
			Include: "**/*.{json,ini,yaml,toml,hcl}",
		},
		client: client,
	}

	got := validator.ValidateFiles(context.TODO(), args)

	assert.Equal(t, 8, len(got))

	invalids := entity.GetInvalidValidations(got)

	fmt.Println(invalids)

	assert.Equal(t, 5, len(invalids))
}
