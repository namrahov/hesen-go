package repo

import (
	"fmt"
	"github.com/namrahov/hesen-go/model"
	log "github.com/sirupsen/logrus"
)

type IApplicationRepo interface {
	GetApplicationById(id int64) (*model.Application, error)
	UpdateApplication(application *model.Application) (*model.Application, error)
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

func (r ApplicationRepo) UpdateApplication(application *model.Application) (*model.Application, error) {
	_, err := Db.Model(application).
		OnConflict("(id) DO UPDATE").
		Insert()
	if err != nil {
		log.Fatal(err)
	}
	return application, nil
}

func (r ApplicationRepo) SaveApplication(application *model.Application) (*model.Application, error) {
	tx, err := Db.Begin()
	if err != nil {
		return nil, err
	}
	_, err = tx.Model(application).Insert()
	if err != nil {
		return nil, err
	}
	for idx := range application.Comments {
		application.Comments[idx].ApplicationId = application.Id
	}
	_, err = tx.Model(&application.Comments).Insert()
	if err != nil {
		return nil, fmt.Errorf("%w, %v", err, tx.Rollback())
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	return application, nil
}
