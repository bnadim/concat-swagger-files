// Unmarshal yaml files with JSON compatible Go structures
package csf

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

func ReadYamlFile(file string) (interface{}, error) {
	content, err := ioutil.ReadFile(file)

	if err != nil {
		return nil, err
	}

	return Unmarshal(content)
}

// Unmarshal a yaml content into a json compatible go structures
/*func Unmarshal(input []byte) ([]byte, error) {
	var yamlData interface{}

	err := yaml.Unmarshal(input, &yamlData)
	if err != nil {
		return nil, err
	}

	data := Yaml2JsonGo(yamlData)

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	return jsonData, nil
}*/

// Unmarshal a yaml content into a json compatible go structures
func Unmarshal(input []byte) (interface{}, error) {
	var yamlData interface{}

	err := yaml.Unmarshal(input, &yamlData)
	if err != nil {
		return nil, err
	}

	result := yaml2GoJson(yamlData)

	return result, nil
}

// Convert a yaml value to a json compatible go value
func yaml2GoJson(yamlData interface{}) interface{} {
	switch yamlData.(type) {
	case map[interface{}]interface{}:
		data := yamlData.(map[interface{}]interface{})
		return yamlMap2GoJson(data)

	case []interface{}:
		data := yamlData.([]interface{})
		return yamlArray2GoJson(data)

	default:
		return yamlData
	}
	return yamlData
}

// Convert a yaml array to a json compatible go array (sub maps should have string keys)
func yamlArray2GoJson(yamlData []interface{}) []interface{} {
	result := make([]interface{}, len(yamlData))
	for i, v := range yamlData {
		result[i] = yaml2GoJson(v)
	}

	return result
}

// Convert a yaml map to a json compatible go map (key must be a string)
func yamlMap2GoJson(yamlData map[interface{}]interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	for k, v := range yamlData {
		key := keyToString(k)
		result[key] = yaml2GoJson(v)
	}

	return result
}

// Convert a map key into a string
func keyToString(key interface{}) string {
	return fmt.Sprintf("%v", key)
}