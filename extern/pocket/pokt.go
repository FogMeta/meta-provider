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
	"math/big"
	"meta-provider/common"
	"net/http"
	"os"
	"strconv"
	"time"
)

const API_POCKET_V1 = "/poktsrv"

type Status struct {
	Version     string
	Address     string
	Height      string
	Balance     string
	Award       string
	Jailed      string
	JailedBlock string
	JailedUntil string
}

type PoktService struct {
	SwanApiUrl           string
	SwanApiKey           string
	SwanAccessToken      string
	PoktApiUrl           string
	PoktScanInterval     time.Duration
	ApiHeartbeatInterval time.Duration
	PoktServerApiUrl     string
	PoktServerApiPort    int

	dkImage    string
	dkName     string
	dkConfPath string

	dkNetworkType string

	dkCli *DockerCli

	PoktAddress string
	AlarmBlc    big.Int
	CurStatus   Status
}

var myPoktSvr *PoktService

func GetMyPoktService() *PoktService {
	if myPoktSvr == nil {
		confPokt := GetConfig()

		myPoktSvr = &PoktService{
			SwanApiUrl:           confPokt.SwanApiUrl,
			SwanApiKey:           confPokt.SwanApiKey,
			SwanAccessToken:      confPokt.SwanAccessToken,
			PoktApiUrl:           confPokt.PoktApiUrl,
			PoktScanInterval:     confPokt.PoktScanInterval,
			ApiHeartbeatInterval: confPokt.PoktHeartbeatInterval,
			PoktServerApiUrl:     confPokt.PoktServerApiUrl,
			PoktServerApiPort:    confPokt.PoktServerApiPort,
			dkImage:              confPokt.PoktDockerImage,
			dkName:               confPokt.PoktDockerName,
			dkConfPath:           confPokt.PoktConfigPath,
			dkNetworkType:        confPokt.PoktNetworkType,
			CurStatus:            Status{},
		}
		// 新目录名称和权限属性
		perm := os.FileMode(0777)
		err := os.MkdirAll(myPoktSvr.dkConfPath, perm)
		if err != nil {
			GetLog().Error("Create ", myPoktSvr.dkConfPath, "error: ", err)
			panic("Create Pocket Data Path Error")
		}
		os.Chmod(myPoktSvr.dkConfPath, perm)

		myPoktSvr.dkCli = GetMyCli(myPoktSvr.dkImage, myPoktSvr.dkName, myPoktSvr.dkConfPath)

		GetLog().Debugf("New myPoktSvr :%+v ", *myPoktSvr)

		return myPoktSvr
	}
	return myPoktSvr
}

func (psvc *PoktService) GetCli() *DockerCli {
	if psvc.dkCli == nil {
		psvc.dkCli = GetMyCli(psvc.dkImage, psvc.dkName, psvc.dkConfPath)
		GetLog().Infof("GetCli New Docker Cli")
	}
	return psvc.dkCli
}

func printPoktUsage() {
	fmt.Println("SUBCMD:")
	fmt.Println("    pocket")
	fmt.Println("USAGE:")
	fmt.Println("    swan-provider pocket start")
	fmt.Println("    swan-provider pocket version")
	fmt.Println("    swan-provider pocket validator")
	fmt.Println("    swan-provider pocket balance --addr=0123456789012345678901234567890123456789")
	fmt.Println("    swan-provider pocket status")
	//fmt.Println("    swan-provider pocket custodial")
	//fmt.Println("                         --fromAddr='0123456789012345678901234567890123456789'")
	//fmt.Println("                         --amount='15100000000'")
	//fmt.Println("                         --relayChainIDs='0001,0021'")
	//fmt.Println("                         --serviceURI='https://pokt.rocks:443'")
	//fmt.Println("                         --networkID='mainnet'")
	//fmt.Println("                         --fee='10000'")
	//fmt.Println("                         --isBefore='false'")
	//fmt.Println("    swan-provider pocket noncustodial")
	//fmt.Println("                         --operatorPublicKey='0123456789012345678901234567890123456789012345678901234567890123'")
	//fmt.Println("                         --outputAddress='0123456789012345678901234567890123456789'")
	//fmt.Println("                         --amount='15100000000'")
	//fmt.Println("                         --relayChainIDs='0001,0021'")
	//fmt.Println("                         --serviceURI='https://pokt.rocks:443'")
	//fmt.Println("                         --networkID='mainnet'")
	//fmt.Println("                         --fee='10000'")
	//fmt.Println("                         --isBefore='false'")

}

