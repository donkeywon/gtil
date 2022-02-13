package statistics

import (
	"fmt"
	"testing"
)

func TestStatistics_Incr(t *testing.T) {
	s := New("abc")

	s.Incr("abc", 111)

	fmt.Println(s.Export())
}
