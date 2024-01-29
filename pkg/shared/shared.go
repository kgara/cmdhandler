package shared

import "fmt"

// Command Our client-server command
type Command struct {
	Action ActionType
	Key    string
	Value  string
}

type ActionType int

const (
	AddItem ActionType = iota
	DeleteItem
	GetItem
	GetAllItems
)

func (a ActionType) String() string {
	switch a {
	case AddItem:
		return "addItem"
	case DeleteItem:
		return "deleteItem"
	case GetItem:
		return "getItem"
	case GetAllItems:
		return "getAllItems"
	default:
		return fmt.Sprintf("unknownAction: %d", a)
	}
}
