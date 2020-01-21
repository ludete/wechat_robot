package wallets

import "fmt"

type IFResponse struct {
	Code    int
	Data    interface{}
	Message string
}

func (r *IFResponse) GetWalletCredentialID() (string, error) {
	if r.Code != 0 || r.Message != "OK" {
		fmt.Printf("%s\n", r.Message)
		return "", fmt.Errorf("code : %d, message : %s\n", r.Code, r.Message)
	}
	data, ok := r.Data.(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("convert interface{} to map[string]string failed\n")
	}
	return data["credential_id"].(string), nil
}

func (r *IFResponse) GetTxID() (string, error) {
	if r.Code != 0 || r.Message != "OK" {
		return "", fmt.Errorf("request send coin failed, code : %d, error : %s", r.Code, r.Message)
	}
	data, ok := r.Data.(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("convert interface{} to map[string]string failed\n")
	}
	return data["txid"].(string), nil
}

func (r *IFResponse) GetBalanceAndAddr(tokenID string) (string, int) {
	if r.Code != 0 || r.Message != "OK" {
		return "", -1
	}
	data, ok := r.Data.(map[string]interface{})
	if !ok {
		return "", -1
	}
	addr := data["slp_address"].(string)
	balances, ok := data["tokens"].([]interface{})
	if !ok {
		return "", -1
	}

	for _, v := range balances {
		balance := v.(map[string]interface{})
		if balance["tokenId"].(string) == tokenID {
			return addr, int(balance["token_balance"].(float64))
		}
	}
	return "", -1
}

type Balance struct {
	SatoshiBalance float64 `json:"satoshis_balance"`
	TokenID        string  `json:"tokenId"`
	TokenBalance   float64 `json:"token_balance"`
}
