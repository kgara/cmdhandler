package consumer

import (
	rand "github.com/dchest/uniuri"
	"github.com/kgara/cmdhandler/pkg/shared"
	"github.com/stretchr/testify/mock"
	"testing"
)

type FileWriterMock struct {
	mock.Mock
}

func (m *FileWriterMock) Write(content string) {
	m.Called(content)
}

func initialize() (*FileWriterMock, *OrderedMapImpl) {
	fileWriterMock := &FileWriterMock{}
	om := NewOrderedMap(fileWriterMock)
	return fileWriterMock, om
}

func TestExecuteCommandAddItemSingle(t *testing.T) {
	fileWriterMock, om := initialize()

	cmd := &shared.Command{
		Action: shared.AddItem,
		Key:    "testKey",
		Value:  "testValue",
	}

	fileWriterMock.On("Write", "AddItem: Added item successfully. Key: testKey, Value: testValue\n").Once()
	om.ExecuteCommand(cmd)

	getItemCmd := &shared.Command{
		Action: shared.GetItem,
		Key:    "testKey",
	}
	fileWriterMock.On("Write", "GetItem: Key: testKey, Value: testValue\n").Times(1)
	om.ExecuteCommand(getItemCmd)

	getAllItemsCmd := &shared.Command{
		Action: shared.GetAllItems,
	}

	fileWriterMock.On("Write", "GetAllItems: Position: 0, Key: testKey, Value: testValue\n").Once()
	om.ExecuteCommand(getAllItemsCmd)

	fileWriterMock.AssertExpectations(t)
}

func TestExecuteCommandAddItemSomeEntriesBeforeAndAfter(t *testing.T) {
	fileWriterMock, om := initialize()

	cmd := &shared.Command{
		Action: shared.AddItem,
		Key:    rand.New(),
		Value:  rand.New(),
	}
	fileWriterMock.On("Write", mock.Anything).Once()
	om.ExecuteCommand(cmd)

	cmd = &shared.Command{
		Action: shared.AddItem,
		Key:    "testKey",
		Value:  "testValue",
	}
	fileWriterMock.On("Write", "AddItem: Added item successfully. Key: testKey, Value: testValue\n").Once()
	om.ExecuteCommand(cmd)

	cmd = &shared.Command{
		Action: shared.AddItem,
		Key:    rand.New(),
		Value:  rand.New(),
	}
	fileWriterMock.On("Write", mock.Anything).Once()
	om.ExecuteCommand(cmd)

	getItemCmd := &shared.Command{
		Action: shared.GetItem,
		Key:    "testKey",
	}
	fileWriterMock.On("Write", "GetItem: Key: testKey, Value: testValue\n").Times(1)
	om.ExecuteCommand(getItemCmd)

	getAllItemsCmd := &shared.Command{
		Action: shared.GetAllItems,
	}

	fileWriterMock.On("Write", mock.Anything).Once()
	fileWriterMock.On("Write", "GetAllItems: Position: 1, Key: testKey, Value: testValue\n").Once()
	fileWriterMock.On("Write", mock.Anything).Once()
	om.ExecuteCommand(getAllItemsCmd)

	fileWriterMock.AssertExpectations(t)
}

func TestExecuteCommandAddItemReplace(t *testing.T) {
	fileWriterMock, om := initialize()

	cmd := &shared.Command{
		Action: shared.AddItem,
		Key:    rand.New(),
		Value:  rand.New(),
	}
	fileWriterMock.On("Write", mock.Anything).Once()
	om.ExecuteCommand(cmd)

	// Add our item
	cmd = &shared.Command{
		Action: shared.AddItem,
		Key:    "testKey",
		Value:  "testValue",
	}
	fileWriterMock.On("Write", "AddItem: Added item successfully. Key: testKey, Value: testValue\n").Once()
	om.ExecuteCommand(cmd)

	cmd = &shared.Command{
		Action: shared.AddItem,
		Key:    rand.New(),
		Value:  rand.New(),
	}
	fileWriterMock.On("Write", mock.Anything).Once()
	om.ExecuteCommand(cmd)

	// Replace the existing item
	cmd = &shared.Command{
		Action: shared.AddItem,
		Key:    "testKey",
		Value:  "differentTestValue",
	}
	fileWriterMock.On("Write", "AddItem: Replaced item successfully. Key: testKey, Value: differentTestValue\n").Once()
	om.ExecuteCommand(cmd)

	getItemCmd := &shared.Command{
		Action: shared.GetItem,
		Key:    "testKey",
	}
	fileWriterMock.On("Write", "GetItem: Key: testKey, Value: differentTestValue\n").Times(1)
	om.ExecuteCommand(getItemCmd)

	getAllItemsCmd := &shared.Command{
		Action: shared.GetAllItems,
	}

	fileWriterMock.On("Write", mock.Anything).Once()
	fileWriterMock.On("Write", "GetAllItems: Position: 1, Key: testKey, Value: differentTestValue\n").Once()
	fileWriterMock.On("Write", mock.Anything).Once()
	om.ExecuteCommand(getAllItemsCmd)

	fileWriterMock.AssertExpectations(t)
}

