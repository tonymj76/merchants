package datastore

import (
	"database/sql/driver"
	"fmt"

	pbM "github.com/Investliftng/ocm-api/merchant/proto/merchant"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/duration"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/jackc/pgx/pgtype"
)

type roleTypeWrapper pbM.RoleType

// Value implements database/sql/driver.Valuer for merchant.RoleType
func (rw roleTypeWrapper) Value() (driver.Value, error) {
	switch pbM.RoleType(rw) {
	case pbM.RoleType_SUPER_ADMIN:
		return "super_admin", nil
	case pbM.RoleType_ADMIN:
		return "admin", nil
	case pbM.RoleType_SUB_ADMIN:
		return "sub_admin", nil
	case pbM.RoleType_SALE_PERSON:
		return "sale_person", nil
	default:
		return nil, fmt.Errorf("invalid Role: %q", rw)
	}
}

// Scan implements database/sql/driver.Scanner for pbpbMs.Role
func (rw *roleTypeWrapper) Scan(in interface{}) error {
	switch in.(string) {
	case "super_admin":
		*rw = roleTypeWrapper(pbM.RoleType_SUPER_ADMIN)
		return nil
	case "admin":
		*rw = roleTypeWrapper(pbM.RoleType_ADMIN)
		return nil
	case "sub_admin":
		*rw = roleTypeWrapper(pbM.RoleType_SUB_ADMIN)
		return nil
	case "sale_person":
		*rw = roleTypeWrapper(pbM.RoleType_SALE_PERSON)
		return nil
	default:
		return fmt.Errorf("invalid Role: %q", in.(string))
	}
}

type phoneTypeWrapper pbM.PhoneType

// Value implements database/sql/driver.Valuer for merchant.RoleType
func (pw phoneTypeWrapper) Value() (driver.Value, error) {
	switch pbM.PhoneType(pw) {
	case pbM.PhoneType_HOME:
		return "home", nil
	case pbM.PhoneType_MOBILE:
		return "mobile", nil
	case pbM.PhoneType_WORK:
		return "work", nil
	default:
		return nil, fmt.Errorf("invalid Role: %q", pw)
	}
}

// Scan implements database/sql/driver.Scanner for pbpbMs.Role
func (pw *phoneTypeWrapper) Scan(in interface{}) error {
	switch in.(string) {
	case "home":
		*pw = phoneTypeWrapper(pbM.PhoneType_HOME)
		return nil
	case "mobile":
		*pw = phoneTypeWrapper(pbM.PhoneType_MOBILE)
		return nil
	case "work":
		*pw = phoneTypeWrapper(pbM.PhoneType_WORK)
		return nil
	default:
		return fmt.Errorf("invalid Role: %q", in.(string))
	}
}

type timeWrapper timestamp.Timestamp

// Value implements database/sql/driver.Valuer for timestamp.Timestamp
func (tw *timeWrapper) Value() (driver.Value, error) {
	return ptypes.Timestamp((*timestamp.Timestamp)(tw))
}

// Scan implements database/sql/driver.Scanner for timestamp.Timestamp
func (tw *timeWrapper) Scan(in interface{}) error {
	var t pgtype.Timestamptz
	err := t.Scan(in)
	if err != nil {
		return err
	}
	tp, err := ptypes.TimestampProto(t.Time)
	if err != nil {
		return err
	}
	*tw = (timeWrapper)(*tp)
	return nil
}

type durationWrapper duration.Duration

// Value implements database/sql/driver.Valuer for duration.Duration
func (dw *durationWrapper) Value() (driver.Value, error) {
	d, err := ptypes.Duration((*duration.Duration)(dw))
	if err != nil {
		return nil, err
	}

	i := pgtype.Interval{
		Microseconds: int64(d) / 1000,
		Status:       pgtype.Present,
	}

	return i.Value()
}
