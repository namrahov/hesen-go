package repo

import (
	"github.com/namrahov/hesen-go/model"
	log "github.com/sirupsen/logrus"
)

type IApplicationRepo interface {
	GetApplicationById(id int64) (*model.Application, error)
}

type ApplicationRepo struct {
}

func (r ApplicationRepo) GetApplicationById(id int64) (*model.Application, error) {
	var application model.Application
	err := Db.Model(&application).
		Column("application.*").
		Where("id = ?", id).
		Select()
	if err != nil {
		log.Fatal(err)
	}

	return &application, nil
}
