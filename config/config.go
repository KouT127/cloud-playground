package config

import (
	"cloud.google.com/go/compute/metadata"
	"log"
	"os"
)

var ProjectID string

func Configure() {
	log.SetFlags(0)
	ProjectID = os.Getenv("GOOGLE_CLOUD_PROJECT")
	if ProjectID == "" {
		ProjectID, _ = metadata.ProjectID()
	}
	if ProjectID == "" {
		log.Println("Could not determine Google Cloud Project. Running without log correlation. For local use set the GOOGLE_CLOUD_PROJECT environment variable.")
	}
}
