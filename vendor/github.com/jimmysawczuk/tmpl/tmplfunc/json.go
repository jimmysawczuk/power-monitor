package tmplfunc

import (
	"encoding/json"
	"io/ioutil"
)

func GetJSON(path string) (interface{}, error) {
	by, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var target interface{}
	err = json.Unmarshal(by, &target)
	if err != nil {
		return nil, err
	}

	return target, nil
}

func JSONify(v interface{}) (string, error) {
	by, err := json.Marshal(v)
	if err != nil {
		return "", err
	}

	return string(by), nil
}
