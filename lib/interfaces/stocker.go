package interfaces

// Stocker must be implemented to be used as Stock APIs
type Stocker interface {
	GetCode() string
	String() string
}
