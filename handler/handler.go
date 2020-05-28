package handler

import (
	"cloud.google.com/go/storage"
	"context"
	"encoding/json"
	"fmt"
	"github.com/KouT127/cloud-playground/middleware"
	"github.com/KouT127/cloud-playground/model"
	"github.com/nfnt/resize"
	"image"
	"image/png"
	"io/ioutil"
	"log"
	"net/http"
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
	Message struct {
		Data []byte `json:"data,omitempty"`
		ID   string `json:"id"`
	} `json:"message"`
	Subscription string `json:"subscription"`
}

type StorageInformation struct {
	FileName      string `json:"file_name"`
	FileExtension string `json:"file_extension"`
	Directory     string `json:"directory"`
	ImagePath     string `json:"image_path"`
}

func ResizeSubscriptionHandler(w http.ResponseWriter, r *http.Request) {
	var (
		m    PubSubMessage
		info StorageInformation
	)
	ctx := r.Context()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("ioutil.ReadAll: %v", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	if err := json.Unmarshal(body, &m); err != nil {
		log.Printf("json.Unmarshal: %v", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	if err := json.Unmarshal(m.Message.Data, &info); err != nil {
		log.Printf("json parse error: %v json: %s", err, m.Message.Data)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	client, err := NewStorageClient()
	if err != nil {
		log.Printf("client error: %v", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	h := client.Bucket(bucketName).Object(info.ImagePath)
	reader, err := h.NewReader(ctx)
	if err != nil {
		log.Printf("storage reader error: %v", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	defer reader.Close()
	image, _, err := image.Decode(reader)
	if err != nil {
		log.Printf("image decode error: %v", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	targetPath := fmt.Sprintf("%s/%s.%s", info.Directory, info.FileName+"1", info.FileExtension)

	thumb := resize.Thumbnail(300, 300, image, resize.NearestNeighbor)
	writer := client.Bucket(bucketName).Object(targetPath).NewWriter(ctx)
	if err := png.Encode(writer, thumb); err != nil {
		log.Printf("png encode error: %v", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	writer.ACL = []storage.ACLRule{{Entity: storage.AllUsers, Role: storage.RoleReader}}
	writer.ContentType = fmt.Sprintf("image/png")
	writer.CacheControl = "public, max-age=86400"
	if err := writer.Close(); err != nil {
		log.Printf("storage writer error: %v", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	trace, ok := ctx.Value(middleware.CloudTraceContext).(string)
	if !ok {
		trace = ""
	}
	log.Println(model.Entry{
		Message:   "resize succeed",
		Component: "image-resize",
		Trace:     trace,
	})

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "ok")
	return
}

func TaskHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	trace, ok := ctx.Value(middleware.CloudTraceContext).(string)
	if !ok {
		trace = ""
	}
	log.Println(model.Entry{
		Message:   "task catch",
		Component: "task",
		Trace:     trace,
	})
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "ok")
	return
}
