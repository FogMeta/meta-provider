package pocket

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/fatih/color"
	"github.com/filswan/go-swan-lib/client/swan"
	"github.com/filswan/go-swan-lib/client/web"
	"github.com/filswan/go-swan-lib/utils"
	"github.com/gin-gonic/gin"
	cors "github.com/itsjamie/gin-cors"
	"strconv"
	"strings"
	"time"
)

var poktService *PoktService

func ParsePoktCmd(cmd []string) {
	if len(cmd) < 2 {
		printPoktUsage()
		return
	}

	subCmd := cmd[1]
	switch subCmd {
	case "start":
		cmdPoktStart(cmd[1:])
		poktHttpServer()
	case "version":
		cmdPoktVersion()
	case "validator":
		cmdPoktNodeAddr()
	case "balance":
		cmdPoktBalance(cmd[1:])
	case "custodial":
		cmdPoktCustodial(cmd[1:])
	case "noncustodial":
		cmdPoktNonCustodial(cmd[1:])
	case "status":
		cmdPoktStatus()
	default:
		printPoktUsage()
	}

}

func cmdPoktStart(op []string) {
	confPokt := GetConfig()
	SetLogLevel(confPokt.PoktLogLevel)

	poktService = GetMyPoktService()

	poktService.StartPoktContainer(op)
	for {
		ver, err := poktService.dkCli.PoktCtnExecVersion()
		if err == nil {
			GetLog().Info("Pocket Node Version Is:", ver.Version)
			time.Sleep(time.Second * 3)
			break
		}
		GetLog().Info("Wait for Pocket Node Available...")
		time.Sleep(time.Second * 3)
	}

	//
	acc, err := poktService.dkCli.PoktCtnExec([]string{"pocket", "accounts", "list"})
	if err != nil {
		GetLog().Error("Get Pocket Accounts error:", err)
	}
	if !strings.Contains(acc, "(0) ") {
		panic("Get Init Pocket Accounts Error.")
	}
	poktService.PoktAddress = strings.Split(acc, "(0) ")[1][0:40]
	poktService.CurStatus.Address = poktService.PoktAddress
	GetLog().Info("Pocket Accounts is:", poktService.PoktAddress)

	// not ready for heartbeat
	go sendHeartbeat2Swan()

	go poktStartScan()
}

func getSwanClient() *swan.SwanClient {
	var err error
	confPokt := GetConfig()

	swanApiUrl := confPokt.SwanApiUrl
	swanApiKey := confPokt.SwanApiKey
	swanAccessToken := confPokt.SwanAccessToken

	if len(swanApiUrl) == 0 {
		GetLog().Error("please set config-pokt->swan_api_url")
	}

	if len(swanApiKey) == 0 {
		GetLog().Error("please set config-pokt->swan_api_key")
	}

	if len(swanAccessToken) == 0 {
		GetLog().Error("please set config-pokt->swan_access_token")
	}

	swanClient, err := swan.GetClient(swanApiUrl, swanApiKey, swanAccessToken, "")
	if err != nil {
		GetLog().Error("Get Client Error: ", err)
	}

	return swanClient
}

func sendHeartbeat2Swan() {
	time.Sleep(time.Second * poktService.ApiHeartbeatInterval)
	//swanClient := getSwanClient()
	swanClient := &swan.SwanClient{}
	GetLog().Info("HeartBeats Start...")
	for {
		poktService.SendPoktHeartbeatRequest(swanClient)
		GetLog().Info("Sleeping...")
		time.Sleep(time.Second * poktService.ApiHeartbeatInterval)
	}
}

func poktStartScan() {
	GetLog().Info("Scan Start...")
	time.Sleep(time.Second * poktService.PoktScanInterval)

	for {
		poktService.StartScan()
		GetLog().Info("Sleeping...")
		//TODO: config
		time.Sleep(time.Second * poktService.PoktScanInterval)
	}
}

func poktHttpServer() {
	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()
	r.Use(cors.Middleware(cors.Config{
		Origins:         "*",
		Methods:         "GET, PUT, POST, DELETE",
		RequestHeaders:  "Origin, Authorization, Content-Type",
		ExposedHeaders:  "",
		MaxAge:          50 * time.Second,
		Credentials:     true,
		ValidateHeaders: false,
	}))

	apiv1 := r.Group(API_POCKET_V1)
	{
		apiv1.GET("/version", HttpGetPoktVersion)
		apiv1.GET("/height", HttpGetPoktCurHeight)
		apiv1.GET("/validator", HttpGetPoktValidatorAddr)
		apiv1.GET("/status", HttpGetPoktStatus)

		apiv1.POST("/balance", HttpGetPoktBalance)
		apiv1.POST("/threshold", HttpGetPoktThreshold)
		apiv1.POST("/set-validator", HttpSetPoktValidator)
		apiv1.POST("/custodial", HttpSetPoktCustodial)
		apiv1.POST("/noncustodial", HttpSetPoktNonCustodial)
	}

	port := GetConfig().PoktServerApiPort
	err := r.Run(":" + strconv.Itoa(port))
	if err != nil {
		GetLog().Fatal(err)
	}
}

