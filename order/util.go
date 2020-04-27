package order

import (
	"sort"

	"github.com/lithammer/fuzzysearch/fuzzy"
)

// CUSTOMER STORE
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
	storeList []CustomerStore, page, pageSize int,
	query string) (int, []CustomerStore) {

	matches := make([]storeMatch, 0)
	for _, store := range storeList {
		score := fuzzy.RankMatchNormalizedFold(query, store.Name)
		if score == -1 {
			continue
		}

		m := storeMatch{
			store: store,
			score: score,
		}
		matches = append(matches, m)
	}

	sort.Sort(storeMatchList(matches))

	result := make([]CustomerStore, 0)
	first := page * pageSize
	last := (page + 1) * pageSize
	if len(matches) < last {
		last = len(matches)
	}

	for i := first; i < last; i++ {
		result = append(result, matches[i].store)
	}
	return len(matches), result
}

// PRODUCT INFO
type productInfoMatch struct {
	product ProductInfo
	score   int
}

type productInfoMatchList []productInfoMatch

func (m productInfoMatchList) Len() int {
	return len(m)
}

func (m productInfoMatchList) Swap(i, j int) {
	tmp := m[i]
	m[i] = m[j]
	m[j] = tmp
}

func (m productInfoMatchList) Less(i, j int) bool {
	return m[i].score < m[j].score
}

func FuzzySearchProductInfo(
	products []ProductInfo, page, pageSize int,
	query string) (int, []ProductInfo) {

	matches := make([]productInfoMatch, 0)
	for _, product := range products {
		score := fuzzy.RankMatchNormalizedFold(query, product.Name)
		if score == -1 {
			continue
		}

		m := productInfoMatch{
			product: product,
			score:   score,
		}
		matches = append(matches, m)
	}

	sort.Sort(productInfoMatchList(matches))

	result := make([]ProductInfo, 0)
	first := page * pageSize
	last := (page + 1) * pageSize
	if len(matches) < last {
		last = len(matches)
	}

	for i := first; i < last; i++ {
		result = append(result, matches[i].product)
	}
	return len(matches), result
}
