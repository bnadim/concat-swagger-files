package csf

import (
	"io/ioutil"

	"bytes"
	"encoding/json"
	"path"
)

func readFile(filePath string) ([]byte, error) {
	content, err := ioutil.ReadFile(filePath)

	if err != nil {
		return nil, err
	}

	return content, nil
}

func writeJsonFile(filePath string, jsonData interface{}) error {
	output, err := json.Marshal(jsonData)
	if err != nil { return err }

	var prettyOutput bytes.Buffer
	json.Indent(&prettyOutput, output, "", "    ")
	return ioutil.WriteFile(filePath, prettyOutput.Bytes(), 0644)
}



func Convert(inputPath string, outputPath string) error {
	jsonData, err := ReadYamlFile(inputPath)
	if err != nil {
		return err
	}

	resolvedData, err := ResolveJsonRefs(jsonData, path.Dir(inputPath))
	//fmt.Printf("Error is %v", err)
	if err != nil {
		return err
	}

	return writeJsonFile(outputPath,resolvedData)
	//fmt.Printf("Result is %#v", result)

	return nil
}