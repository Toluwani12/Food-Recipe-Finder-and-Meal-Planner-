package recipe

type ListResponse struct {
	ID    string `db:"id" json:"id"`
	Name  string `db:"name" json:"name"`
	Liked bool   `db:"liked" json:"liked"`
}
