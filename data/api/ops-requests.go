package dataapi

import (
	"go.proteos.ai/model/meta"
)

// === Tenant Operations ===

// CreateTenantRequest represents the request to create a new tenant schema
type CreateTenantRequest struct {
	OrgId string `json:"org_id" binding:"required" validate:"required"`
}

// CreateTenantResponse represents the response after creating a tenant
type CreateTenantResponse struct {
	Success bool   `json:"success"`
	Schema  string `json:"schema"`
	Message string `json:"message"`
}

// DeleteTenantResponse represents the response after deleting a tenant
type DeleteTenantResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// === Table Operations ===

// CreateTableRequest represents the request to create a new table for an entity.
// The orgId is taken from entity.OrgId.
type CreateTableRequest struct {
	Entity metamodel.EntityWithSchema `json:"entity" binding:"required" validate:"required"`
}

// CreateTableResponse represents the response after creating a table
type CreateTableResponse struct {
	Success   bool   `json:"success"`
	TableName string `json:"table_name"`
	Message   string `json:"message"`
}

// UpdateTableRequest represents the request to update entity in cache (no table changes).
// The orgId is taken from entity.OrgId.
type UpdateTableRequest struct {
	Entity metamodel.EntityWithSchema `json:"entity" binding:"required" validate:"required"`
}

// UpdateTableResponse represents the response after updating a table
type UpdateTableResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
