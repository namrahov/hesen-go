package repo

import (
	"github.com/namrahov/hesen-go/model"
)

type IUserRepo interface {
	GetUserById(id string) (*model.User, error)
	GetSessionBySessionId(sessionId string) (*[]model.Session, error)
	SaveUser(u *model.User) (*model.User, error)
	SaveSession(u *model.Session) error
}

type UserRepo struct {
}

func (r UserRepo) GetUserById(id string) (*model.User, error) {
	var user model.User
	err := Db.Model(&user).
		Where("id = ?", id).
		Select()
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r UserRepo) GetSessionBySessionId(sessionId string) (*[]model.Session, error) {
	var session []model.Session
	err := Db.Model(&session).
		Where("session_id = ?", sessionId).
		Select()

	if err != nil {
		return nil, err
	}

	return &session, nil
}

func (r UserRepo) SaveUser(u *model.User) (*model.User, error) {
	tx, err := Db.Begin()
	if err != nil {
		return nil, err
	}
	_, err = tx.Model(u).Insert()
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (r UserRepo) SaveSession(s *model.Session) error {
	tx, err := Db.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Model(s).Insert()
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}
