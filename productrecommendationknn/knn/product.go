package knn

// Product represent Product data
type Product struct {
	ProductName string   `json:"productName"`
	Url         string   `json:"url"`
	Categories  []string `json:"categories"`
	CategoryNum []int    `json:"-"`
}

// Products represent Product collections
type Products []*Product
