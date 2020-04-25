package models

// Video represents a video record.
type Video struct {
	ReducedMegabytes float64 `json:"reduced_megabytes"`
	Reduced          struct {
		FileName  string  `json:"file_name"`
		Path      string  `json:"path"`
		Megabytes float64 `json:"megabytes"`
	}
	Original struct {
		FileName  string  `json:"file_name"`
		Path      string  `json:"path"`
		Megabytes float64 `json:"megabytes"`
	}
}

// Record represents a video record.
type Record struct {
	ReducedMegabytes float64 `json:"reduced_megabytes"`
	Videos           []Video `json:"items"`
}
