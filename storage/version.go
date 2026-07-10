package storagemodel

import "time"

type Version struct {
	Id             string    `json:"id" example:"0027.795e1a4a-d925-4621-8d17-7dd90f85f1fd"`
	Number         uint32    `json:"number" example:"1"`
	FileId         string    `json:"file_id" example:"0027.8f5b88e0-b875-4aa7-a73e-08995671047e"`
	SizeInBytes    int64     `json:"size_in_bytes" example:"1543255"`
	Hash           string    `json:"hash" example:"5891b5b522d5df086d0ff0b110fbd9d21bb4fc7163af34d08286a2e846f6be03"`
	CreatedAt      time.Time `json:"created_at"`
	LastAccessedAt time.Time `json:"last_accessed_at"`
}
