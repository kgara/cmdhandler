package shared

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCommandSerialization(t *testing.T) {
	command := &Command{Action: AddItem, Key: "key1", Value: "value1"}
	commandMarshaled, err := json.Marshal(command)
	if err != nil {
		t.Errorf("Error decoding JSON: %s\n", err)
	}
	assert.Equal(t, `{"Action":0,"Key":"key1","Value":"value1"}`, string(commandMarshaled))
}

func TestCommandDeserialization(t *testing.T) {
	expectedCommand := &Command{Action: AddItem, Key: "key1", Value: "value1"}
	jsonCommand := `{"Action":0,"Key":"key1","Value":"value1"}`
	actualCommand := &Command{}
	err := json.Unmarshal([]byte(jsonCommand), actualCommand)
	if err != nil {
		t.Errorf("Error decoding JSON: %s\n", err)
	}
	assert.Equal(t, expectedCommand, actualCommand)
}

func TestCommandSerializationNilValue(t *testing.T) {
	command := &Command{Action: GetItem, Key: "key1"}
	commandMarshaled, err := json.Marshal(command)
	if err != nil {
		t.Errorf("Error decoding JSON: %s\n", err)
	}
	assert.Equal(t, `{"Action":2,"Key":"key1","Value":""}`, string(commandMarshaled))
}

func TestCommandDeserializationNilValue(t *testing.T) {
	expectedCommand := &Command{Action: GetItem, Key: "key1"}
	jsonCommand := `{"Action":2,"Key":"key1","Value":""}`
	actualCommand := &Command{}
	err := json.Unmarshal([]byte(jsonCommand), actualCommand)
	if err != nil {
		t.Errorf("Error decoding JSON: %s\n", err)
	}
	assert.Equal(t, expectedCommand, actualCommand)
}
