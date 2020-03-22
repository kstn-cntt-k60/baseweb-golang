package account

import (
	"fmt"
	"log"
	"sort"
	"time"

	"github.com/lithammer/fuzzysearch/fuzzy"
)

type simplePersonMatch struct {
	person SimplePerson
	score  int
}

type simplePersonMatchList []simplePersonMatch

func (m simplePersonMatchList) Len() int {
	return len(m)
}

func (m simplePersonMatchList) Swap(i, j int) {
	a := m[i]
	m[i] = m[j]
	m[j] = a
}

func (m simplePersonMatchList) Less(i, j int) bool {
	return m[i].score < m[j].score
}

func FuzzySearchSimplePerson(
	personList []SimplePerson,
	query string) []SimplePerson {

	matches := make([]simplePersonMatch, 0)
	for _, person := range personList {
		d, err := time.Parse(time.RFC3339, person.BirthDate)
		if err != nil {
			log.Panic(err)
		}

		target := fmt.Sprint(person.LastName, " ",
			person.MiddleName, " ", person.FirstName, " ",
			d.Day(), "/", int(d.Month()), "/", d.Year())

		score := fuzzy.RankMatchNormalizedFold(query, target)
		if score == -1 {
			continue
		}

		matches = append(matches, simplePersonMatch{
			person: person,
			score:  score,
		})
	}

	sort.Sort(simplePersonMatchList(matches))

	result := make([]SimplePerson, 0, 10)
	for i, m := range matches {
		if i >= 10 {
			break
		}
		result = append(result, m.person)
	}
	return result
}

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
