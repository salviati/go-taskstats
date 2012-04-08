package main

func in(needle string, haystack []string) bool {
	for _, i := range haystack {
		if needle == i { return true }
	}
	return false
}


func seq(s1, s2 []string) bool {
	if len(s1) != len(s2) {
		return false
	}
	
	if len(s1) == 0 {
		return true
	}
	
	for i:=0; i<len(s1); i++ {
		if s1[i] != s2[i] {
			return false
		}
	}
	return true
}

