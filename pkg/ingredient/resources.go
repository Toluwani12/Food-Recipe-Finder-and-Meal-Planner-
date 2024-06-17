package ingredient

import (
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
)

type Resource struct {
	db *sqlx.DB
}

// NewResource creates and returns a resource.
func NewResource(db *sqlx.DB) *Resource {
	return &Resource{
		db: db,
	}
}

func (rs *Resource) Router() *chi.Mux {
	r := chi.NewRouter()

	repo := NewRepository(rs.db)
	svc := NewService(repo)
	hndlr := NewHandler(svc)

	//r.Use(auth.AuthMiddleware)
	//r.Post("/add", hndlr.add)
	//r.Delete("/delete/{id}", hndlr.delete)
	//r.Put("/update/{id}", hndlr.update)
	//r.Get("/get/{id}", hndlr.get)
	r.Get("/", hndlr.list)

	return r
}
