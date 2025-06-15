package utils

func RemoveEmptyInterfaces(slice []interface{}) []interface{} {
	var result = make([]interface{}, 0)
	for _, val := range slice {
		if val != nil {
			result = append(result, val)
		}
	}
	return result
}
