package testhelpers

import "io/ioutil"

// GetMsgFromFile gets a
func GetMsgFromFile(filepath string) ([]byte, error) {
	jsonFile, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	return jsonFile, nil
}
