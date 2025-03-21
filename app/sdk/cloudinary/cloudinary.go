package cloudinary

import (
	"context"
	"fmt"
	"io"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

// CloudinaryService struct to hold the Cloudinary instance
type CloudinaryService struct {
	client *cloudinary.Cloudinary
}

// NewCloudinaryService creates a new CloudinaryService instance
func NewCloudinaryService(url string) (*CloudinaryService, error) {
	cld, err := cloudinary.NewFromURL(url)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize cloudinary: %w", err)
	}

	return &CloudinaryService{client: cld}, nil
}

// UploadMedia uploads a file to Cloudinary
func (c *CloudinaryService) UploadMedia(file io.Reader, filename string) (*uploader.UploadResult, error) {
	uploadParams := uploader.UploadParams{
		PublicID:     filename,
		ResourceType: "auto",
	}

	result, err := c.client.Upload.Upload(context.Background(), file, uploadParams)
	if err != nil {
		return nil, fmt.Errorf("error uploading to cloudinary: %w", err)
	}
	return result, nil
}

// DeleteMedia deletes a file from Cloudinary
func (c *CloudinaryService) DeleteMedia(publicId string) error {
	_, err := c.client.Upload.Destroy(context.Background(), uploader.DestroyParams{PublicID: publicId})
	if err != nil {
		return fmt.Errorf("failed to delete asset from cloudinary: %w", err)
	}
	return nil
}
