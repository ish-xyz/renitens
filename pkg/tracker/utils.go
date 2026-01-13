package tracker

func contains(source []string, target string) bool {
	contained := false
	for _, ex := range source {
		if ex == target {
			contained = true
			break
		}
	}
	return contained
}

func fmtError(err interface{}) map[string]interface{} {
	return map[string]interface{}{"error": err}
}

func fmtData(obj interface{}) map[string]interface{} {
	return map[string]interface{}{"data": obj}
}
