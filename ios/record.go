package ios

import (
	"time"

	"github.com/tidwall/gjson"
	"github.com/u03013112/ss-ios-purchase/sql"
)

// Bills : ios 账单记录
type Bills struct {
	sql.BaseModel
	UUID          string    `json:"uuid,omitempty"`
	ProductID     string    `json:"productId,omitempty"`
	TransactionID string    `json:"transactionId,omitempty"`
	Environment   string    `json:"environment,omitempty"`
	PurchaseDate  time.Time `json:"purchaseDate,omitempty"`
}

func recordBills(str string, uuid string) {
	environment := gjson.Get(str, "environment").String()
	v := gjson.Get(str, "latest_receipt_info")
	if len(v.Array()) > 0 {
		b := v.Array()[0].Raw
		productID := gjson.Get(b, "product_id").String()
		transactionID := gjson.Get(b, "transaction_id").String()
		ex := gjson.Get(b, "purchase_date_ms").Int()
		PurchaseDate := time.Unix(ex/int64(1000), 0)

		var bill Bills
		db := sql.GetInstance().First(&bill, "transaction_id = ?", transactionID)
		if db.RowsAffected == 0 {
			bill.UUID = uuid
			bill.ProductID = productID
			bill.TransactionID = transactionID
			bill.Environment = environment
			bill.PurchaseDate = PurchaseDate
			sql.GetInstance().Create(&bill)
		}
	}
}
