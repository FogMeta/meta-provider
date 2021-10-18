package service

import (
	"strings"
	"swan-provider/common/client"
	"swan-provider/logs"
	"time"
)

const ARIA2_TASK_STATUS_ERROR = "error"
const ARIA2_TASK_STATUS_ACTIVE = "active"
const ARIA2_TASK_STATUS_COMPLETE = "complete"

const DEAL_STATUS_CREATED = "Created"
const DEAL_STATUS_WAITING = "Waiting"

const DEAL_STATUS_DOWNLOADING = "Downloading"
const DEAL_STATUS_DOWNLOADED = "Downloaded"
const DEAL_STATUS_DOWNLOAD_FAILED = "DownloadFailed"

const DEAL_STATUS_IMPORT_READY = "ReadyForImport"
const DEAL_STATUS_IMPORTING = "FileImporting"
const DEAL_STATUS_IMPORTED = "FileImported"
const DEAL_STATUS_IMPORT_FAILED = "ImportFailed"
const DEAL_STATUS_ACTIVE = "DealActive"

const ONCHAIN_DEAL_STATUS_ERROR = "StorageDealError"
const ONCHAIN_DEAL_STATUS_ACTIVE = "StorageDealActive"
const ONCHAIN_DEAL_STATUS_NOTFOUND = "StorageDealNotFound"
const ONCHAIN_DEAL_STATUS_WAITTING = "StorageDealWaitingForData"
const ONCHAIN_DEAL_STATUS_ACCEPT = "StorageDealAcceptWait"
const ONCHAIN_DEAL_STATUS_AWAITING = "StorageDealAwaitingPreCommit"

const ARIA2_MAX_DOWNLOADING_TASKS = 10
const LOTUS_IMPORT_NUMNBER = "20" //Max number of deals to be imported at a time
const LOTUS_SCAN_NUMBER = "100"   //Max number of deals to be scanned at a time

var aria2Client = client.GetAria2Client()
var swanClient = client.GetSwanClient()

var swanService = GetSwanService()
var aria2Service = GetAria2Service()
var lotusService = GetLotusService()

func AdminOfflineDeal() {
	checkLotusConfig()
	swanService.UpdateBidConf(swanClient)
	go swanSendHeartbeatRequest()
	go aria2CheckDownloadStatus()
	go aria2StartDownload()
	go lotusStartImport()
	go lotusStartScan()
}

func checkLotusConfig() {
	logs.GetLogger().Info("Start testing lotus config.")

	lotusClient := client.LotusGetClient()
	if len(lotusClient.ApiUrl) == 0 {
		logs.GetLogger().Fatal("please set config:lotus->api_url")
	}

	if len(lotusClient.MinerApiUrl) == 0 {
		logs.GetLogger().Fatal("please set config:lotus->miner_api_url")
	}

	if len(lotusClient.MinerAccessToken) == 0 {
		logs.GetLogger().Fatal("please set config:lotus->miner_access_token")
	}

	response := client.LotusImportData("bafyreib7azyg2yubucdhzn64gvyekdma7nbrbnfafcqvhsz2mcnvbnkktu", "test")

	if strings.Contains(response, "no return") {
		logs.GetLogger().Fatal("please check config:lotus->miner_api_url,lotus->miner_access_token")
	}

	if strings.Contains(response, "(need 'write')") {
		logs.GetLogger().Error("please check config:lotus->miner_access_token")
		logs.GetLogger().Fatal(response)
	}

	currentEpoch := client.LotusGetCurrentEpoch()
	if currentEpoch < 0 {
		logs.GetLogger().Fatal("please check config:lotus->api_url")
	}

	logs.GetLogger().Info("Pass testing lotus config.")
}

func swanSendHeartbeatRequest() {
	for {
		logs.GetLogger().Info("Start...")
		swanService.SendHeartbeatRequest(swanClient)
		logs.GetLogger().Info("Sleeping...")
		time.Sleep(swanService.ApiHeartbeatInterval)
	}
}

func aria2CheckDownloadStatus() {
	for {
		logs.GetLogger().Info("Start...")
		aria2Service.CheckDownloadStatus(aria2Client, swanClient)
		logs.GetLogger().Info("Sleeping...")
		time.Sleep(time.Minute)
	}
}

func aria2StartDownload() {
	for {
		logs.GetLogger().Info("Start...")
		aria2Service.StartDownload(aria2Client, swanClient)
		logs.GetLogger().Info("Sleeping...")
		time.Sleep(time.Minute)
	}
}

func lotusStartImport() {
	for {
		logs.GetLogger().Info("Start...")
		lotusService.StartImport(swanClient)
		logs.GetLogger().Info("Sleeping...")
		time.Sleep(lotusService.ImportIntervalSecond)
	}
}

func lotusStartScan() {
	for {
		logs.GetLogger().Info("Start...")
		lotusService.StartScan(swanClient)
		logs.GetLogger().Info("Sleeping...")
		time.Sleep(lotusService.ScanIntervalSecond)
	}
}