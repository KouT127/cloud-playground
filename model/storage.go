package model

type StorageInformation struct {
	FileName      string `json:"file_name"`
	FileExtension string `json:"file_extension"`
	Directory     string `json:"directory"`
	ImagePath     string `json:"image_path"`
}
