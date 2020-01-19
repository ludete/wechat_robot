package wallets

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
)

type IfWallet struct {
	addr         string
	apiKey       string
	secretKey    string
	credentialID string
}

type TransferNews struct {
	Address string
	Amount  float64
	Denom   string `json:"denom,omitempty"`
}

func (i IfWallet) SendMoney(from string, news []TransferNews) (string, error) {
	client := &http.Client{}
	url := i.addr + "/h5_open/hd_wallet/slp/tx"
	data := make(map[string]interface{})
	data["credential_id"] = i.credentialID
	data["token_id_hex"] = "4de69e374a8ed21cbddd47f2338cc0f479dc58daa2bbe11cd604ca488eca0ddf"
	data["receivers"] = news
	bz, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bz))
	if err != nil {
		log.Errorf("new request failed, error : %s\n", err.Error())
		return "", err
	}
	i.fillReqHeader(req.Header)
	resp, err := client.Do(req)
	defer resp.Body.Close()

	bz, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("read data from wallet response failed : %s\n", err.Error())
		return "", err
	}
	var res Response
	if err = json.Unmarshal(bz, &res); err != nil {
		log.Errorf("unmarshal create wallet response to json failed : %s\n", err.Error())
		return "", err
	}
	txid, err := res.GetTxID()
	if err != nil {
		log.Errorf("get walletID from create wallet response failed : %s\n", err.Error())
		return "", err
	}
	return txid, nil
}

func (i IfWallet) GetAmountOfDenoms(credentialID string, denom []string) (string, int) {
	client := &http.Client{}
	url := i.addr + "/h5_open/hd_wallet/slp/balance"
	data := make(map[string]interface{})
	data["credential_id"] = credentialID
	data["tokens_id_hex"] = denom
	bz, _ := json.Marshal(data)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bz))
	if err != nil {
		log.Errorf("new request failed, error : %s\n", err.Error())
		return "", -1
	}
	i.fillReqHeader(req.Header)
	resp, err := client.Do(req)
	defer resp.Body.Close()

	bz, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("read data from wallet response failed : %s\n", err.Error())
		return "", -1
	}
	var res Response
	if err = json.Unmarshal(bz, &res); err != nil {
		log.Errorf("unmarshal create wallet response to json failed : %s\n", err.Error())
		return "", -1
	}
	addr, balance := res.GetBalanceAndAddr(denom[0])
	if len(addr) == 0 && balance < 0 {
		log.Errorf("get walletID from create wallet response failed : %s\n", err.Error())
		return "", balance
	}
	return addr, balance
}

func (i IfWallet) GetAllAmounts(addr string) map[string]int {
	panic("implement me")
}

func (i IfWallet) CreateUserDenomAddr() string {
	client := &http.Client{}
	url := i.addr + "/h5_open/hd_wallet"
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		log.Errorf("new request failed, error : %s\n", err.Error())
		return ""
	}
	i.fillReqHeader(req.Header)
	resp, err := client.Do(req)
	defer resp.Body.Close()

	bz, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("read data from wallet response failed : %s\n", err.Error())
		return ""
	}
	var res Response
	if err = json.Unmarshal(bz, &res); err != nil {
		log.Errorf("unmarshal create wallet response to json failed : %s\n", err.Error())
		return ""
	}
	id, err := res.GetWalletCredentialID()
	if err != nil {
		log.Errorf("get walletID from create wallet response failed : %s\n", err.Error())
		return ""
	}
	return id
}

func (i IfWallet) fillReqHeader(header http.Header) {
	header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	now := strconv.Itoa(int(time.Now().Unix()))
	header.Set("TIME", now)
	header.Set("APIKEY", i.apiKey)
	header.Set("SIGN", generateSign(i.apiKey, i.secretKey, now))
}

func generateSign(apiKey, secretKey, time string) string {
	bz := make([]byte, 0, len(apiKey)+len(secretKey)+len(time))
	bz = append(bz, []byte(apiKey)...)
	bz = append(bz, []byte(secretKey)...)
	bz = append(bz, []byte(time)...)
	sum := md5.Sum(bz)
	return string(sum[:])
}
