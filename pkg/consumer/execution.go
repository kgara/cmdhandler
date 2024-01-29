package consumer

import (
	"fmt"
	"github.com/kgara/cmdhandler/pkg/shared"
	"sync"
)

type entry struct {
	key, value string
	prev, next *entry
}

type OrderedMapImpl struct {
	items      map[string]*entry
	head, tail *entry
	mu         sync.RWMutex
	// Just the encapsulation for emulating "heavy io operation"
	fileWriter FileWriter
}

func NewOrderedMap(fileWriter FileWriter) *OrderedMapImpl {
	return &OrderedMapImpl{
		items:      make(map[string]*entry),
		fileWriter: fileWriter,
	}
}

func (om *OrderedMapImpl) ExecuteCommand(cmd *shared.Command) {
	switch cmd.Action {
	case shared.AddItem:
		om.mu.Lock()
		// There was no explicit clarification on how do we handle the duplicate keys entries
		// and how we treat the order in that case, so let's just override the value and keep the initial order
		existingEntry, ok := om.items[cmd.Key]
		if ok {
			existingEntry.value = cmd.Value
		} else {
			newEntry := &entry{key: cmd.Key, value: cmd.Value}
			if om.head == nil {
				// First entry
				om.head = newEntry
			} else {
				om.tail.next = newEntry
				newEntry.prev = om.tail
			}
			om.tail = newEntry
			om.items[cmd.Key] = newEntry
		}
		om.mu.Unlock()
		if ok {
			om.fileWriter.Write(fmt.Sprintf("AddItem: Replaced item successfully. Key: %s, Value: %s\n", cmd.Key, cmd.Value))
		} else {
			om.fileWriter.Write(fmt.Sprintf("AddItem: Added item successfully. Key: %s, Value: %s\n", cmd.Key, cmd.Value))
		}
	case shared.DeleteItem:
		om.mu.Lock()
		entry, ok := om.items[cmd.Key]
		if ok {
			if entry.prev != nil {
				entry.prev.next = entry.next
			} else {
				// Head element
				om.head = entry.next
			}

			if entry.next != nil {
				entry.next.prev = entry.prev
			} else {
				// Tail element
				om.tail = entry.prev
			}

			delete(om.items, cmd.Key)
		}
		om.mu.Unlock()
		if ok {
			om.fileWriter.Write(fmt.Sprintf("DeleteItem: Deleted item successfully. Key: %s\n", cmd.Key))
		} else {
			om.fileWriter.Write(fmt.Sprintf("DeleteItem: Key %s not found\n", cmd.Key))
		}
	case shared.GetItem:
		om.mu.RLock()
		entry, ok := om.items[cmd.Key]
		om.mu.RUnlock()
		if ok {
			// Having access to the position here will compromise the O(1) complexity condition
			// forcing us to iterate over keys or keep the index of position which will have to be updated
			// on each `add` or `delete` on all entries as well.
			om.fileWriter.Write(fmt.Sprintf("GetItem: Key: %s, Value: %s\n", cmd.Key, entry.value))
		} else {
			om.fileWriter.Write(fmt.Sprintf("GetItem: Key %s not found\n", cmd.Key))
		}
	case shared.GetAllItems:
		om.mu.RLock()
		var position int
		current := om.head
		for current != nil {
			om.fileWriter.Write(fmt.Sprintf("GetAllItems: Position: %d, Key: %s, Value: %s\n", position, current.key, current.value))
			position++
			current = current.next
		}
		if position == 0 {
			om.fileWriter.Write(fmt.Sprintf("GetAllItems: Empty map\n"))
		}
		om.mu.RUnlock()
	default:
		om.fileWriter.Write(fmt.Sprintf("Action: %s, is not supported\n", cmd.Action))
	}
}

type OrderedMap interface {
	ExecuteCommand(cmd *shared.Command)
}
