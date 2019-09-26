package dialogs

import (
	"testing"
)

func TestStorageSet(t *testing.T) {
	var id = "test_id"
	var memStorage = NewMemoryStorage()

	testState := memStorage.GetState(id)
	if testState != 0 {
		t.Error("Test state start: Expected 0, got ", testState)
	}

	memStorage.SetState(id, 1)
	testState = memStorage.GetState(id)
	if testState != 1 {
		t.Error("Test state : Expected 1, got ", testState)
	}

	memStorage.SetData(id, 1)
	testData := memStorage.GetData(id)
	if testData != 1 {
		t.Error("Test data :Expected data 1, got ", testData)
	}

}
