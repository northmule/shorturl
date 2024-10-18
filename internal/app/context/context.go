package context

type key int

// KeyContext Ключи контекста, для передачи в запросах.
const (
	KeyContext key = iota
)

type UserUUID string
