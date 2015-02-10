package api

import (
	"fmt"
	"testing"
)

func TestNewUUID(t *testing.T) {
	uuid, err := NewUUID()
	if err != nil {
		t.Fatalf("error: %v\n", err)
	}
	fmt.Printf("%s\n", uuid)
}
