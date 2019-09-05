package merger

import "github.com/alldroll/suggest/pkg/utils"

// ListMerger solves `threshold`-occurrence problem:
// For given inverted lists find the set of strings ids, that appears at least
// `threshold` times.
type ListMerger interface {
	// Merge returns list of candidates, that appears at least `threshold` times.
	Merge(rid Rid, threshold int, collector Collector) error
}

// Rid represents inverted lists for ListMerger
type Rid []ListIterator

// Len is the number of elements in the collection.
func (p Rid) Len() int { return len(p) }

// Less reports whether the element with
// index i should sort before the element with index j.
func (p Rid) Less(i, j int) bool { return p[i].Len() < p[j].Len() }

// Swap swaps the elements with indexes i and j.
func (p Rid) Swap(i, j int) { p[i], p[j] = p[j], p[i] }

// MergeCandidate is result of merging Rid
type MergeCandidate uint64

// NewMergeCandidate creates a new instance of MergeCandidate
func NewMergeCandidate(position uint32, overlap int) MergeCandidate {
	return MergeCandidate(utils.Pack(position, uint32(overlap)))
}

// Position returns the given position of the candidate
func (m MergeCandidate) Position() uint32 {
	position, _ := utils.Unpack(uint64(m))
	return position
}

// Overlap returns the current overlap count of the candidate
func (m MergeCandidate) Overlap() int {
	_, overlap := utils.Unpack(uint64(m))
	return int(overlap)
}

// increment increments the overlap value of the candidate
func (m *MergeCandidate) increment() {
	*m = NewMergeCandidate(m.Position(), m.Overlap()+1)
}

// mergerOptimizer internal merger that is aimed to optimize merge workflow
type mergerOptimizer struct {
	merger      ListMerger
	intersector ListIntersector
}

func newMerger(merger ListMerger) ListMerger {
	return &mergerOptimizer{
		merger:      merger,
		intersector: Intersector(),
	}
}

// Merge returns list of candidates, that appears at least `threshold` times.
func (m *mergerOptimizer) Merge(rid Rid, threshold int, collector Collector) error {
	n := len(rid)

	if n < threshold || n == 0 {
		return nil
	}

	if n == threshold {
		return m.intersector.Intersect(rid, collector)
	}

	return m.merger.Merge(rid, threshold, collector)
}