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

// === Tenant Schema Sweep ===

// SweepTenantSchemasRequest triggers the tenant-schema sweep.
//
// It carries NO DDL. The columns and indexes applied are declared in reviewed Go
// source (data-service's postgres.TenantSchemaAdditions), because record tables
// live in per-org schemas created by runtime DDL and accepting column
// definitions over HTTP would make this endpoint an arbitrary-DDL surface.
type SweepTenantSchemasRequest struct {
	// OrgIds optionally narrows the sweep to specific tenants. Empty sweeps every
	// org_* schema — the normal case, since the point is to reach all of them.
	OrgIds []string `json:"org_ids,omitempty"`
	// IsDryRun returns the statements that WOULD run without executing any of
	// them. Run this first on production-sized data.
	IsDryRun bool `json:"is_dry_run,omitempty"`
}

// SweepFailure is one table the sweep could not fully apply. The sweep continues
// past failures so one poisoned table cannot block the remaining tenants.
type SweepFailure struct {
	Schema    string `json:"schema"`
	Table     string `json:"table"`
	Statement string `json:"statement"`
	Error     string `json:"error"`
}

// SweepTenantSchemasReport is the outcome of one sweep run. The sweep is
// idempotent and resumable, so a run with failures is resolved by fixing the
// cause and running it again.
type SweepTenantSchemasReport struct {
	IsDryRun       bool `json:"is_dry_run"`
	TablesScanned  int  `json:"tables_scanned"`
	ColumnsApplied int  `json:"columns_applied"`
	IndexesApplied int  `json:"indexes_applied"`
	// ShareTablesEnsured counts schemas whose per-tenant access_grants table
	// was ensured this run (created or already present).
	ShareTablesEnsured int `json:"share_tables_ensured"`
	// ColumnsRetired counts (table × retired column) drops that ran cleanly.
	ColumnsRetired int `json:"columns_retired"`
	// Statements is populated on a dry run only.
	Statements []string       `json:"statements,omitempty"`
	Failures   []SweepFailure `json:"failures,omitempty"`
	Message    string         `json:"message"`
}

