package image

import (
	"cloud.google.com/go/storage"
	"context"
	"encoding/json"
	"fmt"
	"github.com/nfnt/resize"
	"image"
	"image/png"
	"log"
)

const bucketName = "attendance-manament-d"

func NewStorageClient() (*storage.Client, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	return client, nil
}

type PubSubMessage struct {
	Data []byte `json:"data"`
}

type StorageInformation struct {
	FileName      string `json:"file_name"`
	FileExtension string `json:"file_extension"`
	Directory     string `json:"directory"`
	ImagePath     string `json:"image_path"`
}

func ResizeSubscriber(ctx context.Context, m PubSubMessage) error {
	var info StorageInformation
	err := json.Unmarshal(m.Data, &info)
	if err != nil {
		log.Printf("json parse error: %v", err)
		return err
	}
	client, err := NewStorageClient()
	if err != nil {
		log.Printf("client error: %v", err)
		return err
	}
	h := client.Bucket(bucketName).Object(info.ImagePath)
	r, err := h.NewReader(ctx)
	if err != nil {
		log.Printf("storage reader error: %v", err)
		return err
	}
	defer r.Close()
	image, _, err := image.Decode(r)
	if err != nil {
		log.Printf("image decode error: %v", err)
		return err
	}
	targetPath := fmt.Sprintf("%s/%s.%s", info.Directory, info.FileName+"1", info.FileExtension)

	thumb := resize.Thumbnail(300, 300, image, resize.NearestNeighbor)
	writer := client.Bucket(bucketName).Object(targetPath).NewWriter(ctx)
	if err := png.Encode(writer, thumb); err != nil {
		log.Printf("png encode error: %v", err)
		return err
	}
	writer.ACL = []storage.ACLRule{{Entity: storage.AllUsers, Role: storage.RoleReader}}
	writer.ContentType = fmt.Sprintf("image/png")
	writer.CacheControl = "public, max-age=86400"
	if err := writer.Close(); err != nil {
		log.Printf("storage writer error: %v", err)
		return err
	}

	return nil
}
