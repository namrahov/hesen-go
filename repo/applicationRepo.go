package repo

import (
	"fmt"
	"github.com/namrahov/hesen-go/model"
	log "github.com/sirupsen/logrus"
	"sync"
	"time"
)

type IApplicationRepo interface {
	GetApplicationById(id int64) (*model.Application, error)
	UpdateApplication(application *model.Application) (*model.Application, error)
	SaveApplication(application *model.Application) (*model.Application, error)
	GetDistinctCourtName(wg *sync.WaitGroup) (*[]model.Application, error)
	GetDistinctJudgeName(wg *sync.WaitGroup) (*[]model.Application, error)
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

func (r ApplicationRepo) GetDistinctCourtName(wg *sync.WaitGroup) (*[]model.Application, error) {
	var applications []model.Application
	err := Db.Model(&applications).
		ColumnExpr("DISTINCT court_name").
		Column("application.court_name").
		Select()
	if err != nil {
		log.Fatal(err)
	}
	time.Sleep(3 * time.Second)
	wg.Done()
	return &applications, nil
}

func (r ApplicationRepo) GetDistinctJudgeName(wg *sync.WaitGroup) (*[]model.Application, error) {
	var applications []model.Application
	err := Db.Model(&applications).
		ColumnExpr("DISTINCT judge_name").
		Column("application.judge_name").
		Select()
	if err != nil {
		log.Fatal(err)
	}
	time.Sleep(3 * time.Second)
	wg.Done()
	return &applications, nil
}
