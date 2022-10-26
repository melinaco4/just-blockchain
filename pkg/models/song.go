package song

type Song struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Artist      string `json:"artist"`
	ReleaseDate string `json:"release_date"`
	SKU         string `json:"sku"`
}
