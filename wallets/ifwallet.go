package wallets

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
)

type IfWallet struct {
	addr      string
	apiKey    string
	secretKey string
}

type TransferNews struct {
	Address string `json:"address"`
	Amount  int64  `json:"amount"`
	Denom   string `json:"denom,omitempty"`
}

func NewWallet(url, apiKey, secretKey string) WalletInterface {
	return &IfWallet{
		addr:      url,
		apiKey:    apiKey,
		secretKey: secretKey,
	}
}

func (i *IfWallet) CreateUserWallet() (string, error) {
	url := i.addr + "/h5_open/hd_wallet"
	res, err := i.sendRequest(url, nil)
	if err != nil {
		return "", err
	}
	id, err := res.GetWalletCredentialID()
	if err != nil {
		log.Errorf("get walletID from create wallet response failed : %s\n", err.Error())
		return "", err
	}
	return id, nil
}

func (i IfWallet) SendMoney(walletID string, news []TransferNews) (string, error) {
	url := i.addr + "/h5_open/hd_wallet/slp/tx"
	denomID := getTokenID(news[0].Denom)
	if denomID == "" {
		return "", fmt.Errorf("未知的币种 ：%s\n", news[0].Denom)
	}
	data := make(map[string]interface{})
	data["credential_id"] = walletID
	data["token_id_hex"] = denomID
	data["receivers"] = news
	bz, _ := json.Marshal(data)
	fmt.Printf("send money : %s\n", bz)
	res, err := i.sendRequest(url, bytes.NewBuffer(bz))
	if err != nil {
		return "", err
	}
	txid, err := res.GetTxID()
	if err != nil {
		log.Errorf("get walletID from create wallet response failed : %s\n", err.Error())
		return "", err
	}
	return txid, nil
}

func (i IfWallet) GetAmountOfDenoms(credentialID string, denom string) (string, int, error) {
	url := i.addr + "/h5_open/hd_wallet/slp/balance"
	denomID := getTokenID(denom)
	if denomID == "" {
		return "", -1, fmt.Errorf("未知的币种 ：%s\n", denom)
	}
	data := make(map[string]interface{})
	data["credential_id"] = credentialID
	data["tokens_id_hex"] = []string{denomID}
	bz, _ := json.Marshal(data)
	fmt.Println("send money bz : ", string(bz))
	fmt.Println("send money bz : ", len(string(bz)))
	//fmt.Println("unicode size : ", len(toUnicode(string(bz))))
	res, err := i.sendRequest(url, bytes.NewBuffer(bz))
	if err != nil {
		return "", -1, err
	}
	addr, balance := res.GetBalanceAndAddr(getTokenID(denom))
	if len(addr) == 0 {
		log.Errorf("get balance from wallet response failed, code : %d, message : %s", res.Code, res.Message)
		return "", balance, fmt.Errorf("get balance from wallet response failed, code : %d, message : %s", res.Code, res.Message)
	}
	return addr, balance, nil
}

func (i *IfWallet) fillReqHeader(header http.Header) {
	header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	now := strconv.Itoa(int(time.Now().Unix()))
	header.Set("TIME", now)
	header.Set("APIKEY", i.apiKey)
	header.Set("SIGN", generateSign(i.apiKey, i.secretKey, now))
	header.Set("Content-Type", "application/json")
	header.Set("Accept", "*/*")
	header.Set("Accept-Encoding", "gzip")
	header.Set("Accept-Encoding", "deflate")
	header.Set("Accept-Encoding", "br")
	header.Set("Cache-Control", "no-cache")
	header.Set("Connection", "keep-alive")

}

func generateSign(apiKey, secretKey, time string) string {
	bz := make([]byte, 0, len(apiKey)+len(secretKey)+len(time))
	bz = append(bz, []byte(apiKey)...)
	bz = append(bz, []byte(secretKey)...)
	bz = append(bz, []byte(time)...)
	sum := sha256.Sum256(bz)
	dst := make([]byte, len(sum)*2)
	hex.Encode(dst, sum[:])
	return string(dst)
}

func (i *IfWallet) sendRequest(route string, body io.Reader) (*IFResponse, error) {
	client := &http.Client{}
	req, err := http.NewRequest("POST", route, body)
	if err != nil {
		log.Errorf("new request failed, error : %s\n", err.Error())
		return nil, err
	}
	i.fillReqHeader(req.Header)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	fmt.Println(req.Header)
	fmt.Println(req.Body)

	bz, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("read data from wallet response failed : %s\n", err.Error())
		return nil, fmt.Errorf("read data from wallet response failed : %s\n", err.Error())
	}
	var res IFResponse
	if err = json.Unmarshal(bz, &res); err != nil {
		log.Errorf("unmarshal create wallet response to json failed : %s\n", err.Error())
		return nil, fmt.Errorf("unmarshal create wallet response to json failed : %s\n", err.Error())
	}
	return &res, nil
}

func getTokenID(denom string) string {
	if denom == SPICE {
		return SPICEID
	}
	return ""
}

func toUnicode(str string) string {
	runes := []rune(str)
	res := ""
	for _, r := range runes {
		if r < rune(128) {
			res += string(r)
		} else {
			res += "\\u" + strconv.FormatInt(int64(r), 16)
		}
	}
	return res
}
