package repo

import (
	"github.com/namrahov/hesen-go/model"
	log "github.com/sirupsen/logrus"
)

type IApplicationRepo interface {
	GetApplicationById(id int64) (*model.Application, error)
	SaveApplication(application *model.Application) (*model.Application, error)
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

func (r ApplicationRepo) SaveApplication(application *model.Application) (*model.Application, error) {
	_, err := Db.Model(application).
		OnConflict("(id) DO UPDATE").
		Insert()
	if err != nil {
		log.Fatal(err)
	}
	return application, nil
}
