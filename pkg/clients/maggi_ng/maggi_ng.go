package maggi_ng

import (
	client "Food/pkg/clients"
	"Food/pkg/recipe/model"
	"context"
	"encoding/json"
	"fmt"
	"github.com/chromedp/chromedp"
	log "github.com/sirupsen/logrus"
)

type SearchParams struct {
	ContentType string `json:"content_type"`
	Range       int    `json:"range"`
	SearchPage  bool   `json:"searchpage"`
	Filters     bool   `json:"filters"`
}

type Client struct {
	*client.Client
}

// NewClient creates a new API client for the Maggi NG service wrapped in a custom type.
func NewClient() *Client {
	cl := client.NewClient()
	cl.SetBaseURL("https://www.maggi.ng")
	return &Client{Client: cl}
}

// FetchData method defined on the Client type.
func (c *Client) FetchLinks(params SearchParams) (*model.Request, error) {
	_, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	for page := 0; ; page++ {
		var response SearchRecipeResponse

		err := c.Get(fmt.Sprintf("/search-ajax-results?content_type=srh_recipe&page=%d&range=12&searchpage=false", page), params, &response) // Ensure response is a pointer
		if err != nil {
			log.Printf("Failed to fetch recipes: %v", err)
			return nil, err
		}

		println(response.Results.SearchResults[0].Pagination.TotalPages)
		println(response.Results.SearchResults[0].SearchResultList)

		//res := model.Request{}
		for _, result := range response.Results.SearchResults {
			for _, link := range result.SearchResultList {
				println(link.ContentLink.Link)
			}
			//	link := result.SearchResultList[0].ContentLink.Link
			//	//recipe, err := c.FetchRecipe(ctx, string(link))
			//	if err != nil {
			//		log.Printf("Failed to fetch recipe: %v", err)
			//	}
			//	res = append(res, *recipe)
		}

		if page >= response.Results.SearchResults[0].Pagination.TotalPages {
			break
		}
	}

	return nil, nil
}

func (c *Client) FetchRecipe(ctx context.Context, url string) (*model.RequestData, error) {

	log.Println("Fetching recipe details for URL: ", url)

	// Run the browser and navigate to the page
	var htmlContent string
	err := chromedp.Run(ctx,
		chromedp.Navigate(c.GetBaseURL()+url),
		//chromedp.Sleep(5*time.Second), // Wait for the page to load fully
		chromedp.OuterHTML("html", &htmlContent),
	)
	if err != nil {
		log.Fatalf("Failed to navigate and get HTML content: %v", err)
	}
	fmt.Println("HTML content loaded successfully")

	// Extract the details
	var recipe model.RequestData
	err = chromedp.Run(ctx,
		chromedp.WaitVisible(`h1[data-recipeid]`, chromedp.ByQuery),
		chromedp.Text(`h1[data-recipeid]`, &recipe.Name, chromedp.NodeVisible),
		chromedp.Text(`div.mg-recipe-instructions__header > h2`, &recipe.Description, chromedp.NodeVisible),
		chromedp.Text(`div.recipe-inst-totalmins`, &recipe.CookingTime, chromedp.NodeVisible),
		chromedp.Evaluate(`Array.from(document.querySelectorAll('li.mg-recipe-instructions__steps')).map(step => step.innerText)`, &recipe.Instructions),
		chromedp.AttributeValue(`div.mg-stage-carousel__slide img`, "src", &recipe.ImgUrl, nil),
	)
	if err != nil {
		log.Fatalf("Failed to extract recipe details: %v", err)
	}
	fmt.Println("Basic recipe details extracted successfully")

	// Extract ingredients
	var ingredients []map[string]string
	err = chromedp.Run(ctx,
		chromedp.WaitVisible(`div.mg-ingredients__table`, chromedp.ByQuery),
		chromedp.Evaluate(`Array.from(document.querySelectorAll('div.mg-ingredients__table > div[role="row"]')).map(row => {
			let nameElem = row.querySelector('div.mg-span.recipe-ingredient-name em');
			let quantityElem = row.querySelector('div.mg-span.recipe-ingredient-unit em.mg-ingredient-quantity');
			let unitElem = row.querySelector('div.mg-span.recipe-ingredient-unit em.mg-ingredient-unit');

			return {
				name: nameElem ? nameElem.innerText.trim() : "",
				quantity: quantityElem ? (quantityElem.innerText.trim() + " " + (unitElem ? unitElem.innerText.trim() : "")) : ""
			};
		})`, &ingredients),
	)
	if err != nil {
		log.Fatalf("Failed to extract ingredients: %v", err)
	}

	// Print the extracted ingredients for debugging
	fmt.Printf("Extracted ingredients: %+v\n", ingredients)

	// Convert the ingredients to the desired format
	for _, ing := range ingredients {
		recipe.Ingredients = append(recipe.Ingredients, model.IngredientRequest{
			Name:     ing["name"],
			Quantity: ing["quantity"],
		})
	}
	fmt.Println("Ingredients extracted successfully")

	// Print the extracted recipe as JSON
	recipeJSON, err := json.MarshalIndent(recipe, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal recipe to JSON: %v", err)
	}
	fmt.Println(string(recipeJSON))

	return &recipe, nil
}
