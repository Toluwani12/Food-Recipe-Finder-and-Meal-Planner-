package maggi_ng

type SearchRecipeResponse struct {
	Status  string `json:"status"`
	Results struct {
		SearchResults []struct {
			TabCount   int `json:"tabCount"`
			Pagination struct {
				ItemsPerPage int `json:"items_per_page"`
				TotalPages   int `json:"total_pages"`
				CurrentPage  int `json:"current_page"`
			} `json:"pagination"`
			TabTitle      string `json:"tabTitle"`
			TabType       string `json:"tabType"`
			LoadMoreTitle string `json:"loadMoreTitle"`
			LoadMoreLink  string `json:"loadMoreLink"`
			PagerLink     string `json:"pagerLink"`
			Filters       []struct {
				FilterFacetName   string `json:"filterFacetName"`
				FilterFacetKey    string `json:"filterFacetKey"`
				FilterFacetValues []struct {
					FacetName string `json:"facetName"`
				} `json:"filterFacetValues"`
			} `json:"filters"`
			Sorting []struct {
				Sortingname string `json:"sortingname"`
				SortLink    string `json:"sortLink"`
			} `json:"sorting"`
			SearchResultList []struct {
				ContentTitle string `json:"contentTitle"`
				ContentLink  struct {
					Link  string `json:"link"`
					Title string `json:"title"`
				} `json:"contentLink"`
				ContentImage struct {
					Source struct {
						Xlarge string `json:"xlarge"`
						Large  string `json:"large"`
						Medium string `json:"medium"`
						Small  string `json:"small"`
					} `json:"source"`
					Alt string `json:"alt"`
				} `json:"contentImage"`
				ContentRating     float64 `json:"contentRating,omitempty"`
				ContentVote       int     `json:"contentVote,omitempty"`
				RecipeAjax        string  `json:"recipe_ajax"`
				Contenttime       string  `json:"contenttime"`
				Contentdifficulty string  `json:"contentdifficulty"`
				RecipeScore       int     `json:"recipe_score"`
				Type              string  `json:"type"`
			} `json:"searchResultList"`
		} `json:"searchResults"`
	} `json:"results"`
}

//type T struct {
//	Status  string `json:"status"`
//	Results struct {
//		SearchResults []struct {
//			TabCount   int `json:"tabCount"`
//			Pagination struct {
//				ItemsPerPage int `json:"items_per_page"`
//				TotalPages   int `json:"total_pages"`
//				CurrentPage  int `json:"current_page"`
//			} `json:"pagination"`
//			TabTitle         string `json:"tabTitle"`
//			TabType          string `json:"tabType"`
//			LoadMoreTitle    string `json:"loadMoreTitle"`
//			LoadMoreLink     string `json:"loadMoreLink"`
//			PagerLink        string `json:"pagerLink"`
//			SearchResultList []struct {
//				ContentTitle string `json:"contentTitle"`
//				ContentLink  struct {
//					Link  string `json:"link"`
//					Title string `json:"title"`
//				} `json:"contentLink"`
//				ContentImage struct {
//					Source struct {
//						Xlarge string `json:"xlarge"`
//						Large  string `json:"large"`
//						Medium string `json:"medium"`
//						Small  string `json:"small"`
//					} `json:"source"`
//					Alt string `json:"alt"`
//				} `json:"contentImage"`
//				ContentRating     float64 `json:"contentRating,omitempty"`
//				ContentVote       int     `json:"contentVote,omitempty"`
//				RecipeAjax        string  `json:"recipe_ajax"`
//				Contenttime       string  `json:"contenttime"`
//				Contentdifficulty string  `json:"contentdifficulty"`
//				RecipeScore       int     `json:"recipe_score"`
//				Type              string  `json:"type"`
//			} `json:"searchResultList"`
//		} `json:"searchResults"`
//	} `json:"results"`
//}
