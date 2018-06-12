// Even if the name is jsonref, this module works on any go structure json compatible and fetches
// referenced modules
package csf

import (
	"reflect"
	"fmt"
	"net/url"
	"strings"
	"path"
)

type jsonRefResolver struct {
	parentsRefs map[string]bool // A map of imported references in current cycle
	jsonData interface{} // full json data. Used to fetch internal pointers
	pwd string // initial folder to resolve json refs
}

func ResolveJsonRefs(jsonData interface{}, pwd string) (interface{}, error) {
	resolver := jsonRefResolver{
		parentsRefs: make(map[string]bool),
		jsonData: jsonData,
		pwd: pwd,
	}

	return resolver.resolve(jsonData)
}

func (resolver *jsonRefResolver) resolve(jsonData interface{}) (interface{}, error) {
	switch jsonData.(type) {
	case map[string]interface{}:
		typedJsonData := jsonData.(map[string]interface{})
		return resolver.resolveMap(typedJsonData)

	case []interface{}:
		typedJsonData := jsonData.([]interface{})
		return resolver.resolveArray(typedJsonData)
	}

	return jsonData, nil
}

func (resolver *jsonRefResolver) resolveArray(jsonArray []interface{}) ([]interface{}, error) {
	result := make([]interface{}, len(jsonArray))

	for i, v := range jsonArray {
		value, err := resolver.resolve(v)
		if err != nil {
			return nil, err
		}

		result[i] = value
	}

	return result, nil
}

func (resolver *jsonRefResolver) resolveMap(jsonMap map[string]interface{}) (interface{}, error) {
	if _, isRef := jsonMap["$ref"]; isRef {
		if reflect.TypeOf(jsonMap["$ref"]).Kind() != reflect.String {
			err := fmt.Errorf("error: malformated json reference, value should be a uri but got %#v", jsonMap["$ref"])
			return nil, err
		}

		return resolver.resolveRef(jsonMap["$ref"].(string))
	}

	result := make(map[string]interface{})
	for k, v := range jsonMap {
		jsonData, err := resolver.resolve(v)
		if err != nil {
			return nil, err
		}

		result[k] = jsonData
	}

	return result, nil
}

func (resolver *jsonRefResolver) resolveRef(jsonRef string) (interface{}, error) {
	// Check cyclic dependencies
	if isParentRef, _ := resolver.parentsRefs[jsonRef]; isParentRef {
		return nil, fmt.Errorf("error: cyclic reference for json ref %s", jsonRef)
	}
	resolver.parentsRefs[jsonRef] = true

	// Parse jsonref uri components
	parsedJsonRef, err := url.Parse(jsonRef)
	if err != nil { return nil, err }

	// Check that it has at least a file path and a jsonPointer
	if len(parsedJsonRef.Path) + len(parsedJsonRef.Fragment) <= 0 {
		err := fmt.Errorf("error: a json ref should have at least a path or a key, but none found in %s", jsonRef)
		return nil, err
	}

	// Now resolve file path if need and then json pointer
	jsonData := resolver.jsonData
	fileJsonData := resolver.jsonData
	pwd := resolver.pwd

	if len(parsedJsonRef.Path) > 0 {
		filepath := path.Join(resolver.pwd, parsedJsonRef.Path)
		jsonData, err = ReadYamlFile(filepath)
		pwd = path.Dir(filepath)
		if err != nil { return nil, err }
		fileJsonData = jsonData
	}

	if len(parsedJsonRef.Fragment) > 0 {
		jsonData, err = valueForJsonPointer(jsonData, parsedJsonRef.Fragment)
		if err != nil { return nil, err }
	}

	// Now resolve sub references
	subResolver := jsonRefResolver{
		parentsRefs: resolver.parentsRefs,
		jsonData: fileJsonData,
		pwd: pwd,
	}

	result, err := subResolver.resolve(jsonData)
	if err != nil { return nil, err }

	// Don't forget to say we are not a parent anymore
	resolver.parentsRefs[jsonRef] = false

	return result, nil
}

func valueForJsonPointer(jsonData interface{}, jsonPointer string) (interface{}, error) {
	currentData := jsonData

	ptrComponents := strings.Split(jsonPointer, "/")

	for _, ptrComponent := range ptrComponents[1:] {
		switch currentData.(type) {
		case map[string]interface{}:
			typedCurrentData := currentData.(map[string]interface{})
			if _, exists := typedCurrentData[ptrComponent]; !exists {
				return nil, fmt.Errorf("error: ref error, can not find component %s of json pointer %s in data %v", ptrComponent, jsonPointer, currentData)
			}
			currentData = typedCurrentData[ptrComponent]

		default:
			return nil, fmt.Errorf("error: ref error, can not find component %s of json pointer %s in data %v", ptrComponent, jsonPointer, currentData)
		}
	}

	return currentData, nil
}


