package dao

import (
	"github.com/itering/subscan-plugin/storage"
	"github.com/itering/subscan/plugins/imtoken/model"
)

func NewTransfer(db storage.DB, transfer *model.EventTransfer) error {
	err := db.Create(transfer)
	if err != nil {
		return err
	}

	return nil
}
