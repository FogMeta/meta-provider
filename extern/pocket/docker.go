package pocket

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	dc "github.com/docker/docker/client"
	"github.com/filswan/go-swan-lib/logs"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"time"
)

const (
	POCKET_CONFIG_PATH = "/home/app"
)

var myCli *DockerCli

type DockerCli struct {
	Image    string
	Name     string
	DataPath string

	Client *dc.Client
	Ctx    context.Context

	Cid string
}

func GetMyCli(image string, name string, path string) *DockerCli {

	if myCli == nil {
		cli := &DockerCli{
			Image:    image,
			Name:     name,
			DataPath: path,
		}

		client, err := dc.NewClientWithOpts(dc.FromEnv)
		if err != nil {
			GetLog().Error(err)
			return nil
		}
		cli.Client = client
		cli.Ctx = context.Background()

		clist, err := cli.Client.ContainerList(cli.Ctx, types.ContainerListOptions{All: true})
		if err != nil {
			GetLog().Error(err)
			return nil
		}

		finded := false
		for _, container := range clist {
			if "/"+cli.Name == container.Names[0] {
				cli.Cid = container.ID
				finded = true
				break
			}
		}
		if !finded {
			//
		}

		myCli = cli
		return myCli
	}

	return myCli
}

func GetDockerCli(image string, name string, path string) *DockerCli {
	dockerCli := &DockerCli{
		Image:    image,
		Name:     name,
		DataPath: path,
	}

	client, err := dc.NewClientWithOpts(dc.FromEnv)
	if err != nil {
		GetLog().Error(err)
		return nil
	}

	dockerCli.Client = client
	dockerCli.Ctx = context.Background()
	dockerCli.Cid = ""

	return dockerCli
}

func (cli *DockerCli) UpdateCid() (bool, error) {

	clist, err := cli.Client.ContainerList(cli.Ctx, types.ContainerListOptions{All: true})
	if err != nil {
		GetLog().Error(err)
		return false, err
	}

	found := false
	for _, container := range clist {
		if "/"+cli.Name == container.Names[0] {
			cli.Cid = container.ID
			found = true
			break
		}
	}

	if !found {
		return false, errors.New("do not find container")
	}

	return true, nil
}

func (cli *DockerCli) PoktCtnCreate() bool {
	out, err := cli.Client.ImagePull(cli.Ctx, cli.Image, types.ImagePullOptions{})
	if err != nil {
		GetLog().Error("Image Pull Error:", err)
		return false
	}
	defer out.Close()
	io.Copy(os.Stdout, out)
	//GetLog().Warn("### Image Pull Skip ###")

	body, err := cli.Client.ContainerCreate(
		cli.Ctx,
		&container.Config{Image: cli.Image, Tty: true},
		&container.HostConfig{NetworkMode: "host", Binds: []string{cli.DataPath + ":/home/app/.pocket"}},
		nil,
		nil,
		cli.Name)
	if err != nil {
		GetLog().Error("Container Create Error:", err)
		return false
	}

	GetLog().Debug("Container Create Id:", body.ID[:10])
	cli.Cid = body.ID
	return true
}

func (cli *DockerCli) PoktCtnPullAndCreate(cmd, env []string, autoRemove bool) bool {
	out, err := cli.Client.ImagePull(cli.Ctx, cli.Image, types.ImagePullOptions{})
	if err != nil {
		GetLog().Error("Image Pull Error:", err)
		return false
	}
	defer out.Close()
	io.Copy(os.Stdout, out)
	//GetLog().Warn("### Image Pull Skip ###")

	GetLog().Debug("Container Create Para DataPath=", cli.DataPath, " autoRemove=", autoRemove)
	body, err := cli.Client.ContainerCreate(
		cli.Ctx,
		&container.Config{
			Cmd:   cmd,
			Env:   env,
			Image: cli.Image,
			Tty:   true},
		&container.HostConfig{
			NetworkMode: "host",
			Binds:       []string{cli.DataPath + ":/home/app/.pocket"},
			AutoRemove:  autoRemove},
		nil,
		nil,
		cli.Name)
	if err != nil {
		GetLog().Error("Container Create Error:", err)
		return false
	}

	GetLog().Info("Account Container Create Id:", body.ID[:10])
	cli.Cid = body.ID
	return true
}

