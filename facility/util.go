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

type storeMatch struct {
	store CustomerStore
	score int
}

type storeMatchList []storeMatch

func (m storeMatchList) Len() int {
	return len(m)
}

func (m storeMatchList) Swap(i, j int) {
	tmp := m[i]
	m[i] = m[j]
	m[j] = tmp
}

func (m storeMatchList) Less(i, j int) bool {
	return m[i].score < m[j].score
}

func FuzzySearchCustomerStore(
	storeList []CustomerStore, page, pageSize uint,
	query string) (uint, []CustomerStore) {

	matches := make([]storeMatch, 0)
	for _, s := range storeList {
		score := fuzzy.RankMatchNormalizedFold(query, s.Name)
		if score == -1 {
			continue
		}

		m := storeMatch{
			store: s,
			score: score,
		}
		matches = append(matches, m)
	}

	sort.Sort(storeMatchList(matches))

	result := make([]CustomerStore, 0)
	first := page * pageSize
	last := (page + 1) * pageSize
	count := uint(len(matches))
	if count < last {
		last = count
	}

	for i := first; i < last; i++ {
		result = append(result, matches[i].store)
	}
	return count, result
}

type simpleCustomerMatch struct {
	simpleCustomer SimpleCustomer
	score          int
}

type simpleCustomerMatchList []simpleCustomerMatch

func (m simpleCustomerMatchList) Len() int {
	return len(m)
}

func (m simpleCustomerMatchList) Swap(i, j int) {
	tmp := m[i]
	m[i] = m[j]
	m[j] = tmp
}

func (m simpleCustomerMatchList) Less(i, j int) bool {
	return m[i].score < m[j].score
}

func FuzzySearchSimpleCustomer(
	simpleCustomerList []SimpleCustomer,
	query string) []SimpleCustomer {

	matches := make([]simpleCustomerMatch, 0)
	for _, w := range simpleCustomerList {
		score := fuzzy.RankMatchNormalizedFold(query, w.Name)
		if score == -1 {
			continue
		}

		m := simpleCustomerMatch{
			simpleCustomer: w,
			score:          score,
		}
		matches = append(matches, m)
	}

	sort.Sort(simpleCustomerMatchList(matches))

	result := make([]SimpleCustomer, 0)
	first := 0
	last := 10
	count := len(matches)
	if count < last {
		last = count
	}

	for i := first; i < last; i++ {
		result = append(result, matches[i].simpleCustomer)
	}
	return result
}