func (psvc *PoktService) StartPoktContainer(op []string) {

	cli := psvc.dkCli
	if !cli.PoktCtnExist() {

		GetLog().Debug("Init Pocket Container ... ")
		fs := flag.NewFlagSet("Start", flag.ExitOnError)
		passwd := fs.String("passwd", "", "password for create account")
		err := fs.Parse(op[1:])
		if *passwd == "" || err != nil {
			printPoktUsage()
			panic("need password for create account.")
			return
		}

		pass := *passwd
		GetLog().Debug("POCKET_PASSPHRASE=", pass)
		env := []string{"POCKET_PASSPHRASE=" + pass}

		//accCmd := []string{"pocket", "accounts", "create"}
		accCmd := []string{""}
		cli.PoktCtnPullAndCreate(accCmd, env, true)
		cli.PoktCtnStart()

		for {
			if cli.PoktCtnExist() {
				GetLog().Info("Wait for Creating Account...")
				cli.PoktCtnList()
				time.Sleep(time.Second * 3)
				continue
			}
			break
		}
		GetLog().Debug("Init Creating Account Over")

		runCmd := []string{""}
		if psvc.dkNetworkType == "TESTNET" {
			env = []string{"POCKET_TESTNET='true'"}
		} else if psvc.dkNetworkType == "MAINNET" {
			env = []string{"POCKET_MAINNET='true'"}
		} else if psvc.dkNetworkType == "SIMULATE" {
			env = []string{"POCKET_SIMULATE='true'"}
		}

		chains := ReadPocketChains()
		if chains != "" {
			env = append(env, "POCKET_CORE_CHAINS="+chains)
		}

		GetLog().Info("Create Pocket ", psvc.dkNetworkType, "")

		cli.PoktCtnCreateRun(runCmd, env, false)

	}

	if !cli.PoktCtnStart() {
		GetLog().Error("Pocket Start FALSE")
	}
}

func (psvc *PoktService) StartScan() {
	url := psvc.PoktApiUrl
	height, err := PoktApiGetHeight(url)
	if err != nil {
		GetLog().Error(err)
	}

	synced, err := PoktApiGetSync()
	if err != nil {
		GetLog().Error(err)
	}
	GetLog().Info("Pokt Get Current Height=", height, " and Synced=", synced)

}