func cmdPoktVersion() {
	params := ""
	confPokt := GetConfig()
	selfUrl := utils.UrlJoin(confPokt.PoktServerApiUrl, API_POCKET_V1)

	apiUrl := utils.UrlJoin(selfUrl, "version")
	response, err := web.HttpGetNoToken(apiUrl, params)
	if err != nil {
		fmt.Printf("Get Pocket Version err: %s \n", err)
		return
	}

	res := &VersionResponse{}
	err = json.Unmarshal(response, res)
	if err != nil {
		fmt.Printf("Parse Response (%s) err: %s \n", response, err)
		return
	}
	title := color.New(color.FgGreen).Sprintf("%s", "Pocket Version")
	value := color.New(color.FgYellow).Sprintf("%s", res.Data.Version)
	fmt.Printf("%s\t: %s\n", title, value)

}

func cmdPoktNodeAddr() {
	params := ""
	confPokt := GetConfig()
	selfUrl := utils.UrlJoin(confPokt.PoktServerApiUrl, API_POCKET_V1)

	apiUrl := utils.UrlJoin(selfUrl, "validator")
	response, err := web.HttpGetNoToken(apiUrl, params)
	if err != nil {
		fmt.Printf("Get Pocket Node Address err: %s \n", err)
		return
	}

	res := &Response{}
	err = json.Unmarshal(response, res)
	if err != nil {
		fmt.Printf("Parse Response (%s) err: %s \n", response, err)
		return
	}

	title := color.New(color.FgGreen).Sprintf("%s", "Validator Address")
	value := color.New(color.FgYellow).Sprintf("%s", res.Data)
	fmt.Printf("%s\t: %s\n", title, value)

}

func cmdPoktBalance(op []string) {
	fs := flag.NewFlagSet("Balance", flag.ExitOnError)
	addr := fs.String("addr", "", "address to lookup")
	err := fs.Parse(op[1:])
	if *addr == "" || err != nil {
		printPoktUsage()
		return
	}

	params := &BalancePoktParams{Height: 0, Address: *addr}
	confPokt := GetConfig()
	selfUrl := utils.UrlJoin(confPokt.PoktServerApiUrl, API_POCKET_V1)
	apiUrl := utils.UrlJoin(selfUrl, "balance")

	response, err := web.HttpPostNoToken(apiUrl, params)
	if err != nil {
		fmt.Printf("Get Pocket Balance err: %s params: %+v \n", err, params)
		return
	}

	res := &BalanceHttpResponse{}
	err = json.Unmarshal(response, res)
	if err != nil {
		fmt.Printf("Parse Response (%s) err: %s  params: %+v \n", response, err, params)
		return
	}

	title := color.New(color.FgGreen).Sprintf("%s", "Address")
	value := color.New(color.FgYellow).Sprintf("%s", res.Data.Address)
	fmt.Printf("%s\t: %s\n", title, value)

	title = color.New(color.FgGreen).Sprintf("%s", "Balance")
	value = color.New(color.FgYellow).Sprintf("%s", res.Data.Balance)
	fmt.Printf("%s\t: %s\n", title, value)

}

func cmdPoktCustodial(op []string) {

	fs := flag.NewFlagSet("custodial", flag.ExitOnError)
	fromAddr := fs.String("fromAddr", "", "")
	amount := fs.String("amount", "", "")
	relayChainIDs := fs.String("relayChainIDs", "", "")
	serviceURI := fs.String("serviceURI", "", "")
	networkID := fs.String("networkID", "", "")
	fee := fs.String("fee", "", "")
	isBefore := fs.String("isBefore", "", "")
	passwd := fs.String("passwd", "", "")

	err := fs.Parse(op[1:])
	if *fromAddr == "" || *amount == "" || *relayChainIDs == "" || *serviceURI == "" || *networkID == "" || *fee == "" || *isBefore == "" || err != nil {
		printPoktUsage()
		return
	}

	params := &CustodialParams{
		Address:       *fromAddr,
		Amount:        *amount,
		RelayChainIDs: *relayChainIDs,
		ServiceURI:    *serviceURI,
		NetworkID:     *networkID,
		Fee:           *fee,
		IsBefore:      *isBefore,
		Passwd:        *passwd,
	}

	confPokt := GetConfig()
	selfUrl := utils.UrlJoin(confPokt.PoktServerApiUrl, API_POCKET_V1)

	apiUrl := utils.UrlJoin(selfUrl, "custodial")
	response, err := web.HttpPostNoToken(apiUrl, params)
	if err != nil {
		fmt.Printf("Get Pocket Custodial err: %s \n", err)
		return
	}

	res := &CustodialResponse{}
	err = json.Unmarshal(response, res)
	if err != nil {
		fmt.Printf("Parse Response (%s) err: %s \n", response, err)
		return
	}

	fmt.Printf("Pocket Custodial Result is: %+v \n", res.Data)
}