func TestExecuteCommandDeleteItemSingle(t *testing.T) {
	fileWriterMock, om := initialize()

	addItemCmd := &shared.Command{
		Action: shared.AddItem,
		Key:    "testKey",
		Value:  "testValue",
	}
	fileWriterMock.On("Write", mock.Anything).Times(1)
	om.ExecuteCommand(addItemCmd)

	deleteItemCmd := &shared.Command{
		Action: shared.DeleteItem,
		Key:    "testKey",
	}
	fileWriterMock.On("Write", "DeleteItem: Deleted item successfully. Key: testKey\n").Times(1)
	om.ExecuteCommand(deleteItemCmd)

	getItemCmd := &shared.Command{
		Action: shared.GetItem,
		Key:    "testKey",
	}
	fileWriterMock.On("Write", "GetItem: Key testKey not found\n").Times(1)
	om.ExecuteCommand(getItemCmd)

	getAllItemsCmd := &shared.Command{
		Action: shared.GetAllItems,
	}

	fileWriterMock.On("Write", "GetAllItems: Empty map\n").Times(1)
	om.ExecuteCommand(getAllItemsCmd)

	fileWriterMock.AssertExpectations(t)
}

func TestExecuteCommandDeleteItemMiddle(t *testing.T) {
	fileWriterMock, om := initialize()

	cmd := &shared.Command{
		Action: shared.AddItem,
		Key:    rand.New(),
		Value:  rand.New(),
	}
	fileWriterMock.On("Write", mock.Anything).Once()
	om.ExecuteCommand(cmd)

	// Add our item
	cmd = &shared.Command{
		Action: shared.AddItem,
		Key:    "testKey",
		Value:  "testValue",
	}
	fileWriterMock.On("Write", "AddItem: Added item successfully. Key: testKey, Value: testValue\n").Once()
	om.ExecuteCommand(cmd)

	cmd = &shared.Command{
		Action: shared.AddItem,
		Key:    rand.New(),
		Value:  rand.New(),
	}
	fileWriterMock.On("Write", mock.Anything).Once()
	om.ExecuteCommand(cmd)

	// Delete item
	cmd = &shared.Command{
		Action: shared.DeleteItem,
		Key:    "testKey",
	}
	fileWriterMock.On("Write", "DeleteItem: Deleted item successfully. Key: testKey\n").Once()
	om.ExecuteCommand(cmd)

	getItemCmd := &shared.Command{
		Action: shared.GetItem,
		Key:    "testKey",
	}
	fileWriterMock.On("Write", "GetItem: Key testKey not found\n").Times(1)
	om.ExecuteCommand(getItemCmd)

	getAllItemsCmd := &shared.Command{
		Action: shared.GetAllItems,
	}

	fileWriterMock.On("Write", mock.Anything).Once()
	fileWriterMock.On("Write", mock.Anything).Once()
	om.ExecuteCommand(getAllItemsCmd)

	fileWriterMock.AssertExpectations(t)
}

func TestExecuteCommandDeleteItemHead(t *testing.T) {
	fileWriterMock, om := initialize()

	// Add our item
	cmd := &shared.Command{
		Action: shared.AddItem,
		Key:    "testKey",
		Value:  "testValue",
	}
	fileWriterMock.On("Write", "AddItem: Added item successfully. Key: testKey, Value: testValue\n").Once()
	om.ExecuteCommand(cmd)

	cmd = &shared.Command{
		Action: shared.AddItem,
		Key:    rand.New(),
		Value:  rand.New(),
	}
	fileWriterMock.On("Write", mock.Anything).Once()
	om.ExecuteCommand(cmd)

	// Delete item
	cmd = &shared.Command{
		Action: shared.DeleteItem,
		Key:    "testKey",
	}
	fileWriterMock.On("Write", "DeleteItem: Deleted item successfully. Key: testKey\n").Once()
	om.ExecuteCommand(cmd)

	getItemCmd := &shared.Command{
		Action: shared.GetItem,
		Key:    "testKey",
	}
	fileWriterMock.On("Write", "GetItem: Key testKey not found\n").Times(1)
	om.ExecuteCommand(getItemCmd)

	getAllItemsCmd := &shared.Command{
		Action: shared.GetAllItems,
	}

	fileWriterMock.On("Write", mock.Anything).Once()
	om.ExecuteCommand(getAllItemsCmd)

	// Just adding one more item to be sure, that we did not break the map somehow
	cmd = &shared.Command{
		Action: shared.AddItem,
		Key:    rand.New(),
		Value:  rand.New(),
	}
	fileWriterMock.On("Write", mock.Anything).Once()
	om.ExecuteCommand(cmd)

	fileWriterMock.On("Write", mock.Anything).Once()
	fileWriterMock.On("Write", mock.Anything).Once()
	om.ExecuteCommand(getAllItemsCmd)

	fileWriterMock.AssertExpectations(t)
}

