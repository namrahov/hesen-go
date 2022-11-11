package util

import (
	"errors"
	"fmt"
	"github.com/namrahov/hesen-go/model"
	log "github.com/sirupsen/logrus"
)

type IValidationUtil interface {
	ValidationApplicationStatus(applicationStatus model.Status, requestStatus model.Status) error
}

type ValidationUtil struct {
}

func (v *ValidationUtil) ValidationApplicationStatus(applicationStatus model.Status, requestStatus model.Status) error {
	if !canBeChangeTo(applicationStatus, requestStatus) {
		log.Error(fmt.Sprintf("ActionLog.ValidationApplicationStatus.error: %s -> %s is not possible", applicationStatus, requestStatus))
		return errors.New(fmt.Sprintf("Invalid stattus from %s to %s", applicationStatus, requestStatus))
	}
	return nil
}

func canBeChangeTo(applicationStatus model.Status, requestStatus model.Status) bool {
	permissions := statusChangePermissions(applicationStatus)
	for _, permission := range permissions {
		if permission == requestStatus {
			return true
		}
	}
	return false
}

func statusChangePermissions(applicationStatus model.Status) []model.Status {
	var permissions []model.Status
	switch applicationStatus {
	case model.Received:
		permissions = append(permissions, model.Inprogress, model.Hold)
		break
	case model.Inprogress:
		permissions = append(permissions, model.Sent, model.Hold)
		break
	case model.Sent:
		break
	case model.Hold:
		permissions = append(permissions, model.Sent, model.Inprogress)
		break
	}
	return permissions
}
