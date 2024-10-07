package main

import (
	"bytes"
	"crypto/hmac"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	errors "github.com/BlockChain-Passion/go-project/pkg/error"
	"github.com/BlockChain-Passion/go-project/pkg/logger"
	"github.com/BlockChain-Passion/go-project/pkg/newrelic"
	"github.com/BlockChain-Passion/go-project/pkg/runtimevars"
	"github.com/lpernett/godotenv"
)

func init() {
	if err := godotenv.Load("./.env"); err != nil {
		log.Println("Env file not found in recommendation ", err)
	}
}

type recommendationOrchService struct {
	// interdependencies
	Logger logger.ILogger
	RV     runtimevars.IRV
	NR     newrelic.INewrelic
	DBMap  map[string]postgres.IDatabase
	Error  errors.IErrors
	Hmac   hmac.IHmac
	// Service Info
	ServiceName     string
	ServiceSummary  string
	ServiceOnline   bool
	ServiceProtocol string
	ServiceVersion  string
	ServiceBaseUrl  string
	ServiceRoutes   map[string]interface{}
	GatewayURL      string
}

func (ds *recommendationOrchService) Setup() {
	// Load Runtime Vars
	ds.RV.Load()
	// Setup Logger w/RV
	err := ds.Logger.Setup(ds.RV)
	if err != nil {
		log.Fatal("error while setting up logger: ", err)
	}
	logg := ds.Logger.GetLogrous()
	// setup HMAC KEYS
	logg.Info("HMAC keys are getting initialized")
	//ds.Hmac = &hmac.HmacKeys{Logger: logg, ServiceErrors: ds.Error}
	//err = ds.HMac.LoadKeys{}
	//if err != nil {
	// log.Fatal("error loading HMAC keys: ",err)
	//}

	// Connect and Setup DBMap w/ RV
}

func (ds *recommendationOrchService) Run() {
	logg := ds.Logger.GetLogrous()

	logg.Info("Loading router ...")
	router := routers(routerProps{
		Logger: ds.Logger,
		RV:     ds.RV,
		NR:     ds.NR,
		DBMap:  ds.DBMap,
		HMac:   ds.Hmac,
		Error:  ds.Error,
	})

	// start service
	logg.Info("starting service....")
	logg.Info("micro-service running on PORT: ", ds.RV.Get("PORT"))
	logg.Debug(http.ListenAndServe(":" + ds.RV.Get("PORT")))
}

func (ds *recommendationOrchService) registerAtGateway() error {
	log.Println("registering at gateway....")

	err := json.Unmarshal([]byte(os.Getenv("SERVICE_ROUTES")), &ds.ServiceRoutes)

	if err != nil {
		return fmt.Errorf("error parsing routes json: %v ", err)
	}

	body, err := json.Marshal(struct {
		ServiceName     string                 `json:"service_name"`
		ServiceSummary  string                 `json:"service_summary"`
		ServiceOnline   bool                   `json:"service_online"`
		ServiceProtocol string                 `json:"service_protocol"`
		ServiceVersion  string                 `json:"service_version"`
		BaseUrl         string                 `json:"base_url"`
		Routes          map[string]interface{} `json:"routes"`
	}{
		ds.ServiceName, ds.ServiceSummary, ds.ServiceOnline, ds.ServiceProtocol, ds.ServiceVersion, ds.ServiceBaseUrl, ds.ServiceRoutes,
	})

	if err != nil {
		return fmt.Errorf("error in building request body: %v ", err)
	}

	req, err := http.NewRequest("POST", ds.GatewayURL, bytes.NewBuffer(body))

	if err != nil {
		return fmt.Errorf("error creating request: %v ", err)
	}

	lastestKey := ds.Hmac.GetLatestKey()
	hmacByte := ds.Hmac.CreateHash(req, lastestKey)
	hmac64 := base64.StdEncoding.EncodeToString(hmacByte)
	req.Header.Add("X-HMAC-HASH", hmac64)

	// Perform request
	c := http.DefaultClient
	res, err := c.Do(req)

	if err != nil {
		return fmt.Errorf("error registering service. Status: %v | err: %v ", res.StatusCode, err)
	}

	log.Println("register complete at: ", ds.GatewayURL)
	return nil
}

func main() {
	log.Println("this is from main in recommendations ")
	service := recommendationOrchService{
		Logger: &logger.DefaultLogger{LogLvlKey: "LOG_LEVEL"},
		RV:     &runtimevars.RV{},
		NR:     &newrelic.DefaultNewRelic{AppNameKey: "NEW_RELIC_APP_NAME", LicenseKey: "NEW_REILC_LICENSE", DisplayNameKey: "NEW_RELIC_DISPLAY_NAME"},
		DBMap:  make(map[string]postgres.IDatabase),
		Error:  &errors.ServiceErrors{},
		//Hmac   hmac.IHmac
		// Service Info
		ServiceName:     os.Getenv("SERVICE_NAME"),
		ServiceSummary:  os.Getenv("SERVICE_SUMMARY"),
		ServiceOnline:   true,
		ServiceProtocol: os.Getenv("SERVICE_PROTOCOL"),
		ServiceVersion:  os.Getenv("SERVICE_VERSION"),
		ServiceBaseUrl:  os.Getenv("SERVICE_BASE_URL"),
		ServiceRoutes:   make(map[string]interface{}),
		GatewayURL:      os.Getenv("GATEWAY_URL"),
	}

	service.Setup()

	var err error
	if service.GatewayURL != "" {
		err = service.registerAtGateway()
	}

	if err != nil {
		if string.Contains(err.Error(), "409") {
			log.Println("no error: serivce already registered ")
		} else {
			log.Fatalf("error registering service: %v", err)
		}
	}

	service.Run()
}
