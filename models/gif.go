package models

type Gif struct {
	ContentURL string
}

type CustomSearchAPIResponse struct {
	Context struct {
		Title string `json:"title"`
	} `json:"context"`
	Items []struct {
		DisplayLink string `json:"displayLink"`
		HTMLSnippet string `json:"htmlSnippet"`
		HTMLTitle   string `json:"htmlTitle"`
		Image       struct {
			ByteSize        int    `json:"byteSize"`
			ContextLink     string `json:"contextLink"`
			Height          int    `json:"height"`
			ThumbnailHeight int    `json:"thumbnailHeight"`
			ThumbnailLink   string `json:"thumbnailLink"`
			ThumbnailWidth  int    `json:"thumbnailWidth"`
			Width           int    `json:"width"`
		} `json:"image"`
		Kind    string `json:"kind"`
		Link    string `json:"link"`
		Mime    string `json:"mime"`
		Snippet string `json:"snippet"`
		Title   string `json:"title"`
	} `json:"items"`
	Kind    string `json:"kind"`
	Queries struct {
		NextPage []struct {
			Count          int    `json:"count"`
			Cx             string `json:"cx"`
			InputEncoding  string `json:"inputEncoding"`
			OutputEncoding string `json:"outputEncoding"`
			Safe           string `json:"safe"`
			SearchTerms    string `json:"searchTerms"`
			SearchType     string `json:"searchType"`
			StartIndex     int    `json:"startIndex"`
			Title          string `json:"title"`
			TotalResults   string `json:"totalResults"`
		} `json:"nextPage"`
		Request []struct {
			Count          int    `json:"count"`
			Cx             string `json:"cx"`
			InputEncoding  string `json:"inputEncoding"`
			OutputEncoding string `json:"outputEncoding"`
			Safe           string `json:"safe"`
			SearchTerms    string `json:"searchTerms"`
			SearchType     string `json:"searchType"`
			StartIndex     int    `json:"startIndex"`
			Title          string `json:"title"`
			TotalResults   string `json:"totalResults"`
		} `json:"request"`
	} `json:"queries"`
	SearchInformation struct {
		FormattedSearchTime   string  `json:"formattedSearchTime"`
		FormattedTotalResults string  `json:"formattedTotalResults"`
		SearchTime            float64 `json:"searchTime"`
		TotalResults          string  `json:"totalResults"`
	} `json:"searchInformation"`
	URL struct {
		Template string `json:"template"`
		Type     string `json:"type"`
	} `json:"url"`
}
