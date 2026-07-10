package dataapi

import (
	"encoding/json"
	"go.proteos.ai/model/common"
	"go.proteos.ai/model/data"
)

// CreateRecordRequest represents the request to create a new record
type CreateRecordRequest struct {
	Data map[string]any `json:"data" validate:"required"`
}

// UpdateRecordRequest represents the request to update an existing record
type UpdateRecordRequest struct {
	Data map[string]any `json:"data" validate:"required"`
}

// GetManyRecordsQuery represents the query parameters for listing records
type GetManyRecordsQuery struct {
	Query *common.ClientQuery `json:"query,omitempty"`
	common.Pagination
	common.Sorting
}

// GetManyRecordsResponse represents the response for listing records
type GetManyRecordsResponse struct {
	Meta common.ResponseMeta `json:"meta"`
	Data []datamodel.Record  `json:"data"`
}

// === Batch Operations ===

// BatchTransactionStatus represents the status of a batch transaction
type BatchTransactionStatus string

const (
	BatchTransactionStatusSuccess BatchTransactionStatus = "success"
	BatchTransactionStatusError   BatchTransactionStatus = "error"
)

// BatchTransactionError represents an error in a batch transaction
type BatchTransactionError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// --- Batch Upsert ---

// BatchUpsertTransaction represents a single transaction in a batch upsert request
type BatchUpsertTransaction struct {
	TransactionId string          `json:"transaction_id"`
	Data          json.RawMessage `json:"data"`
}

// BatchUpsertTransactionResult represents the result of a single batch upsert transaction
type BatchUpsertTransactionResult struct {
	TransactionId string                 `json:"transaction_id"`
	Status        BatchTransactionStatus `json:"status"`
	Record        datamodel.Record       `json:"record,omitempty"`
	Error         *BatchTransactionError `json:"error,omitempty"`
}

// BatchUpsertRecordsResponse represents the response for a batch upsert operation
type BatchUpsertRecordsResponse struct {
	Results []BatchUpsertTransactionResult `json:"results"`
}

// --- Batch Create ---

// BatchCreateTransaction represents a single transaction in a batch create request
type BatchCreateTransaction struct {
	TransactionId string          `json:"transaction_id"`
	Data          json.RawMessage `json:"data"`
}

// BatchCreateTransactionResult represents the result of a single batch create transaction
type BatchCreateTransactionResult struct {
	TransactionId string                 `json:"transaction_id"`
	Status        BatchTransactionStatus `json:"status"`
	Record        datamodel.Record       `json:"record,omitempty"`
	Error         *BatchTransactionError `json:"error,omitempty"`
}

// BatchCreateRecordsResponse represents the response for a batch create operation
type BatchCreateRecordsResponse struct {
	Results []BatchCreateTransactionResult `json:"results"`
}

// --- Batch Update ---

// BatchUpdateTransaction represents a single transaction in a batch update request
type BatchUpdateTransaction struct {
	TransactionId string          `json:"transaction_id"`
	Id            string          `json:"id"`
	Data          json.RawMessage `json:"data"`
}

// BatchUpdateTransactionResult represents the result of a single batch update transaction
type BatchUpdateTransactionResult struct {
	TransactionId string                 `json:"transaction_id"`
	Status        BatchTransactionStatus `json:"status"`
	Record        datamodel.Record       `json:"record,omitempty"`
	Error         *BatchTransactionError `json:"error,omitempty"`
}

// BatchUpdateRecordsResponse represents the response for a batch update operation
type BatchUpdateRecordsResponse struct {
	Results []BatchUpdateTransactionResult `json:"results"`
}

// --- Batch Delete ---

// BatchDeleteTransaction represents a single transaction in a batch delete request
type BatchDeleteTransaction struct {
	TransactionId string `json:"transaction_id"`
	Id            string `json:"id"`
}

// BatchDeleteTransactionResult represents the result of a single batch delete transaction
type BatchDeleteTransactionResult struct {
	TransactionId string                 `json:"transaction_id"`
	Status        BatchTransactionStatus `json:"status"`
	Error         *BatchTransactionError `json:"error,omitempty"`
}

// BatchDeleteRecordsResponse represents the response for a batch delete operation
type BatchDeleteRecordsResponse struct {
	Results []BatchDeleteTransactionResult `json:"results"`
}
