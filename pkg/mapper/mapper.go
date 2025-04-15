package mapper

type Mapping struct {
	Key   string
	Value struct {
		Function string
		Params   []string
	}
}
