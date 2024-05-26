package entities

import (
	"reflect"
	"testing"
)

func TestReqURL(t *testing.T) {
	expectedFields := map[string]string{
		"ID":          "correlation_id",
		"OriginalURL": "original_url,omitempty",
		"ShortURL":    "short_url",
	}

	validateFields(t, ReqURL{}, expectedFields)
}

func TestReqListAll(t *testing.T) {
	expectedFields := map[string]string{
		"ShortURL":    "short_url",
		"OriginalURL": "original_url",
	}

	validateFields(t, ReqListAll{}, expectedFields)
}

func TestStorage(t *testing.T) {
	expectedTags := map[string]string{
		"UUID":        "id",
		"ShortURL":    "shorturl",
		"OriginalURL": "originalurl",
		"UserID":      "userid",
		"DeletedFlag": "deleted",
	}

	s := Storage{}

	structType := reflect.TypeOf(s)

	for fieldName, expectedTag := range expectedTags {
		field, ok := structType.FieldByName(fieldName)
		if !ok {
			t.Errorf("Field %s not found in structure Storage", fieldName)
			continue
		}

		tag := field.Tag.Get("db")
		if tag != expectedTag {
			t.Errorf("Field %s in structure Storage has incorrect db tag: got %s, want %s", fieldName, tag, expectedTag)
		}
	}
}

// validateFields helper that checks tags
func validateFields(t *testing.T, s interface{}, expectedFields map[string]string) {
	t.Helper()

	structType := reflect.TypeOf(s)

	for fieldName, expectedTag := range expectedFields {
		field, ok := structType.FieldByName(fieldName)
		if !ok {
			t.Errorf("Field %s not found in structure %s", fieldName, structType.Name())
			continue
		}

		tag := field.Tag.Get("json")
		if tag != expectedTag {
			t.Errorf("Field %s in structure %s has incorrect tag: got %s, want %s", fieldName, structType.Name(), tag, expectedTag)
		}
	}
}
