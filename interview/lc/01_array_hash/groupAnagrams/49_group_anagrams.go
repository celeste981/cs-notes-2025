package groupAnagrams

func GroupAnagrams(strs []string) [][]string {
	res := make([][]string, 0)

	groupMap := make(map[[26]int][]string)

	for _, str := range strs {
		var tmp [26]int
		for _, v := range str {
			tmp[v-'a'] = tmp[v-'a'] + 1
		}
		groupMap[tmp] = append(groupMap[tmp], str)
	}

	for _, strings := range groupMap {
		res = append(res, strings)
	}

	return res
}
