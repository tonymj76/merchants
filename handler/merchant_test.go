package handler

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/Investliftng/ocm-api/merchant/datastore"
	pbM "github.com/Investliftng/ocm-api/merchant/proto/merchant"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/sirupsen/logrus"
)

var (
	link string
	log  *logrus.Logger
)

func TestMain(m *testing.M) {
	code := 0
	defer func() {
		os.Exit(code)
	}()

	log = logrus.New()
	log.Formatter = &logrus.TextFormatter{
		TimestampFormat: time.RFC3339,
		FullTimestamp:   true,
		ForceColors:     true,
	}

	pool, err := dockertest.NewPool("")
	if err != nil {
		log.WithError(err).Fatal("Could not connect to docker")
	}

	src := map[string]string{
		"user":     "postgres",
		"password": "password",
		"db":       "merchantTest",
	}

	runOpts := dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "latest",
		Env: []string{
			"POSTGRES_USER=" + src["user"],
			"POSTGRES_PASSWORD=" + src["password"],
			"POSTGRES_DB=" + src["db"],
		},
	}
	resource, err := pool.RunWithOptions(&runOpts)
	if err != nil {
		log.WithError(err).Fatal("could not start postgres container")
	}
	defer func() {
		err = pool.Purge(resource)
		if err != nil {
			log.WithError(err).Error("Could not purge resource")
		}
	}()
	logWaiter, err := pool.Client.AttachToContainerNonBlocking(docker.AttachToContainerOptions{
		Container:    resource.Container.ID,
		OutputStream: log.Writer(),
		ErrorStream:  log.Writer(),
		Stderr:       true,
		Stdout:       true,
		Stream:       true,
	})
	if err != nil {
		log.WithError(err).Fatal("could not connect to postgres container log output")
	}
	defer func() {
		err = logWaiter.Close()
		if err != nil {
			log.WithError(err).Error("Could not wait for container log to close")
		}
	}()

	pool.MaxWait = 10 * time.Second
	link = fmt.Sprintf("postgres://%s:%s@localhost:%s/%s?sslmode=disable", src["user"], src["password"], resource.GetPort("5432/tcp"), src["db"])
	if err = pool.Retry(func() error {
		db, err := sql.Open("postgres", link)
		if err != nil {
			return err
		}
		return db.Ping()
	}); err != nil {
		log.WithError(err).Fatal("Could not connect to postgres server")
	}
	code = m.Run()
}

func TestGetListAddDelectMerchant(t *testing.T) {
	// srv := micro.NewService(
	// 	micro.Name("service.merchant"),
	// )
	// srv.Init()
	connect, err := datastore.NewConnection(log, link)
	if err != nil {
		t.Fatalf("Failed to create a new connection: %s\n", err)
	}
	defer func() {
		err = connect.Close()
		if err != nil {
			t.Errorf("Failed to close directory: %s", err)
		}
	}()

	// pbM.RegisterMerchantServiceHandler(srv.Server(), &Service{connect})
	// client := pbM.NewMerchantService("service.merchant", srv.Client())
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// if err := srv.Run(); err != nil {
	// 	t.Fatalf("Failed to Run micro connection: %s\n", err)
	// }
	t.Run("creating a merchant", func(t *testing.T) {
		mr := new(pbM.MerchantRequest)
		role := pbM.RoleType_SALE_PERSON
		phone := &pbM.PhoneNumber{
			Number: "09099999",
			Type:   pbM.PhoneType_MOBILE,
		}
		mr.BusinessName = "lifthub"
		mr.Email = "lift@g.com"
		mr.NumberOfOutlet = 1
		mr.NumberOfProduct = 33
		mr.UserId = 34
		mr.Role = role
		mr.Phone = phone
		// m1, err := client.CreateMerchant(ctx, &pbM.CreateRequest{
		// 	Password: "new password2",
		// 	Merchant: mr,
		// })
		// if err != nil {
		// 	t.Fatalf("failed to add a merchant: %s", err)
		// }
		// if m1.GetMerchant().Role != role {
		// 	t.Errorf("Got role %q, want role %q", m1.Merchant.GetRole(), role)
		// }

		err := connect.Create(ctx, &pbM.CreateRequest{
			Password: "new password2",
			Merchant: mr,
		})

		if err != nil {
			t.Fatalf("failed to add a merchant: %s", err)
		}
	})
}
