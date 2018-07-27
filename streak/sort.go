package streak

import (
	"sort"
)

func sortFactors(f []Factor) {
	sort.Sort(byAfter(f))
}

// byAfter implements sort.Interface for []Factor based on
// the After field (increasing order).
type byAfter []Factor

func (a byAfter) Len() int           { return len(a) }
func (a byAfter) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byAfter) Less(i, j int) bool { return a[i].After < a[j].After }
