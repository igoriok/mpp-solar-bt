package input

type Input interface {
	Read() (map[string]interface{}, error)
}
