package song-checkout

type SongCheckout struct {
	SongID       string `json:"song_id"`
	User         string `json:"user"`
	CheckoutDate string `json:"checkout_date"`
	IsGenesis    bool   `jsons:"is_genesis"`
}
