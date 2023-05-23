package url

// Values结构体
type Json struct {
	values map[string]interface{}
}

func NewJson() *Json {
	value := make(map[string]interface{})
	return &Json{values: value}
}
