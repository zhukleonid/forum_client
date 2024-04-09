package convertor

import (
	"encoding/json"
	"lzhuk/clients/model"
	"net/http"
)

func ConvertGetPosts(resp *http.Response) (*model.GetPost, error) {
	getPosts := &model.GetPost{}
	err := json.NewDecoder(resp.Body).Decode(getPosts)
	if err != nil {
		return nil, err
	}

	return getPosts, nil
}