func (psvc *PoktService) SendPoktHeartbeatRequest(swanClient *swan.SwanClient) {

	params := ""
	confPokt := GetConfig()
	selfUrl := utils.UrlJoin(confPokt.PoktServerApiUrl, API_POCKET_V1)

	apiUrl := utils.UrlJoin(selfUrl, "status")
	response, err := web.HttpGetNoToken(apiUrl, params)
	if err != nil {
		fmt.Printf("Heartbeat Get Pocket Status err: %s \n", err)
		return
	}

	res := &StatusResponse{}
	err = json.Unmarshal(response, res)
	if err != nil {
		fmt.Printf("Heartbeat Parse Response (%s) err: %s \n", response, err)
		return
	}
	// Swan Server Is Not Ready!
	//stat := swan.PocketHeartbeatOnlineParams{
	//	Address:     res.Data.Address,
	//	Version:     res.Data.Version,
	//	Height:      res.Data.Height,
	//	Balance:     res.Data.Balance,
	//	Award:       res.Data.Award,
	//	Jailed:      res.Data.Jailed,
	//	JailedBlock: res.Data.JailedBlock,
	//	JailedUntil: res.Data.JailedUntil,
	//}
	//err = swanClient.SendPoktHeartbeatRequest(stat)
	//if err != nil {
	//	fmt.Printf("Heartbeat Send err: %s \n", err)
	//	return
	//}

	{
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

	return
}

func HttpGetPoktVersion(c *gin.Context) {
	poktSvr := GetMyPoktService()
	cmdOut, err := poktSvr.GetCli().PoktCtnExecVersion()
	if err != nil {
		GetLog().Error(err)
		c.JSON(http.StatusInternalServerError, common.CreateErrorResponse("-1", err.Error()))
		return
	}
	c.JSON(http.StatusOK, common.CreateSuccessResponse(cmdOut))
}

func HttpGetPoktCurHeight(c *gin.Context) {
	poktSvr := GetMyPoktService()
	cmdOut, err := poktSvr.GetCli().PoktCtnExecHeight()
	if err != nil {
		GetLog().Error(err)
		c.JSON(http.StatusInternalServerError, common.CreateErrorResponse("-1", err.Error()))
		return
	}
	c.JSON(http.StatusOK, common.CreateSuccessResponse(cmdOut))
}

func HttpGetPoktValidatorAddr(c *gin.Context) {
	poktSvr := GetMyPoktService()
	cmdOut, err := poktSvr.GetCli().PoktCtnExecValidatorAddress()
	if err != nil {
		GetLog().Error(err)
		c.JSON(http.StatusInternalServerError, common.CreateErrorResponse("-1", err.Error()))
		return
	}
	c.JSON(http.StatusOK, common.CreateSuccessResponse(cmdOut))
}

func HttpGetPoktStatus(c *gin.Context) {
	poktSvr := GetMyPoktService()
	data := &StatusData{}

	versionData, err := poktSvr.GetCli().PoktCtnExecVersion()
	if err != nil {
		GetLog().Error(err)
	} else {
		data.Version = versionData.Version
	}

	heightData, err := poktSvr.GetCli().PoktCtnExecHeight()
	if err != nil {
		GetLog().Error(err)
	} else {
		data.Height = heightData.Height
	}

	data.Synced, _ = PoktApiGetSync()

	//address, err := poktSvr.GetCli().PoktCtnExecNodeAddress()
	address, err := poktSvr.GetCli().PoktCtnExecInitAddress()
	if err != nil {
		GetLog().Error(err)
	} else {
		data.Address = address
	}

	balanceData, err := poktSvr.GetCli().PoktCtnExecBalance(address)
	if err != nil {
		GetLog().Error(err)
	} else {
		data.Balance = balanceData.Balance
	}

	nodeData, err := poktSvr.GetCli().PoktCtnExecNode(address)
	if err != nil {
		GetLog().Error(err)
	} else {
		data.Staking = nodeData.StakedTokens
		data.PublicKey = nodeData.PublicKey
		data.Jailed = nodeData.Jailed
	}

	signData, err := poktSvr.GetCli().PoktCtnExecSignInfo(address)
	if err != nil || len(signData) == 0 {
		GetLog().Warn(err)
	} else {
		signInfo := signData[0]
		data.JailedUntil = signInfo.JailedUntil
		data.JailedBlock = signInfo.JailedBlocksCounter
	}

	c.JSON(http.StatusOK, common.CreateSuccessResponse(data))
}

func HttpGetPoktBalance(c *gin.Context) {
	var params BalancePoktParams
	err := c.BindJSON(&params)
	if err != nil {
		GetLog().Error(err)
		c.JSON(http.StatusOK, common.CreateErrorResponse("-1", err.Error()))
		return
	}

	poktSvr := GetMyPoktService()
	cmdOut, err := poktSvr.GetCli().PoktCtnExecBalance(params.Address)
	if err != nil {
		GetLog().Error(err)
		c.JSON(http.StatusInternalServerError, common.CreateErrorResponse("-1", err.Error()))
		return
	}

	GetLog().Debug("pocket query balance result:", cmdOut)

	data := &BalanceCmdData{
		Height:  params.Height,
		Address: params.Address,
		Balance: strconv.FormatUint(cmdOut.Balance, 10)}
	c.JSON(http.StatusOK, common.CreateSuccessResponse(data))
}

func HttpGetPoktThreshold(c *gin.Context) {
	var params ThresholdParams
	err := c.BindJSON(&params)
	if err != nil {
		GetLog().Error(err)
		c.JSON(http.StatusOK, common.CreateErrorResponse("-1", err.Error()))
		return
	}

	data := &ThresholdData{
		Address:   params.Address,
		Threshold: params.Threshold,
		Active:    true,
	}
	c.JSON(http.StatusOK, common.CreateSuccessResponse(data))
}

///////////////////////////////////////////////////////////////////////////////

func HttpSetPoktValidator(c *gin.Context) {
	var params ValidatorParams
	err := c.BindJSON(&params)
	if err != nil {
		GetLog().Error(err)
		c.JSON(http.StatusOK, common.CreateErrorResponse("-1", err.Error()))
		return
	}

	poktSvr := GetMyPoktService()
	result, err := poktSvr.GetCli().PoktCtnExecSetValidator(
		params.Address,
		params.Passwd,
	)
	if err != nil {
		GetLog().Error(err)
		c.JSON(http.StatusInternalServerError, common.CreateErrorResponse("-1", err.Error()))
		return
	}

	data := &ValidatorData{
		Result: result,
	}
	c.JSON(http.StatusOK, common.CreateSuccessResponse(data))
}

///////////////////////////////////////////////////////////////////////////////

func HttpSetPoktCustodial(c *gin.Context) {
	var params CustodialParams
	err := c.BindJSON(&params)
	if err != nil {
		GetLog().Error(err)
		c.JSON(http.StatusOK, common.CreateErrorResponse("-1", err.Error()))
		return
	}

	poktSvr := GetMyPoktService()
	result, err := poktSvr.GetCli().PoktCtnExecCustodial(
		params.Address,
		params.Amount,
		params.RelayChainIDs,
		params.ServiceURI,
		params.NetworkID,
		params.Fee,
		params.IsBefore,
		params.Passwd,
	)
	if err != nil {
		GetLog().Error(err)
		c.JSON(http.StatusInternalServerError, common.CreateErrorResponse("-1", err.Error()))
		return
	}

	data := &CustodialData{
		Result: result,
	}
	c.JSON(http.StatusOK, common.CreateSuccessResponse(data))
}

///////////////////////////////////////////////////////////////////////////////

func HttpSetPoktNonCustodial(c *gin.Context) {
	var params NonCustodialParams
	err := c.BindJSON(&params)
	if err != nil {
		GetLog().Error(err)
		c.JSON(http.StatusOK, common.CreateErrorResponse("-1", err.Error()))
		return
	}

	poktSvr := GetMyPoktService()
	result, err := poktSvr.GetCli().PoktCtnExecNonCustodial(
		params.PubKey,
		params.OutputAddr,
		params.Amount,
		params.RelayChainIDs,
		params.ServiceURI,
		params.NetworkID,
		params.Fee,
		params.IsBefore,
		params.Passwd,
	)
	if err != nil {
		GetLog().Error(err)
		c.JSON(http.StatusInternalServerError, common.CreateErrorResponse("-1", err.Error()))
		return
	}

	data := &NonCustodialData{
		Result: result,
	}
	c.JSON(http.StatusOK, common.CreateSuccessResponse(data))
}
