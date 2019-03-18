package provider

import (
	"io/ioutil"
	"net/http"
)

type WebRequestSingleValue struct {
	Url   string
	Value string
}

func (this *WebRequestSingleValue) GetValue() (*string, error) {
	if len(this.Value) == 0 {
		req, err := http.NewRequest("GET", this.Url, nil)
		if err != nil {
			return nil, err
		}

		client := &http.Client{}
		response, err := client.Do(req)

		defer response.Body.Close()
		if err != nil {
			return nil, err
		}

		contents, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return nil, err
		}

		this.Value = string(contents)
	}

	return &this.Value, nil
}
