package account

import (
	"sort"

	"github.com/lithammer/fuzzysearch/fuzzy"
)

type userLoginMatch struct {
	userLogin ClientUserLogin
	score     int
}

type userLoginMatchList []userLoginMatch

func (list userLoginMatchList) Len() int {
	return len(list)
}

func (list userLoginMatchList) Less(i, j int) bool {
	return list[i].score < list[j].score
}

func (list userLoginMatchList) Swap(i, j int) {
	t := list[i]
	list[i] = list[j]
	list[j] = t
}

func FuzzySearchUserLogin(list []ClientUserLogin,
	query string, page, pageSize uint) (uint, []ClientUserLogin) {

	matches := make([]userLoginMatch, 0, len(list))
	for _, userLogin := range list {
		target := userLogin.Username
		score := fuzzy.RankMatchNormalizedFold(query, target)
		if score == -1 {
			continue
		}

		match := userLoginMatch{
			userLogin: userLogin,
			score:     score,
		}
		matches = append(matches, match)
	}

	sort.Sort(userLoginMatchList(matches))

	result := make([]ClientUserLogin, 0, pageSize)

	first := page * pageSize
	last := (page + 1) * pageSize
	if uint(len(matches)) < last {
		last = uint(len(matches))
	}

	for i := first; i < last; i++ {
		result = append(result, matches[i].userLogin)
	}

	return uint(len(matches)), result
}
