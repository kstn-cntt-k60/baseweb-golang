package facility

import (
	"sort"

	"github.com/lithammer/fuzzysearch/fuzzy"
)

type warehouseMatch struct {
	warehouse Warehouse
	score     int
}

type warehouseMatchList []warehouseMatch

func (m warehouseMatchList) Len() int {
	return len(m)
}

func (m warehouseMatchList) Swap(i, j int) {
	tmp := m[i]
	m[i] = m[j]
	m[j] = tmp
}

func (m warehouseMatchList) Less(i, j int) bool {
	return m[i].score < m[j].score
}

func FuzzySearchWarehouse(
	warehouseList []Warehouse, page, pageSize uint,
	query string) (uint, []Warehouse) {

	matches := make([]warehouseMatch, 0)
	for _, w := range warehouseList {
		score := fuzzy.RankMatchNormalizedFold(query, w.Name)
		if score == -1 {
			continue
		}

		m := warehouseMatch{
			warehouse: w,
			score:     score,
		}
		matches = append(matches, m)
	}

	sort.Sort(warehouseMatchList(matches))

	result := make([]Warehouse, 0)
	first := page * pageSize
	last := (page + 1) * pageSize
	count := uint(len(matches))
	if count < last {
		last = count
	}

	for i := first; i < last; i++ {
		result = append(result, matches[i].warehouse)
	}
	return count, result
}
