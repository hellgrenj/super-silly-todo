package features

// Flags is the struct for our feature flags
type Flags struct {
	DelegateAddListItemToMicroservice bool
}

// GetFlags returns a struct with our feature flags
func GetFlags() *Flags {
	flags := &Flags{DelegateAddListItemToMicroservice: true}
	return flags
}
