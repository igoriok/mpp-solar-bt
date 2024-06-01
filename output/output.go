package output

type Output interface {
	Write(data map[string]interface{}) error
}
