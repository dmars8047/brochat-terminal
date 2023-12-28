package state

type Property = string

const (
	UserSessionProp Property = "user-session"
)

type ApplicationState = map[Property]interface{}

func NewApplicationState() ApplicationState {
	return make(map[Property]interface{})
}

// Set creates an entry in the application state using the key and the value
func Set(m ApplicationState, key Property, value any) {
	m[key] = value
}

// Get retrieves a value from the map based on the key and the expected type.
func Get[T any](m ApplicationState, key Property) (*T, bool) {
	rawValue, exists := m[key]
	if !exists {
		return nil, false
	}

	value, ok := rawValue.(*T)

	if !ok {
		return nil, false
	}

	return value, true
}
