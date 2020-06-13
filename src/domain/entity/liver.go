package entity

import "github.com/shuymn/nijisanji-db-collector/src/domain"

type Liver struct {
	ID        string       `validate:"required"`
	YouTube   []*YouTube   `validate:"omitempty"`
	Twitter   []*Twitter   `validate:"omitempty"`
	Bilibili  []*Bilibili  `validate:"omitempty"`
	Twitch    []*Twitch    `validate:"omitempty"`
	Facebook  []*Facebook  `validate:"omitempty"`
	Instagram []*Instagram `validate:"omitempty"`
}

func (l *Liver) Validate() error {
	return domain.Validator.Struct(l)
}

type YouTube struct {
	ID string `validate:"omitempty"`
}

type Twitter struct {
	ID string `validate:"omitempty"`
}

type Bilibili struct {
	ID string `validate:"omitempty"`
}

type Twitch struct {
	ID string `validate:"omitempty"`
}

type Facebook struct {
	ID  string `validate:"omitempty"`
	URL string `validate:"omitempty"`
}

type Instagram struct {
	ID  string `validate:"omitempty"`
	URL string `validate:"omitempty"`
}
