package storagemodel

type FileCopy struct {
	SourceFile File `json:"source_file"`
	File       File `json:"file"`
}
