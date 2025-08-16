package rest_controller

import (
	"encoding/json"
	"net/http"
)

type multipartData interface {
	//entity.FileSign
}

func multipartParseByData[D multipartData](r *http.Request) (
	data D, err error,
) {
	err = r.ParseMultipartForm(32 << 20)
	if err != nil {
		return data, err
	}

	formData := r.MultipartForm.Value["data"]

	err = json.Unmarshal([]byte(formData[0]), &data)
	if err != nil {
		return data, err
	}

	return data, nil
}
