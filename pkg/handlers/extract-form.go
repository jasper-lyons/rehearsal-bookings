package handlers

import (
	"net/http"
	"io/ioutil"
	"encoding/json"
)


func ExtractForm[FormType any](r *http.Request) (*FormType, error) {
		var form *FormType 
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(body, &form)

		return form, nil
}
