package storage

type URL struct {
	UUID     string `db:"uuid" json:"uuid"`
	Short    string `db:"short_url" json:"short_url"`
	Original string `db:"original_url" json:"original_url"`
}
