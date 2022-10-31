package schemes

import (
	"math/rand"
	"testing"
)

func TestPaginator_Offset(t *testing.T) {
	page, limit := rand.Int31(), rand.Int31()
	expectOffset := (page - 1) * limit

	p := Paginator{
		Page:  page,
		Limit: limit,
	}
	gotOffset := p.Offset()

	if expectOffset != gotOffset {
		t.Errorf("Expected %d, got %d", expectOffset, gotOffset)
	}
}