func (cli *DockerCli) PoktCtnCreateRun(cmd, env []string, autoRemove bool) bool {

	body, err := cli.Client.ContainerCreate(
		cli.Ctx,
		&container.Config{
			Cmd:   cmd,
			Env:   env,
			Image: cli.Image,
			Tty:   true},
		&container.HostConfig{
			NetworkMode: "host",
			Binds:       []string{cli.DataPath + ":/home/app/.pocket"},
			AutoRemove:  autoRemove},
		nil,
		nil,
		cli.Name)
	if err != nil {
		GetLog().Error("Container Create Error:", err)
		return false
	}

	GetLog().Debug("Container Create Id:", body.ID[:10])
	cli.Cid = body.ID
	return true
}

func (cli *DockerCli) PoktCtnExist() bool {
	clist, err := cli.Client.ContainerList(cli.Ctx, types.ContainerListOptions{All: true})
	if err != nil {
		GetLog().Error(err)
		return false
	}

	for _, container := range clist {
		if "/"+cli.Name == container.Names[0] {
			cli.Cid = container.ID
			return true
		}
	}

	GetLog().Debug("Can Not Find Container:", " Name=", cli.Name)
	return false
}

func (cli *DockerCli) PoktCtnList() bool {
	clist, err := cli.Client.ContainerList(cli.Ctx, types.ContainerListOptions{All: true})
	if err != nil {
		GetLog().Error(err)
		return false
	}

	for _, container := range clist {
		GetLog().Info("Container Create Name:", container.Names[0], " ID=", container.ID[:10])
	}
	return true
}

func (cli *DockerCli) PoktCtnStart() bool {

	containers, err := cli.Client.ContainerList(cli.Ctx, types.ContainerListOptions{All: true})
	if err != nil {
		GetLog().Error(err)
		return false
	}

	for _, container := range containers {
		if "/"+cli.Name == container.Names[0] {
			cli.Cid = container.ID

			if !strings.Contains(container.Status, "Up") {
				err := cli.Client.ContainerStart(context.Background(), cli.Cid, types.ContainerStartOptions{})
				if err != nil {
					GetLog().Error("Container Start Error:", err)
					return false
				}
				GetLog().Debug("Container start:", " id=", cli.Cid, " name=", cli.Name)
				return true
			} else {

				GetLog().Debug("Container Already Running:", " id=", cli.Cid[:10], " name=", cli.Name, " status=", container.Status)
				return true
			}

			break
		}
	}

	return false
}

func (cli *DockerCli) PoktCtnStop() bool {
	timeout := time.Second * 5
	err := cli.Client.ContainerStop(cli.Ctx, cli.Cid, &timeout)
	if err != nil {
		GetLog().Error("Container Stop Error:", err)
		return false
	}

	GetLog().Debug("Stop Container Id:", cli.Cid[:10])
	return true
}

func (cli *DockerCli) PoktCtnExec(cmd []string) (string, error) {

	rst, err := cli.Client.ContainerExecCreate(
		cli.Ctx, cli.Cid,
		types.ExecConfig{AttachStdout: true, AttachStderr: true, Cmd: cmd})
	if err != nil {
		GetLog().Error("Container Exec Create Error:", err)
		return "", err
	}

	response, err := cli.Client.ContainerExecAttach(cli.Ctx, rst.ID, types.ExecStartCheck{})
	if err != nil {
		GetLog().Error("Container Exec Attach Error:", err)
		return "", err
	}
	defer response.Close()

	data, _ := ioutil.ReadAll(response.Reader)
	GetLog().Debug("Container Exec Response:", string(data))

	return string(data), nil
}

