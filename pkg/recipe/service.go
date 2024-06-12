package recipe

import (
	liberror "Food/internal/errors"
	"Food/pkg"
	"Food/pkg/recipe/crawler"
	"Food/pkg/recipe/model"
	"context"
	"errors"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"net/url"
)

type Service struct {
	repo        *Repository
	crawlerList []crawler.ICrawler
}

func NewService(repo *Repository, crawlerList []crawler.ICrawler) *Service {
	return &Service{
		repo:        repo,
		crawlerList: crawlerList,
	}
}

func (s Service) crawl(ctx context.Context) (map[string]bool, error) {
	recipeList := make([]model.RequestData, 0, len(s.crawlerList))
	for _, c := range s.crawlerList {
		data, err := c.CrawlRecipe()
		if err != nil {
			log.WithFields(log.Fields{"service": "recipes/fetchAndSaveRecipes", "crawler": c}).WithError(err)
		}
		recipeList = append(recipeList, *data...)
	}
	//marshal, err := json.Marshal(recipeList)
	//if err != nil {
	//	return nil, err
	//}

	//println(string(marshal))
	req := model.Request(recipeList)

	return s.save(ctx, req)
}

// AddRecipes adds a list of new recipes along with their ingredients to the database.
func (s Service) save(ctx context.Context, recipes model.Request) (map[string]bool, error) {

	// Ensure all IDs are set for recipes, ingredients, and their alternatives
	for i := range recipes {
		if recipes[i].ID == uuid.Nil {
			recipes[i].ID = uuid.New() // Ensure the recipe has a UUID
		}
	}
	successMap, err := s.repo.processRecipesAndIngredients(ctx, recipes)
	return successMap, liberror.CoverErr(err,
		errors.New("service temporarily unavailable. Please try again later"),
		log.WithFields(log.Fields{"service": "recipes/add",
			"repo": "recipes/add",
		}).WithError(err))
}

//func (s Service) update(ctx context.Context, id string, data AddRequest) (*Recipe, error) {
//	recipe := Recipe{
//		Id:           id,
//		Name:         data.Name,
//		CookingTime:  data.CookingTime,
//		Instructions: data.Instructions,
//		UpdatedAt:    time.Now(),
//	}
//	resp, err := s.repo.update(ctx, recipe)
//	return resp, liberror.CoverErr(err,
//		errors.New("service temporarily unavailable. Please try again later"),
//		log.WithFields(log.Fields{"service": "recipes/update", "repo": "recipes/update"}).WithError(err))
//}

func (s Service) delete(ctx context.Context, id string) (string, error) {
	resp, err := s.repo.delete(ctx, id)
	return resp, liberror.CoverErr(err,
		errors.New("service temporarily unavailable. Please try again later"),
		log.WithFields(log.Fields{"service": "recipes/delete", "repo": "recipes/delete"}).WithError(err))
}

func (s Service) get(ctx context.Context, id string) (*Recipe, error) {
	resp, err := s.repo.get(ctx, id)
	return resp, liberror.CoverErr(err,
		errors.New("service temporarily unavailable. Please try again later"),
		log.WithFields(log.Fields{"service": "recipes/get", "repo": "recipes/get"}).WithError(err))
}

func (s Service) list(ctx context.Context, userID string, recipeName string) ([]ListResponse, error) {
	resp, err := s.repo.list(ctx, userID, recipeName)
	return resp, liberror.CoverErr(err,
		errors.New("service temporarily unavailable. Please try again later"),
		log.WithFields(log.Fields{"service": "recipes/list", "repo": "recipes/list"}).WithError(err))
}

func (s Service) search(ctx context.Context, ingredients []string, queryParams url.Values, userId string) (interface{}, *pkg.Pagination, error) {
	recipes, pg, err := s.repo.search(ctx, ingredients, queryParams, userId)
	return recipes, pg, liberror.CoverErr(err,
		errors.New("service temporarily unavailable. Please try again later"),
		log.WithFields(log.Fields{"service": "recipes/findRecipes", "repo": "recipes/findRecipes"}).WithError(err))
}
