package datastore

import (
	"database/sql"
	"errors"

	"github.com/Investliftng/ocm-api/merchant/datastore/migrations"
	pbM "github.com/Investliftng/ocm-api/merchant/proto/merchant"
	"github.com/Masterminds/squirrel"
	migrate "github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	bindata "github.com/golang-migrate/migrate/v4/source/go_bindata"
	"github.com/golang/protobuf/ptypes/timestamp"
	pbError "github.com/micro/go-micro/v2/errors"
	"github.com/sirupsen/logrus"
)

// version defines the current migration version. This ensures the app
// is always compatible with the version of the database.
const version = 4

func validateSchema(db *sql.DB) error {
	sourceInstance, err := bindata.WithInstance(bindata.Resource(migrations.AssetNames(), migrations.Asset))
	if err != nil {
		return err
	}

	targetInstance, err := postgres.WithInstance(db, new(postgres.Config))
	if err != nil {
		return err
	}
	m, err := migrate.NewWithInstance("go-bindata", sourceInstance, "postgres", targetInstance)
	if err != nil {
		return err
	}
	err = m.Migrate(version) // current version
	if err != nil && err != migrate.ErrNoChange {
		return err
	}
	return sourceInstance.Close()
}

func scanMerchant(row squirrel.RowScanner) (*pbM.Merchant, error) {
	m := pbM.Merchant{}
	phone := pbM.PhoneNumber{}
	m.CreatedAt = new(timestamp.Timestamp)
	m.UpdatedAt = new(timestamp.Timestamp)
	m.LastLogin = new(timestamp.Timestamp)
	err := row.Scan(
		&m.Id,
		&m.NumberOfProduct,
		(*roleTypeWrapper)(&m.Role),
		&m.Email,
		&phone.Number,
		(*phoneTypeWrapper)(&phone.Type),
		&m.UserId,
		&m.NumberOfOutlet,
		&m.BusinessName,
		&m.IsSuspended,
		&m.IsEmailVerified,
		(*timeWrapper)(m.LastLogin),
		(*timeWrapper)(m.CreatedAt),
		(*timeWrapper)(m.UpdatedAt),
	)
	m.Phone = &phone
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logrus.Error("no merchant found", err.Error())
			return nil, pbError.BadRequest("404", "no merchant found", err.Error())
		}
		return nil, err
	}
	logrus.WithFields(logrus.Fields{
		"ID":      m.Id,
		"role":    m.GetRole(),
		"created": m.CreatedAt,
		"phone":   &phone,
	}).Info("scaning into response")
	return &m, nil
}

func scanOutlet(row squirrel.RowScanner) (*pbM.Outlet, error) {
	m := pbM.Outlet{}
	phone := pbM.PhoneNumber{}
	m.CreatedAt = new(timestamp.Timestamp)
	m.UpdatedAt = new(timestamp.Timestamp)
	err := row.Scan(
		&m.Id,
		&m.MerchantId,
		&phone.Number,
		(*phoneTypeWrapper)(&phone.Type),
		&m.Latitude,
		&m.Longitude,
		&m.Position,
		&m.CityId,
		&m.CountryId,
		&m.Address,
		&m.Available,
		(*timeWrapper)(m.CreatedAt),
		(*timeWrapper)(m.UpdatedAt),
	)
	m.Phone = &phone
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logrus.Error("No Outlet found", err.Error())
			return nil, pbError.BadRequest("404", "No Outlet found", err.Error())
		}
		return nil, err
	}
	logrus.WithFields(logrus.Fields{
		"ID":      m.Id,
		"created": m.CreatedAt,
		"phone":   &phone,
	}).Info("scaning into response")
	return &m, nil
}

func scanTerminal(row squirrel.RowScanner) (*pbM.Terminal, error) {
	m := pbM.Terminal{}
	m.CreatedAt = new(timestamp.Timestamp)
	m.UpdatedAt = new(timestamp.Timestamp)
	err := row.Scan(
		&m.Id,
		&m.MerchantId,
		&m.UserId,
		&m.OutletId,
		&m.Name,
		(*timeWrapper)(m.CreatedAt),
		(*timeWrapper)(m.UpdatedAt),
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logrus.Error("No Terminal found", err.Error())
			return nil, pbError.BadRequest("404", "No Terminal found", err.Error())
		}
		return nil, err
	}
	logrus.WithFields(logrus.Fields{
		"ID":      m.Id,
		"created": m.CreatedAt,
	}).Info("scaning into response")
	return &m, nil
}

func scanMerchantID(row squirrel.RowScanner) (*pbM.Merchant, error) {
	m := pbM.Merchant{}
	err := row.Scan(
		&m.Id,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logrus.Error("no merchant found with that id", err.Error())
			return nil, pbError.BadRequest("404", "no merchant found with that id", err.Error())
		}
		return nil, err
	}
	logrus.WithFields(logrus.Fields{
		"ID": m.Id,
	}).Info("scaning into response")
	return &m, nil
}

func scanOutletMerchantID(row squirrel.RowScanner) (uint64, error) {
	var merchantID uint64

	err := row.Scan(&merchantID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logrus.Error("no merchant found with that id", err.Error())
			return 0, pbError.BadRequest("404", "no merchant id found", err.Error())
		}
		return 0, err
	}
	return merchantID, nil
}

// will be use in merchant login
// func scanMerchantEmailAndPassword(row squirrel.RowScanner) (*pbM.Merchant, error) {
// 	m := pbM.MerchantEmailAndPassword{}
// 	err := row.Scan(
// 		&m.Email,
// 		&
// 	)

// 	if err != nil {
// 		if errors.Is(err, sql.ErrNoRows) {
// 			logrus.Error("no merchant found", err.Error())
// 			return nil, pbError.BadRequest("404", "no merchant found", err.Error())
// 		}
// 		return nil, err
// 	}
// 	logrus.WithFields(logrus.Fields{
// 		"ID":      m.Id,
// 		"role":    m.GetRole(),
// 		"created": m.CreatedAt,
// 	}).Info("scaning into response")
// 	return &m, nil
// }
