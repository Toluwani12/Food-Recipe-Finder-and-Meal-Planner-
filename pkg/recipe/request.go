package recipe

import (
	log "github.com/sirupsen/logrus"
	"net/http"
)

type ListResponse struct {
	ID    string `db:"id" json:"id"`
	Name  string `db:"name" json:"name"`
	Liked bool   `db:"liked" json:"liked"`
}

type SearchRequest struct {
	Ingredients []string `json:"ingredients"`
}

func (s SearchRequest) Bind(r *http.Request) error {
	log.Printf("Binding search request: %v", s)
	return nil
}
