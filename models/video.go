package models

// Reduced represents a video record.
type Reduced struct {
	FileName         string `json:"file_name"`
	Path             string `json:"path"`
	FileSize         string `json:"file_size"`
	ReducedMegabytes int64  `json:"reduced_megabytes"`
	Original         struct {
		FileName string `json:"file_name"`
		Path     string `json:"path"`
		FileSize string `json:"file_size"`
	}
}
