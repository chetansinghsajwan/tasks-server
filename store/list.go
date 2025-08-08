package store

type List struct {
	ID   uint64 `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
}

type CreateListParams struct {
	Name string `db:"name" json:"name"`
}

type UpdateListParams struct {
	Name *string `db:"name" json:"name"`
}
