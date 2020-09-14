package repository

import (
	"encoding/json"
	"io/ioutil"

	"github.com/pkg/errors"

	"github.com/shipa988/fanouter/internal/domain/entity"
)

const (
	ErrLoadJson = "can't load json config"
)

var _ entity.FanParamRepo = (*JSONRepo)(nil)

type JSONRepo struct {
	jsonPath string
}

func NewJSONRepo(path string) *JSONRepo {
	return &JSONRepo{jsonPath: path}
}

func (r *JSONRepo) Load() (*entity.FanParam, error) {
	dat, err := ioutil.ReadFile(r.jsonPath)
	if err != nil {
		return nil, errors.Wrapf(err, ErrLoadJson)
	}
	var urls entity.FanParam
	err = json.Unmarshal(dat, &urls)
	if err != nil {
		return nil, errors.Wrapf(err, ErrLoadJson)
	}
	return &urls, nil
}
