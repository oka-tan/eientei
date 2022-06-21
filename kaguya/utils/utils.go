package utils

func ToBool(i *int8) bool {
	return i != nil && *i != 0
}
