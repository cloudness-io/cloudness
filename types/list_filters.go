package types

// ListQueryFilter has pagination related info and a query param.
type ListQueryFilter struct {
	Pagination
	Query string `json:"query"`
}

type CreatedFilter struct {
	CreatedGt int64 `json:"created_gt"`
	CreatedLt int64 `json:"created_lt"`
}

type UpdatedFilter struct {
	UpdatedGt int64 `json:"updated_gt"`
	UpdatedLt int64 `json:"updated_lt"`
}
