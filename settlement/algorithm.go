package settlement

type PurchaseAlgo interface {
	Purchase(marketPrice int, usdtRMBPrice int, purchaseRmbAmount int)int
}