func cmdPoktNonCustodial(op []string) {

	fs := flag.NewFlagSet("non-custodial", flag.ExitOnError)
	operatorPublicKey := fs.String("operatorPublicKey", "", "")
	outputAddress := fs.String("outputAddress", "", "")
	amount := fs.String("amount", "", "")
	relayChainIDs := fs.String("relayChainIDs", "", "")
	serviceURI := fs.String("serviceURI", "", "")
	networkID := fs.String("networkID", "", "")
	fee := fs.String("fee", "", "")
	isBefore := fs.String("isBefore", "", "")

	err := fs.Parse(op[1:])
	if *operatorPublicKey == "" || *outputAddress == "" || *amount == "" || *relayChainIDs == "" || *serviceURI == "" || *networkID == "" || *fee == "" || *isBefore == "" || err != nil {
		printPoktUsage()
		return
	}

	params := ""
	confPokt := GetConfig()
	selfUrl := utils.UrlJoin(confPokt.PoktServerApiUrl, API_POCKET_V1)

	apiUrl := utils.UrlJoin(selfUrl, "nonCustodial")
	response, err := web.HttpPostNoToken(apiUrl, params)
	if err != nil {
		fmt.Printf("Get Pocket NonCustodial err: %s \n", err)
		return
	}

	res := &StatusResponse{}
	err = json.Unmarshal(response, res)
	if err != nil {
		fmt.Printf("Parse Response (%s) err: %s \n", response, err)
		return
	}

	fmt.Printf("Pocket NonCustodial is: %+v \n", res)
}

func cmdPoktStatus() {
	params := ""
	confPokt := GetConfig()
	selfUrl := utils.UrlJoin(confPokt.PoktServerApiUrl, API_POCKET_V1)

	apiUrl := utils.UrlJoin(selfUrl, "status")
	response, err := web.HttpGetNoToken(apiUrl, params)
	if err != nil {
		fmt.Printf("Get Pocket Status err: %s \n", err)
		return
	}

	res := &StatusResponse{}
	err = json.Unmarshal(response, res)
	if err != nil {
		fmt.Printf("Parse Response (%s) err: %s \n", response, err)
		return
	}

	title := color.New(color.FgGreen).Sprintf("%s", "Version")
	value := color.New(color.FgYellow).Sprintf("%s", res.Data.Version)
	fmt.Printf("%s\t\t: %s\n", title, value)

	title = color.New(color.FgGreen).Sprintf("%s", "Height")
	value = color.New(color.FgYellow).Sprintf("%d", res.Data.Height)
	fmt.Printf("%s\t\t: %s\n", title, value)

	title = color.New(color.FgGreen).Sprintf("%s", "Synced")
	value = color.New(color.FgRed).Sprintf("%t", res.Data.Synced)
	fmt.Printf("%s\t\t: %s\n", title, value)

	title = color.New(color.FgGreen).Sprintf("%s", "Address")
	value = color.New(color.FgYellow).Sprintf("%s", res.Data.Address)
	fmt.Printf("%s\t\t: %s\n", title, value)

	title = color.New(color.FgGreen).Sprintf("%s", "PublicKey")
	value = color.New(color.FgYellow).Sprintf("%s", res.Data.PublicKey)
	fmt.Printf("%s\t: %s\n", title, value)

	title = color.New(color.FgGreen).Sprintf("%s", "Balance")
	value = color.New(color.FgYellow).Sprintf("%d", res.Data.Balance)
	fmt.Printf("%s\t\t: %s\n", title, value)

	title = color.New(color.FgGreen).Sprintf("%s", "Staking")
	value = color.New(color.FgYellow).Sprintf("%s", res.Data.Staking)
	fmt.Printf("%s\t\t: %s\n", title, value)

	title = color.New(color.FgGreen).Sprintf("%s", "Jailed")
	value = color.New(color.FgRed).Sprintf("%t", res.Data.Jailed)
	fmt.Printf("%s\t\t: %s\n", title, value)

	title = color.New(color.FgGreen).Sprintf("%s", "JailedBlock")
	value = color.New(color.FgRed).Sprintf("%d", res.Data.JailedBlock)
	fmt.Printf("%s\t: %s\n", title, value)

	title = color.New(color.FgGreen).Sprintf("%s", "JailedUntil")
	value = color.New(color.FgRed).Sprintf("%s", res.Data.JailedUntil)
	fmt.Printf("%s\t: %s\n", title, value)
}
