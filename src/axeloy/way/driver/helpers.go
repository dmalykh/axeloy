package driver

func GetValue(params Params, name string) (string, bool) {
	if param, exists := params[name]; exists && len(param) > 0 {
		return param[0], true
	}
	return ``, false
}

func GetValues(params Params, name string) ([]string, bool) {
	param, exists := params[name]
	return param, exists
}
