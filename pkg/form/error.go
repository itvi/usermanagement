package form

// Define a new errors type
type errors map[string][]string

// Add method add error messages for a given field to the map key
func (e errors) Add(field, message string) {
	e[field] = append(e[field], message)
}

// Get method retrieve the first error message for a given field from
// the map.
func (e errors) Get(field string) string {
	es := e[field]
	if len(es) == 0 {
		return ""
	}
	return es[0]
}
