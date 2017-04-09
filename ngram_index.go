package suggest

/*
 * inspired by
 * http://www.aaai.org/ocs/index.php/AAAI/AAAI10/paper/viewFile/1939/2234
 * http://nlp.stanford.edu/IR-book/html/htmledition/k-gram-indexes-for-wildcard-queries-1.html
 * http://bazhenov.me/blog/2012/08/04/autocomplete.html
 */

import (
	"container/heap"
)

type invertedListsT map[string][]int

type NGramIndex struct {
	k             int
	invertedLists invertedListsT
	dictionary    []string
	cardinalities []int
	index         int
}

func NewNGramIndex(k int) *NGramIndex {
	if k < 2 || k > 4 {
		panic("k should be in [2, 4]")
	}

	return &NGramIndex{
		k, make(invertedListsT), make([]string, 0), make([]int, 0), 0,
	}
}

// Add given word to invertedList
func (self *NGramIndex) AddWord(word string) {
	prepared := prepareString(word)
	ngrams := self.getNGramSet(prepared)
	cardinality := len(ngrams)
	for _, ngram := range ngrams {
		self.invertedLists[ngram] = append(self.invertedLists[ngram], self.index)
	}

	self.dictionary = append(self.dictionary, word)
	self.cardinalities = append(self.cardinalities, cardinality)
	self.index++
}

// Return top-k similar strings
func (self *NGramIndex) Suggest(word string, topK int) []string {
	candidates := self.search(word, topK)
	result := make([]string, 0, topK)
	for candidates.Len() > 0 {
		r := heap.Pop(candidates).(*rank)
		result = append([]string{self.dictionary[r.id]}, result...)
	}

	return result
}

//1. try to receive corresponding inverted list for word's ngrams
//2. calculate distance between current word and candidates
//3. return rankHeap
func (self *NGramIndex) search(word string, topK int) *rankHeap {
	preparedWord := prepareString(word)
	set := self.getNGramSet(preparedWord)
	lenA := len(set)

	// find max word id for memory optimize
	maxId := 0
	for _, ngram := range set {
		list := self.invertedLists[ngram]
		if len(list) > 0 {
			curMaxId := list[len(list)-1]
			if curMaxId > maxId {
				maxId = curMaxId
			}
		}
	}

	// calc intersections between current word's ngram set and candidates */
	counts := make([]int, maxId+1)
	for _, ngram := range set {
		for _, id := range self.invertedLists[ngram] {
			counts[id]++
		}
	}

	// use heap search for finding top k items in a list efficiently
	// see http://stevehanov.ca/blog/index.php?id=122
	h := &rankHeap{}
	for id, inter := range counts {
		if inter == 0 {
			continue
		}

		lenB := self.cardinalities[id]
		// use jaccard distance as metric for calc words similarity
		// 1 - |intersection| / |union| = 1 - |intersection| / (|A| + |B| - |intersection|)
		distance := 1 - float64(inter)/float64(lenA+lenB-inter)
		if h.Len() < topK || h.Min().(*rank).distance > distance {
			if h.Len() == topK {
				heap.Pop(h)
			}

			heap.Push(h, &rank{id, distance})
		}
	}

	return h
}

// Return unique ngrams with frequency
func (self *NGramIndex) getNGramSet(word string) []string {
	return GetNGramSet(word, self.k)
}
