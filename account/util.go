package account

import (
	"fmt"
	"log"
	"sort"
	"time"

	"github.com/lithammer/fuzzysearch/fuzzy"
)

type match struct {
	person SimplePerson
	score  int
}

type matchList []match

func (m matchList) Len() int {
	return len(m)
}

func (m matchList) Swap(i, j int) {
	a := m[i]
	m[i] = m[j]
	m[j] = a
}

func (m matchList) Less(i, j int) bool {
	return m[i].score < m[j].score
}

func FuzzySearchSimplePerson(
	personList []SimplePerson,
	query string) []SimplePerson {

	matches := make([]match, 0)
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

		matches = append(matches, match{
			person: person,
			score:  score,
		})
	}

	sort.Sort(matchList(matches))

	result := make([]SimplePerson, 0, 10)
	for i, m := range matches {
		if i >= 10 {
			break
		}
		result = append(result, m.person)
	}
	return result
}
