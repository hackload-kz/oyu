package converter

import (
	"encoding/json"
)

type Constraints interface {
	[]any |
		[]int |
		[]int32 |
		map[string]any |
		[]map[string]any
	//*router.GetPdFForSignEgovRequest
}

func Convert[c Constraints](object any) (convertedObject c, err error) {
	bytes, err := json.Marshal(object)
	if err != nil {
		return convertedObject, err
	}

	err = json.Unmarshal(bytes, &convertedObject)
	if err != nil {
		return convertedObject, err
	}

	return convertedObject, nil
}
