package store

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Store is the data-access interface for all indexer and API persistence.
// All methods accept a context so callers can enforce deadlines and
// propagate cancellation.
type Store interface {
	// UpsertContract inserts or updates a contract row.
	UpsertContract(ctx context.Context, c Contract) error

	// GetContract returns the contract with the given ID, or ErrNotFound.
	GetContract(ctx context.Context, contractID string) (Contract, error)

	// ListContracts returns a cursor-paginated list of all tracked contracts.
	// cursor is opaque; pass "" for the first page. Returns the next cursor
	// (empty string when there are no more pages) as the second return value.
	ListContracts(ctx context.Context, cursor string, limit int) ([]Contract, string, error)

	// BatchInsertEvents inserts events, ignoring duplicates by primary key.
	// All rows are sent in a single network round-trip.
	BatchInsertEvents(ctx context.Context, events []Event) error

	// BatchInsertInvocations inserts invocations, ignoring duplicates.
	BatchInsertInvocations(ctx context.Context, invocations []Invocation) error

	// UpsertStorageEntries inserts or updates storage entries for a contract.
	UpsertStorageEntries(ctx context.Context, entries []StorageEntry) error

	// GetSyncState returns the sync cursor for a contract, or a zero-value
	// SyncState (LastLedger == 0) if no state has been recorded yet.
	GetSyncState(ctx context.Context, contractID string) (SyncState, error)

	// UpsertSyncState writes the sync cursor for a contract.
	UpsertSyncState(ctx context.Context, s SyncState) error

	// GetGlobalStats returns aggregate counts across all tracked contracts.
	// Computed with a single SQL query.
	GetGlobalStats(ctx context.Context) (GlobalStats, error)

	// RecordContractVersion appends a new entry to the contract_versions table
	// if the given wasm_hash has not been seen before for this contract.
	// It is a no-op (returns nil) when the (contract_id, wasm_hash) pair already
	// exists, making repeated indexer calls idempotent.
	RecordContractVersion(ctx context.Context, v ContractVersion) error

	// ListContractVersions returns all recorded Wasm hash entries for the given
	// contract, sorted chronologically by first_seen_ledger ascending.
	ListContractVersions(ctx context.Context, contractID string) ([]ContractVersion, error)

	// GetLatestContractVersion returns the most recently seen ContractVersion for
	// the given contract. Returns ErrNotFound when no version has been recorded yet.
	GetLatestContractVersion(ctx context.Context, contractID string) (ContractVersion, error)
}

// NewStore returns a Store backed by the given pgxpool.Pool.
func NewStore(pool *pgxpool.Pool) Store {
	return &postgresStore{pool: pool}
}
