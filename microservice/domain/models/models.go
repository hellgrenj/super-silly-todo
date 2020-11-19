package models

// Item is the struct for an item in a todo list
type Item struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Done bool   `json:"done"`
}

// OK is the validation function for the struct Item
func (i *Item) OK() error {
	if len(i.Name) == 0 {
		return ErrMissingField("Name")
	}
	return nil
}
