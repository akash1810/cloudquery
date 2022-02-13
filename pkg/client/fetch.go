package client

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/google/uuid"
)

// FetchSummary includes a summarized report of fetch, such as fetch id, fetch start and finish,
// resources fetch results
type FetchSummary struct {
	CqId uuid.UUID `db:"id"`
	//  Unique Id of fetch session
	FetchId            uuid.UUID              `db:"fetch_id"`
	CreatedAt          *time.Time             `db:"created_at"`
	Start              *time.Time             `db:"start"`
	Finish             *time.Time             `db:"finish"`
	IsSuccess          bool                   `db:"is_success"`
	TotalResourceCount uint64                 `db:"total_resource_count"`
	TotalErrorsCount   uint64                 `db:"total_errors_count"`
	ProviderName       string                 `db:"provider_name"`
	ProviderAlias      string                 `db:"provider_alias"`
	ProviderVersion    string                 `db:"provider_version"`
	CoreVersion        string                 `db:"core_version"`
	Resources          ResourceFetchSummaries `db:"results"`
}

type ResourceFetchSummaries []ResourceFetchSummary

// Value implements Valuer interface required by goqu
func (r ResourceFetchSummaries) Value() (driver.Value, error) {
	if len(r) == 0 {
		return nil, nil
	}
	return json.Marshal(r)
}

// ResourceFetchSummary includes a data about fetching specific resource
type ResourceFetchSummary struct {
	ResourceName string `json:"resource_name"`
	// map of resources that have finished fetching
	FinishedResources map[string]bool `json:"finished_resources"`
	// Error value if any, if returned the stream will be canceled
	Error string `json:"error"`
	// Execution status of resource
	Status string `json:"status"`
	// Total Amount of resources collected by this resource
	ResourceCount uint64 `json:"resource_count"`
}

// SaveFetchSummary saves fetch summary into fetches database
func (c *Client) SaveFetchSummary(ctx context.Context, fs *FetchSummary) error {
	id, err := uuid.NewUUID()
	if err != nil {
		return err
	}
	fs.CqId = id
	q := goqu.Dialect("postgres").Insert("cloudquery.fetches").Rows(fs)
	sql, args, err := q.ToSQL()
	if err != nil {
		return err
	}

	return c.db.Exec(ctx, sql, args...)
}
