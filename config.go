package ffclient

import (
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"net/http"

	"github.com/thomaspoignant/go-feature-flag/internal/retriever"
)

// Config is the configuration of go-feature-flag.
// PollInterval is the interval in seconds where we gonna read the file to update the cache.
// You should also have a retriever to specify where to read the flags file.
type Config struct {
	PollInterval  int // Poll every X seconds
	LocalFile     string
	HTTPRetriever *HTTPRetriever
	S3Retriever   *S3Retriever
}

// HTTPRetriever is a configuration struct for an HTTP endpoint retriever.
type HTTPRetriever struct {
	URL    string
	Method string
	Body   string
	Header http.Header
}

// S3Retriever is a configuration struct for a S3 retriever.
type S3Retriever struct {
	Bucket    string
	Item      string
	AwsConfig aws.Config
}

// GetRetriever is used to get the retriever we will use to load the flags file.
func (c *Config) GetRetriever() (retriever.FlagRetriever, error) {
	if c.S3Retriever != nil {
		// Create an AWS session
		sess, err := session.NewSession(&c.S3Retriever.AwsConfig)
		if err != nil {
			return nil, err
		}

		// Create a new AWS S3 downloader
		downloader := s3manager.NewDownloader(sess)
		return retriever.NewS3Retriever(
			downloader,
			c.S3Retriever.Bucket,
			c.S3Retriever.Item,
		), nil
	}

	if c.HTTPRetriever != nil {
		return retriever.NewHTTPRetriever(
			http.DefaultClient,
			c.HTTPRetriever.URL,
			c.HTTPRetriever.Method,
			c.HTTPRetriever.Body,
			c.HTTPRetriever.Header,
		), nil
	}

	if c.LocalFile != "" {
		return retriever.NewLocalRetriever(c.LocalFile), nil
	}
	return nil, errors.New("please add a config to get the flag config file")
}