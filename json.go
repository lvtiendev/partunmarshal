package partunmarhal

import (
	"encoding/json"
	"reflect"
	"strings"
)

// JSON partially unmarshal payload into obj
func JSON(obj interface{}, payload []byte) error {
	v := reflect.ValueOf(obj)
	if v.Kind() != reflect.Ptr {
		return ErrPointerExpected
	}

	elem := v.Elem()
	typ := elem.Type()

	// Unmarshal to a temporary map just to check if a key is present in payload
	unmarshalled := map[string]interface{}{}
	if err := json.Unmarshal(payload, &unmarshalled); err != nil {
		return err
	}

	// Umarshall the payload into newObj of same type as obj
	// so that any custom field's type unmarshalling works as normal.
	newObj := reflect.New(typ).Interface()
	if err := json.Unmarshal(payload, newObj); err != nil {
		return err
	}
	newElem := reflect.ValueOf(newObj).Elem()

	// Hold the to-be-update field values
	toBeUpdate := make(map[int]reflect.Value, typ.NumField())

	// In the first loop, check for error and set toBeUpdate
	for i := 0; i < typ.NumField(); i++ {
		typField := typ.Field(i)
		elemField := elem.Field(i)
		newField := newElem.Field(i)

		if typField.Tag.Get(tagU) != updatable {
			continue
		}

		if !elemField.CanSet() {
			return ErrFieldCannotBeSet
		}

		jsonTagList := strings.Split(typField.Tag.Get(tagJSON), ",")
		jsonKey := jsonTagList[0]

		if jsonKey == "" {
			return ErrNoTagJSON
		}

		// Ignore if the payload does not contain the key
		if _, ok := unmarshalled[jsonKey]; !ok {
			continue
		}

		// Here we don't set the newField to elemField immidately
		// because any error in the loop can cause obj to be inconsistent
		toBeUpdate[i] = newField
	}

	// In the 2nd loop, set the field values into obj
	for i, value := range toBeUpdate {
		elem.Field(i).Set(value)
	}

	return nil
}
