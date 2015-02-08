package json

import (
	"testing"
)

func TestGetByPath(t *testing.T) {

	jsonTest := map[string]interface{}{
		"x": true,
		"y": "y",
		"z": 1,
	}

	jsonArray := []interface{}{0, true, "o", jsonTest}

	jsonObject := map[string]interface{}{
		"a": "a",
		"b": map[string]interface{}{
			"x": true,
			"y": 2,
		},
		"c": jsonArray,
	}

	if GetByPath(jsonObject, "a") != "a" {
		t.Fatal("Failed to get jsonObject.a")
	}
	if GetByPath(jsonObject, "b", "x") != true {
		t.Fatal("Failed to get jsonObject.b.x")
	}
	if GetByPath(jsonObject, "b", "y") != 2 {
		t.Fatal("Failed to get jsonObject.b.y")
	}
	if GetByPath(jsonObject, "c", 0) != 0 {
		t.Fatal("Failed to get jsonObject.c[0]")
	}
	if GetByPath(jsonObject, "c", 1) != true {
		t.Fatal("Failed to get jsonObject.c[1]")
	}
	if GetByPath(jsonObject, "c", 2) != "o" {
		t.Fatal("Failed to get jsonObject.c[2]")
	}
	if GetByPath(jsonArray, 0) != 0 {
		t.Fatal("Failed to get jsonArray[0]")
	}
	if GetByPath(jsonArray, 3, "x") != true {
		t.Fatal("Failed to get jsonArray[3].x")
	}
	if GetByPath(jsonArray, 4, "x") != nil {
		t.Fatal("Failed to get jsonArray[4].x")
	}
	if GetByPath(jsonObject, "d", 1) != nil {
		t.Fatal("Failed to get jsonObject.d[1]")
	}
}

func TestSetByPath(t *testing.T) {
	jsonTest := map[string]interface{}{
		"x": true,
		"y": "y",
		"z": 1,
	}

	jsonArray := []interface{}{0, true, "o", jsonTest}

	jsonObject := map[string]interface{}{
		"a": "a",
		"b": map[string]interface{}{
			"x": true,
			"y": 2,
		},
		"c": jsonArray,
	}

	if !SetByPath(jsonObject, 1, "a") || GetByPath(jsonObject, "a") != 1 {
		t.Fatal("Failed to set jsonObject.a = 1")
	}

	if !SetByPath(jsonObject, 1, "b", "x") || GetByPath(jsonObject, "b", "x") != 1 {
		t.Fatal("Failed to set jsonObject.b.x = 1")
	}
	if !SetByPath(jsonObject, false, "c", 0) || GetByPath(jsonObject, "c", 0) != false {
		t.Fatal("Failed to set jsonObject.c[0]= false")
	}
	if !SetByPath(jsonObject, "z", "b", "z") || GetByPath(jsonObject, "b", "z") != "z" {
		t.Fatal("Failed to set jsonObject.b.z = 'z'")
	}
	if !SetByPath(jsonObject, 123, "c", 3, "x") || GetByPath(jsonObject, "c", 3, "x") != 123 {
		t.Fatal("Failed to set jsonObject.c[3].x = 123")
	}

	if SetByPath(jsonObject, "z", "b", 1) {
		t.Fatal("Failed to set jsonObject.b[1] = x")
	}
	if SetByPath(jsonObject, 123, "c", 5, "x") {
		t.Fatal("Failed to set jsonObject.c[3].x = 123")
	}
}
