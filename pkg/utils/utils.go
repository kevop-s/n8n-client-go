package utils

func RemoveEmptyInterfaces(slice []interface{}) []interface{} {
	var result []interface{}
	for _, val := range slice {
		if val != nil {
			result = append(result, val)
		}
	}
	return result
}
