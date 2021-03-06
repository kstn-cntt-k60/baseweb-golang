package product

import (
	"context"
	"sort"

	"github.com/lithammer/fuzzysearch/fuzzy"
)

type productMatch struct {
	product ClientProduct
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
	productList []ClientProduct, page, pageSize uint,
	query string) (uint, []ClientProduct) {

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

	result := make([]ClientProduct, 0)
	first := page * pageSize
	last := (page + 1) * pageSize
	if uint(len(matches)) < last {
		last = uint(len(matches))
	}

	for i := first; i < last; i++ {
		result = append(result, matches[i].product)
	}
	return uint(len(matches)), result
}

func ViewProductFuzzy(
	ctx context.Context, repo *Repo,
	page, pageSize uint,
	sortedBy, sortOrder string,
	search string) (uint, []ClientProduct, error) {

	if search == "" {
		return repo.ViewProduct(
			ctx, uint(page), uint(pageSize), sortedBy, sortOrder)
	} else {
		count, products, err := repo.SelectProduct(ctx)
		if err != nil {
			return count, products, err
		}
		count, products = FuzzySearchProduct(products,
			uint(page), uint(pageSize), search)
		return count, products, err
	}
}
