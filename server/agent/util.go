package agent

func stringListToMapKeys(list []string) map[string]bool {
	ret := make(map[string]bool)
	for _, item := range list {
		ret[item] = true
	}
	return ret
}
