package _8_random

func groupAnagrams(strs []string) [][]string {
	var res [][]string
	maps := make(map[[26]int][]string)
	for _, str := range strs {
		var key [26]int
		for _, ch := range str {
			key[ch-'a']++
		}
		maps[key] = append(maps[key], str)
	}
	for _, strings := range maps {
		res = append(res, strings)
	}
	return res
}
