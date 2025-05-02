package common

type Enum []string

func (s Enum) Contains(str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}