func (cli *DockerCli) PoktCtnExecVersion() (*VersionData, error) {
	strRes, err := cli.PoktCtnExec([]string{"pocket", "version"})
	if err != nil {
		GetLog().Error("Exec Pocket Version Error:", err)
		return &VersionData{}, err
	}

	GetLog().Info("Exec Pocket Version:", strRes)
	index := strings.Index(strRes, ":")
	if index < 0 {
		GetLog().Error("Exec Pocket Version Error: No version info return")
		return &VersionData{}, errors.New("No Version Info ")
	}

	jOut := &VersionData{
		Version: strings.TrimSuffix(strRes[index+2:], "\n"),
	}
	GetLog().Debug("pocket query version result:", jOut.Version)

	return jOut, nil
}

func (cli *DockerCli) PoktCtnExecInitAddress() (string, error) {
	strRes, err := cli.PoktCtnExec([]string{"pocket", "accounts", "list"})
	if err != nil {
		GetLog().Error("Exec Pocket Node Account Error:", err)
		return "", err
	}

	index := strings.Index(strRes, ")")
	index += 2
	jOut := strRes[index : index+40]
	GetLog().Debug("pocket query node account result:", jOut)

	_, finded := os.LookupEnv("TEST_POCKET_MODE")
	if finded {
		//ONLY FOR TEST
		jOut = "ffad090789253ad0439c56b7b9c301f90424d5b7"
	}

	return jOut, nil
}

func (cli *DockerCli) PoktCtnExecValidatorAddress() (string, error) {
	strRes, err := cli.PoktCtnExec([]string{"pocket", "accounts", "get-validator"})
	if err != nil {
		GetLog().Error("Exec Pocket Node Account Error:", err)
		return "", err
	}

	index := strings.Index(strRes, "Validator Address:") + len("Validator Address:")
	jOut := strRes[index : index+40]
	GetLog().Debug("pocket query validator address result:", jOut)

	_, finded := os.LookupEnv("TEST_POCKET_MODE")
	if finded {
		//ONLY FOR TEST
		jOut = "ffad090789253ad0439c56b7b9c301f90424d5b7"
	}

	return jOut, nil
}

func (cli *DockerCli) PoktCtnExecSetAccount(address string) (string, error) {
	rsp, err := cli.PoktCtnExec([]string{"pocket", "accounts", "set-validator", address})
	if err != nil {
		GetLog().Error("Exec Pocket Set Account Error:", err)
		return "", err
	}
	return rsp, nil
}

func (cli *DockerCli) PoktCtnExecHeight() (*HeightData, error) {
	jOut := &HeightData{}
	strRes, err := cli.PoktCtnExec([]string{"pocket", "query", "height"})
	if err != nil {
		GetLog().Error("Exec Pocket Height Error:", err)
		return jOut, err
	}

	index := strings.Index(strRes, "{")
	GetLog().Debug("pocket query height result for json:", strRes[index:])
	err = json.Unmarshal([]byte(strRes[index:]), jOut)
	if err != nil {
		logs.GetLogger().Error(err)
		return jOut, err
	}

	return jOut, nil
}

func (cli *DockerCli) PoktCtnExecBalance(address string) (*BalanceData, error) {
	jOut := &BalanceData{}
	strRes, err := cli.PoktCtnExec([]string{"pocket", "query", "balance", address})
	if err != nil {
		GetLog().Error("Exec Pocket Balance Error:", err)
		return jOut, err
	}

	index := strings.Index(strRes, "{")
	GetLog().Debug("pocket query balance result for json:", strRes[index:])
	err = json.Unmarshal([]byte(strRes[index:]), jOut)
	if err != nil {
		logs.GetLogger().Error(err)
		return jOut, err
	}

	return jOut, nil
}

