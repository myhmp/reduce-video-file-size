package models

// Video represents a video record.
type Video struct {
	ReducedBytes int64 `json:"reduced_bytes"`
	Reduced      struct {
		FileName string `json:"file_name"`
		Path     string `json:"path"`
		Bytes    int64  `json:"bytes"`
	}
	Original struct {
		FileName string `json:"file_name"`
		Path     string `json:"path"`
		Bytes    int64  `json:"bytes"`
	}
}
