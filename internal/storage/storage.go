package storage

import (
	"kube/internal/config"
)

type Client struct {
	config *config.Config
}

func Init(cfg *config.Config) *Client {
	return &Client{
		config: cfg,
	}
}

func (c *Client) UploadFile(filePath string, data []byte) error {
	// TODO: Implement file upload logic
	return nil
}

func (c *Client) DownloadFile(filePath string) ([]byte, error) {
	// TODO: Implement file download logic
	return nil, nil
}

func (c *Client) DeleteFile(filePath string) error {
	// TODO: Implement file deletion logic
	return nil
}
