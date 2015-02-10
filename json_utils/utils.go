package json_utils

func GetByPath(jsonObj interface{}, pathFragments ...interface{}) interface{} {

	obj := jsonObj

	for _, key := range pathFragments {
		switch key.(type) {
		case int:
			array, ok := obj.([]interface{})
			i := key.(int)
			if ok && i >= 0 && i < len(array) {
				obj = array[i]
			} else {
				return nil
			}
			break
		case string:
			jsonData, ok := obj.(map[string]interface{})
			if ok {
				obj = jsonData[key.(string)]
			} else {
				return nil
			}
			break
		default:
			//Unsupported type
			return nil
		}

	}

	return obj
}

func SetByPath(jsonObj interface{}, value interface{}, pathFragments ...interface{}) bool {

	result := false

	size := len(pathFragments)

	parent := jsonObj

	if size == 0 {
		return false
	} else if size > 1 {
		parent = GetByPath(jsonObj, pathFragments[:size-1]...)
	}
	if parent != nil {
		last := size - 1
		key := pathFragments[last]
		switch key.(type) {
		case int:
			array, ok := parent.([]interface{})
			i := key.(int)
			if ok && i >= 0 && i < len(array) {
				array[key.(int)] = value
				result = true
			} else {
				result = false
			}
			break
		case string:
			jsonData, ok := parent.(map[string]interface{})
			if ok {
				jsonData[key.(string)] = value
				result = true
			} else {
				result = false
			}
			break
		default:
			//Unsupported type
		}
	}
	return result
}
