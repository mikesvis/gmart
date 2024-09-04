package service_test

import (
	"bytes"
	"context"
	"log"
	"net/http"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/mikesvis/gmart/internal/config"
	accrualExchange "github.com/mikesvis/gmart/internal/exchange/accrual"
	"github.com/mikesvis/gmart/internal/logger"
	"github.com/mikesvis/gmart/internal/postgres"
	"github.com/mikesvis/gmart/internal/queue"
	"github.com/mikesvis/gmart/internal/randomstring"
	accrualService "github.com/mikesvis/gmart/internal/service/accrual"
	orderService "github.com/mikesvis/gmart/internal/service/order"
	userService "github.com/mikesvis/gmart/internal/service/user"
	accrualStorage "github.com/mikesvis/gmart/internal/storage/accrual"
	orderStorage "github.com/mikesvis/gmart/internal/storage/order"
	userStorage "github.com/mikesvis/gmart/internal/storage/user"
	"github.com/mikesvis/gmart/pkg/luhn"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"golang.org/x/exp/rand"
)

type User struct {
	Login    string
	Password string
	ID       uint64
}

var UserService *userService.Service
var OrderService *orderService.Service
var AccrualService *accrualService.Service
var EmptyOrdersExistingUserID uint64
var NotEmptyOrdersExistingUserID uint64
var ExistingOrderNumberWithAccrual uint64
var ExistingOrderNumberWithWithdrawal uint64
var ExistingOrderProcessedAtWithWithdrawal time.Time
var ExistingUser *User

func TestService(t *testing.T) {
	RegisterFailHandler(Fail)

	BeforeSuite(func() {
		os.Setenv("RUN_ADDRESS", "127.0.0.1:8989")
		os.Setenv("DATABASE_URI", "host=0.0.0.0 port=5432 user=postgres password=postgres dbname=gmart sslmode=disable")
		os.Setenv("ACCRUAL_SYSTEM_ADDRESS", "http://127.0.0.1:8081")
		ctx := context.Background()

		config := config.NewConfig()
		db, err := postgres.NewPostgres(config)
		if err != nil {
			panic(err)
		}
		logger := logger.NewLogger()

		rand.Seed(uint64(GinkgoRandomSeed()))

		userStorage := userStorage.NewStorage(db, logger)
		UserService = userService.NewService(userStorage, logger)
		orderStorage := orderStorage.NewStorage(db, logger)

		queue := queue.NewQueue()
		accrualStorage := accrualStorage.NewStorage(db, logger)
		accrualExchange := accrualExchange.NewAccrualExchange(config, logger)

		UserService = userService.NewService(userStorage, logger)
		OrderService = orderService.NewService(orderStorage, logger)
		AccrualService = accrualService.NewService(accrualStorage, accrualExchange, queue, logger)
		go AccrualService.RunQueue(context.TODO())

		/* create random user */
		login := randomstring.RandStringRunes(10)
		password := "password"
		userID, _ := UserService.RegisterUser(ctx, login, password)

		ExistingUser = &User{
			Login:    login,
			Password: password,
			ID:       userID,
		}

		/* create random user with no orders */
		login = randomstring.RandStringRunes(10)
		password = "password"
		userID, _ = UserService.RegisterUser(ctx, login, password)
		EmptyOrdersExistingUserID = userID

		/* create random user with orders */
		login = randomstring.RandStringRunes(10)
		password = "password"
		userID, _ = UserService.RegisterUser(ctx, login, password)
		NotEmptyOrdersExistingUserID = userID

		/* create order with accrual */
		ExistingOrderNumberWithAccrual = GenerateRandomLuhn()
		createAccrualClientAndRegisterOrder(config, ExistingOrderNumberWithAccrual)
		OrderService.CreateOrder(ctx, ExistingOrderNumberWithAccrual, NotEmptyOrdersExistingUserID)
		AccrualService.PushToAccural(ExistingOrderNumberWithAccrual)
		time.Sleep(1 * time.Second)

		/* create order with withdrawal */
		ExistingOrderNumberWithWithdrawal = GenerateRandomLuhn()
		OrderService.CreateOrder(ctx, ExistingOrderNumberWithWithdrawal, NotEmptyOrdersExistingUserID)
		ExistingOrderProcessedAtWithWithdrawal = time.Now()
		AccrualService.WithdrawToOrderID(ctx, ExistingOrderNumberWithWithdrawal, -10000)
		time.Sleep(1 * time.Second)
	})

	RunSpecs(t, "Service Suite")
}

func GenerateRandomLuhn() uint64 {
	var randLuhn uint64
	for {
		randLuhn = uint64(rand.Uint32() + 10000)
		if luhn.IsValid(randLuhn) {
			break
		}
	}

	return randLuhn
}

func createAccrualClientAndRegisterOrder(cfg *config.Config, orderNumber uint64) {
	client := &http.Client{}
	good := randomstring.RandStringRunes(10)

	/* register good in accrual */
	body := []byte(`{"match" : "` + good + `","reward": 10,"reward_type": "%"}`)
	request, _ := http.NewRequest(http.MethodPost, cfg.AccrualSystemAddress+"/api/goods", bytes.NewBuffer(body))
	request.Header.Set("Content-Type", "application/json")
	log.Printf("creating good %s in accrual, body is %s", good, string(body))
	response, err := client.Do(request)
	if err != nil {
		log.Panicf("%v", err)
	}
	defer response.Body.Close()
	log.Printf("%s", response.Status)
	time.Sleep(1 * time.Second)

	/* register order in accrual */
	body = []byte(`{"order": "` + strconv.FormatUint(orderNumber, 10) + `","goods": [{"description": "` + good + `","price": 7000}]}`)
	request, _ = http.NewRequest(http.MethodPost, cfg.AccrualSystemAddress+"/api/orders", bytes.NewBuffer(body))
	request.Header.Set("Content-Type", "application/json")
	log.Printf("registering order %d in accrual, body is %s", orderNumber, string(body))
	response, err = client.Do(request)
	if err != nil {
		log.Panicf("%v", err)
	}
	defer response.Body.Close()
	log.Printf("%s", response.Status)
	time.Sleep(1 * time.Second)
}
