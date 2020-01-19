package wallets

import "fmt"

type Response struct {
	Code    int
	Data    interface{}
	Message string
}

func (r *Response) GetWalletCredentialID() (string, error) {
	if r.Code != 0 || r.Message != "OK" {
		return "", fmt.Errorf("request create wallet failed\n")
	}
	data, ok := r.Data.(map[string]string)
	if !ok {
		return "", fmt.Errorf("convert interface{} to map[string]string failed\n")
	}
	return data["credential_id"], nil
}

func (r *Response) GetTxID() (string, error) {
	if r.Code != 0 || r.Message != "OK" {
		return "", fmt.Errorf("request send coin failed\n")
	}
	data, ok := r.Data.(map[string]string)
	if !ok {
		return "", fmt.Errorf("convert interface{} to map[string]string failed\n")
	}
	return data["txid"], nil
}

func (r *Response) GetBalanceAndAddr(tokenID string) (string, int) {
	if r.Code != 0 || r.Message != "OK" {
		return "", -1
	}
	data, ok := r.Data.(map[string]interface{})
	if !ok {
		return "", -1
	}
	addr := data["slp_address"].(string)
	balances, ok := data["tokens"].([]Balance)
	if !ok {
		return "", -1
	}
	for _, v := range balances {
		if v.TokenID == tokenID {
			return addr, v.TokenBalance
		}
	}
	return "", -1
}

type Balance struct {
	SatoshiBalance int    `json:"satoshis_balance"`
	TokenID        string `json:"tokenId"`
	TokenBalance   int    `json:"token_balance"`
}
