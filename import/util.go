package importProduct

import (
	"sort"

	"github.com/lithammer/fuzzysearch/fuzzy"
)

// PRODUCT
type productMatch struct {
	product Product
	score   int
}

type productMatchList []productMatch

func (m productMatchList) Len() int {
	return len(m)
}

func (m productMatchList) Swap(i, j int) {
	tmp := m[i]
	m[i] = m[j]
	m[j] = tmp
}

func (m productMatchList) Less(i, j int) bool {
	return m[i].score < m[j].score
}

func FuzzySearchProduct(
	productList []Product, page, pageSize int,
	query string) (uint, []Product) {

	matches := make([]productMatch, 0)
	for _, product := range productList {
		score := fuzzy.RankMatchNormalizedFold(query, product.Name)
		if score == -1 {
			continue
		}

		m := productMatch{
			product: product,
			score:   score,
		}
		matches = append(matches, m)
	}

	sort.Sort(productMatchList(matches))

	result := make([]Product, 0)
	first := page * pageSize
	last := (page + 1) * pageSize
	if len(matches) < last {
		last = len(matches)
	}

	for i := first; i < last; i++ {
		result = append(result, matches[i].product)
	}
	return uint(len(matches)), result
}