func (cli *DockerCli) PoktCtnExecSignInfo(address string) ([]*SignInfo, error) {
	jOut := &SignInfoResponse{}
	strRes, err := cli.PoktCtnExec([]string{"pocket", "query", "signing-info", address})
	if err != nil {
		GetLog().Error("Exec Pocket Sign Info Error:", err)
		return jOut.Result, err
	}

	index := strings.Index(strRes, "{")
	GetLog().Debug("pocket query signing-info result for json:", strRes[index:])
	err = json.Unmarshal([]byte(strRes[index:]), jOut)
	if err != nil {
		logs.GetLogger().Error(err)
		return jOut.Result, err
	}

	return jOut.Result, nil
}

func (cli *DockerCli) PoktCtnExecNode(address string) (*NodeData, error) {
	jOut := &NodeData{}
	strRes, err := cli.PoktCtnExec([]string{"pocket", "query", "node", address})
	if err != nil {
		GetLog().Error("Exec Pocket Node Error:", err)
		return jOut, err
	}

	index := strings.Index(strRes, "{")
	GetLog().Debug("pocket query node result for json:", strRes[index:])
	err = json.Unmarshal([]byte(strRes[index:]), jOut)
	if err != nil {
		logs.GetLogger().Error(err)
		return jOut, err
	}

	return jOut, nil
}

func (cli *DockerCli) PoktCtnExecSetValidator(address, passwd string) (string, error) {
	rsp, err := cli.PoktCtnExec([]string{"expect", POCKET_CONFIG_PATH + "/set-validator.sh", address, passwd})
	if err != nil {
		GetLog().Error("Exec Pocket Set Validator Error:", err)
		return "", err
	}
	GetLog().Debug("Exec Pocket Set Validator Result:", rsp)

	return rsp, nil
}

type LogsHashData struct {
	Logs string `json:"logs"`
	Hash string `json:"hash"`
}

func (cli *DockerCli) PoktCtnExecCustodial(address, amount, relayChainIDs, serviceURI, networkID, fee, isBefore, passwd string) (string, error) {
	rsp, err := cli.PoktCtnExec([]string{"expect", POCKET_CONFIG_PATH + "/custodial.sh", address, amount, relayChainIDs, serviceURI, networkID, fee, isBefore, passwd})
	if err != nil {
		GetLog().Error("Exec Pocket Custodial Error:", err)
		return "", err
	}
	GetLog().Debug("Exec Pocket Custodial Result:", rsp)

	resultStr := rsp
	re := regexp.MustCompile(`\"logs\": (.*?),(?:\r\n\s+)?\"txhash\": \"(\S+)\"`)
	match := re.FindStringSubmatch(resultStr)
	fmt.Println(match)
	if len(match) == 3 {
		logsHash := LogsHashData{Logs: match[1], Hash: match[2]}
		result, err := json.MarshalIndent(logsHash, "", "  ")
		if err == nil {
			resultStr = string(result)
		}
	}

	return resultStr, nil
}

func (cli *DockerCli) PoktCtnExecNonCustodial(pubKey, outputAddr, amount, relayChainIDs, serviceURI, networkID, fee, isBefore, passwd string) (string, error) {
	rsp, err := cli.PoktCtnExec([]string{"expect", POCKET_CONFIG_PATH + "/noncustodial.sh", pubKey, outputAddr, amount, relayChainIDs, serviceURI, networkID, fee, isBefore, passwd})
	if err != nil {
		GetLog().Error("Exec Pocket NonCustodial Error:", err)
		return "", err
	}

	resultStr := rsp
	re := regexp.MustCompile(`\"logs\": (.*?),(?:\r\n\s+)?\"txhash\": \"(\S+)\"`)
	match := re.FindStringSubmatch(resultStr)
	if len(match) == 3 {
		logsHash := LogsHashData{Logs: match[1], Hash: match[2]}
		result, err := json.MarshalIndent(logsHash, "", "  ")
		if err == nil {
			resultStr = string(result)
		}
	}

	return resultStr, nil
}
