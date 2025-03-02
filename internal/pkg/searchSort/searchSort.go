package searchSort

import (
	"myapp/internal/app/models"
	"sort"

	"github.com/texttheater/golang-levenshtein/levenshtein"
)

func SearchSort(target string, searchResponse []interface{}) []interface{} {
	sort.Slice(searchResponse, func(i, j int) bool {
		var value1, value2 string

		if userCfg, ok := searchResponse[i].(models.UserCfg); ok {
			value1 = userCfg.User
		} else if post, ok := searchResponse[i].(models.Post); ok {
			value1 = post.Header
		}

		if userCfg, ok := searchResponse[j].(models.UserCfg); ok {
			value2 = userCfg.User
		} else if post, ok := searchResponse[i].(models.Post); ok {
			value2 = post.Header
		}

		return levenshtein.DistanceForStrings([]rune(target), []rune(value1), levenshtein.DefaultOptions) <
			levenshtein.DistanceForStrings([]rune(target), []rune(value2), levenshtein.DefaultOptions)
	})
	return searchResponse
}
