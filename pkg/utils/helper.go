package utils

import "encoding/json"

func MarshalUnmarshal(reqStructure, resStructure interface{}) error {
	byte, err := json.Marshal(reqStructure)
	if err != nil {
		return err
	}
	err = json.Unmarshal(byte, resStructure)
	return err
}
