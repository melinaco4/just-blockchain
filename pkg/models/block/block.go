package block

type Block struct {
	Position  int
	Data      SongCheckout
	TimeStamp string
	Hash      string
	PrevHash  string
}