func TestExecuteCommandDeleteItemTail(t *testing.T) {
	fileWriterMock, om := initialize()

	cmd := &shared.Command{
		Action: shared.AddItem,
		Key:    rand.New(),
		Value:  rand.New(),
	}
	fileWriterMock.On("Write", mock.Anything).Once()
	om.ExecuteCommand(cmd)

	// Add our item
	cmd = &shared.Command{
		Action: shared.AddItem,
		Key:    "testKey",
		Value:  "testValue",
	}
	fileWriterMock.On("Write", "AddItem: Added item successfully. Key: testKey, Value: testValue\n").Once()
	om.ExecuteCommand(cmd)

	// Delete item
	cmd = &shared.Command{
		Action: shared.DeleteItem,
		Key:    "testKey",
	}
	fileWriterMock.On("Write", "DeleteItem: Deleted item successfully. Key: testKey\n").Once()
	om.ExecuteCommand(cmd)

	getItemCmd := &shared.Command{
		Action: shared.GetItem,
		Key:    "testKey",
	}
	fileWriterMock.On("Write", "GetItem: Key testKey not found\n").Times(1)
	om.ExecuteCommand(getItemCmd)

	getAllItemsCmd := &shared.Command{
		Action: shared.GetAllItems,
	}

	fileWriterMock.On("Write", mock.Anything).Once()
	om.ExecuteCommand(getAllItemsCmd)

	// Just adding one more item to be sure, that we did not break the map somehow
	cmd = &shared.Command{
		Action: shared.AddItem,
		Key:    rand.New(),
		Value:  rand.New(),
	}
	fileWriterMock.On("Write", mock.Anything).Once()
	om.ExecuteCommand(cmd)

	fileWriterMock.On("Write", mock.Anything).Once()
	fileWriterMock.On("Write", mock.Anything).Once()
	om.ExecuteCommand(getAllItemsCmd)

	fileWriterMock.AssertExpectations(t)
}

func TestExecuteCommandGetItem(t *testing.T) {
	fileWriterMock, om := initialize()

	addItemCmd := &shared.Command{
		Action: shared.AddItem,
		Key:    "testKey",
		Value:  "testValue",
	}
	fileWriterMock.On("Write", mock.Anything).Times(1)
	om.ExecuteCommand(addItemCmd)

	getItemCmd := &shared.Command{
		Action: shared.GetItem,
		Key:    "testKey",
	}
	fileWriterMock.On("Write", "GetItem: Key: testKey, Value: testValue\n").Times(1)
	om.ExecuteCommand(getItemCmd)

	fileWriterMock.AssertExpectations(t)
}

func TestExecuteCommandGetItemDoesNotExist(t *testing.T) {
	fileWriterMock, om := initialize()

	getItemCmd := &shared.Command{
		Action: shared.GetItem,
		Key:    "testKey",
	}
	fileWriterMock.On("Write", "GetItem: Key testKey not found\n").Times(1)
	om.ExecuteCommand(getItemCmd)

	addItemCmd := &shared.Command{
		Action: shared.AddItem,
		Key:    "testKey",
		Value:  "testValue",
	}
	fileWriterMock.On("Write", mock.Anything).Times(1)
	om.ExecuteCommand(addItemCmd)

	//Different key
	getItemCmd = &shared.Command{
		Action: shared.GetItem,
		Key:    "wrongKey",
	}
	fileWriterMock.On("Write", "GetItem: Key wrongKey not found\n").Times(1)
	om.ExecuteCommand(getItemCmd)

	fileWriterMock.AssertExpectations(t)
}

func TestExecuteCommandGetAllItems(t *testing.T) {
	fileWriterMock, om := initialize()

	getAllItemsCmd := &shared.Command{
		Action: shared.GetAllItems,
	}

	// Nothing on the first call
	fileWriterMock.On("Write", "GetAllItems: Empty map\n").Times(1)
	om.ExecuteCommand(getAllItemsCmd)

	addItemCmd1 := &shared.Command{
		Action: shared.AddItem,
		Key:    "key1",
		Value:  "value1",
	}
	fileWriterMock.On("Write", mock.Anything).Times(1)
	om.ExecuteCommand(addItemCmd1)

	addItemCmd2 := &shared.Command{
		Action: shared.AddItem,
		Key:    "key2",
		Value:  "value2",
	}
	fileWriterMock.On("Write", mock.Anything).Times(1)
	om.ExecuteCommand(addItemCmd2)

	fileWriterMock.On("Write", "GetAllItems: Position: 0, Key: key1, Value: value1\n").Once()
	fileWriterMock.On("Write", "GetAllItems: Position: 1, Key: key2, Value: value2\n").Once()
	om.ExecuteCommand(getAllItemsCmd)

	fileWriterMock.AssertExpectations(t)
}

func TestExecuteCommandNotSupported(t *testing.T) {
	fileWriterMock, om := initialize()

	cmd := &shared.Command{
		Action: 13,
		Key:    rand.New(),
		Value:  rand.New(),
	}
	fileWriterMock.On("Write", "Action: unknownAction: 13, is not supported\n").Once()
	om.ExecuteCommand(cmd)
}
