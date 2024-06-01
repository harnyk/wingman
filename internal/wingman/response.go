package wingman

import (
	"encoding/json"
	"strings"
)

type Response struct {
	Command     string `json:"command"`
	Explanation string `json:"explanation"`
}

func ParseResponseJSON(response string) (Response, error) {

	response = strings.Trim(response, "\n\r ")
	res := Response{}

	if err := json.Unmarshal([]byte(response), &res); err != nil {
		return Response{}, err
	}

	return res, nil
}
