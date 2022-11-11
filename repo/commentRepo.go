package repo

import "github.com/namrahov/hesen-go/model"

type ICommentRepo interface {
	SaveComment(comment *model.Comment) error
}

type CommentRepo struct {
}

func (r CommentRepo) SaveComment(comment *model.Comment) error {
	tx, err := Db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Model(comment).Insert()
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}
