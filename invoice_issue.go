package ezpay_invoice

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"log"
	"strconv"
	"time"
)

const (
	TestInvoiceIssueURL     = "https://cinv.ezpay.com.tw/Api/invoice_issue"
	InvoiceIssueURL         = "https://inv.ezpay.com.tw/Api/invoice_issue"
	InvoiceIssueVersion     = "1.5"
	InvoiceIssueRespondType = "JSON"
)

type InvoiceIssueCall struct {
	HashKey             string
	HashIV              string
	InvoiceIssueRequest *InvoiceIssueRequest
}

type InvoiceIssueRequest struct {
	MerchantID_ string
	PostData_   string
}

type InvoiceIssueRequestPostData_ struct {
	//回傳格式
	RespondType string
	//串接程式版本
	Version string
	//時間戳記
	TimeStamp string
	//ezPay 平台 交易序號
	TransNum string
	//自訂編號
	MerchantOrderNo string
	//開立發票方式
	Status string
	//預計開立日期
	CreateStatusTime string
	//發票種類
	Category string
	//買受人名稱
	BuyerName string
	//買受人統一編號
	BuyerUBN string
	//買受人地址
	BuyerAddress string
	//買受人電子信箱
	BuyerEmail string
	//載具類別
	CarrierType string
	//載具編號
	CarrierNum string
	//捐贈碼
	LoveCode string
	//索取紙本發票
	PrintFlag string
	//是否開放至合作超商 Kiosk 列印
	KioskPrintFlag string
	//課稅別
	TaxType string
	//稅率
	TaxRate float32
	//報關標記
	CustomsClearance string
	//銷售額合計
	Amt int
	//銷售額 (課稅別應稅)
	AmtSales int
	//銷售額 (課稅別零稅率)
	AmtZero int
	//銷售額 (課稅別免稅)
	AmtFree int
	//稅額
	TaxAmt int
	//發票金額
	TotalAmt int
	//商品名稱
	ItemName string
	//商品數量
	ItemCount string
	//商品單位
	ItemUnit string
	//商品單價
	ItemPrice string
	//商品小計
	ItemAmt string
	//商品課稅別
	ItemTaxType string
	//備註
	Comment string
}

type InvoiceIssueResponse struct {
	Status  string
	Message string
	Result  json.RawMessage
}

type InvoiceIssueResponseResult struct {
	//商店代號
	MerchantID string
	//ezPay 電子發票 開立序號
	InvoiceTransNo string
	//自訂編號
	MerchantOrderNo string
	//發票金額
	TotalAmt int
	//發票號碼
	InvoiceNumber string
	//發票防偽隨機碼
	RandomNum string
	//開立發票時間
	CreateTime string
	//檢查碼
	CheckCode string
	//發票條碼
	BarCode string
	//發票 QRCode(左)
	QRcodeL string
	//發票 QRCode(右)
	QRcodeR string
}

func (i *InvoiceIssueRequestPostData_) SetInvoiceData(MerchantOrderNo, TaxType, CustomsClearance, Comment string, Amt, TaxAmt, TotalAmt int, TaxRate float32) *InvoiceIssueRequestPostData_ {
	i.MerchantOrderNo = MerchantOrderNo
	i.TaxType = TaxType
	i.TaxRate = TaxRate
	i.CustomsClearance = CustomsClearance
	i.Comment = Comment
	i.Amt = Amt
	i.TotalAmt = TotalAmt
	i.TaxAmt = TaxAmt
	return i
}
func (i *InvoiceIssueRequestPostData_) SetAmtSales(AmtSales int) *InvoiceIssueRequestPostData_ {
	i.AmtSales = AmtSales
	return i
}
func (i *InvoiceIssueRequestPostData_) SetAmtZero(AmtZero int) *InvoiceIssueRequestPostData_ {
	i.AmtZero = AmtZero
	return i
}
func (i *InvoiceIssueRequestPostData_) SetAmtFree(AmtFree int) *InvoiceIssueRequestPostData_ {
	i.AmtFree = AmtFree
	return i
}

func (i *InvoiceIssueRequestPostData_) IssueNow() *InvoiceIssueRequestPostData_ {
	i.Status = "1"
	return i
}
func (i *InvoiceIssueRequestPostData_) IssueWait() *InvoiceIssueRequestPostData_ {
	i.Status = "2"
	return i
}
func (i *InvoiceIssueRequestPostData_) IssueAppointment(CreateStatusTime string) *InvoiceIssueRequestPostData_ {
	i.Status = "3"
	i.CreateStatusTime = CreateStatusTime
	return i
}

func (i *InvoiceIssueRequestPostData_) B2B(BuyerName, BuyerUBN, BuyerAddress, BuyerEmail string) *InvoiceIssueRequestPostData_ {
	i.Category = "B2B"
	i.BuyerName = BuyerName
	i.BuyerUBN = BuyerUBN
	i.BuyerAddress = BuyerAddress
	i.BuyerEmail = BuyerEmail
	i.PrintFlag = "Y"
	return i
}

// ezpay會員載具
func (i *InvoiceIssueRequestPostData_) B2CEZPAYMember(BuyerName, BuyerAddress, BuyerEmail, CarrierNum string) *InvoiceIssueRequestPostData_ {
	i.Category = "B2C"
	i.BuyerName = BuyerName
	i.CarrierType = "2"
	i.BuyerAddress = BuyerAddress
	i.BuyerEmail = BuyerEmail
	i.CarrierNum = CarrierNum
	i.PrintFlag = "N"
	return i
}

// 手機載具
func (i *InvoiceIssueRequestPostData_) B2CPhoneCarrier(BuyerName, BuyerAddress, BuyerEmail, CarrierNum string) *InvoiceIssueRequestPostData_ {
	i.Category = "B2C"
	i.BuyerName = BuyerName
	i.CarrierType = "0"
	i.BuyerAddress = BuyerAddress
	i.BuyerEmail = BuyerEmail
	i.CarrierNum = CarrierNum
	i.PrintFlag = "N"
	return i
}

// 自然人憑證
func (i *InvoiceIssueRequestPostData_) B2CCertificate(BuyerName, BuyerAddress, BuyerEmail, CarrierNum string) *InvoiceIssueRequestPostData_ {
	i.Category = "B2C"
	i.BuyerName = BuyerName
	i.CarrierType = "1"
	i.BuyerAddress = BuyerAddress
	i.BuyerEmail = BuyerEmail
	i.CarrierNum = CarrierNum
	i.PrintFlag = "N"
	return i
}

// 捐贈
func (i *InvoiceIssueRequestPostData_) B2CDonation(BuyerName, BuyerAddress, BuyerEmail, LoveCode string) *InvoiceIssueRequestPostData_ {
	i.Category = "B2C"
	i.BuyerName = BuyerName
	i.BuyerAddress = BuyerAddress
	i.BuyerEmail = BuyerEmail
	i.LoveCode = LoveCode
	i.PrintFlag = "N"
	return i
}

// 皆為空值
func (i *InvoiceIssueRequestPostData_) B2CNothing(BuyerName, BuyerAddress, BuyerEmail string) *InvoiceIssueRequestPostData_ {
	i.Category = "B2C"
	i.BuyerName = BuyerName
	i.BuyerAddress = BuyerAddress
	i.BuyerEmail = BuyerEmail
	i.PrintFlag = "Y"
	return i
}

func (i *InvoiceIssueRequestPostData_) SetItem(ItemName, ItemCount, ItemUnit, ItemPrice, ItemAmt, ItemTaxType string) {
	i.ItemName = ItemName
	i.ItemCount = ItemCount
	i.ItemUnit = ItemUnit
	i.ItemPrice = ItemPrice
	i.ItemAmt = ItemAmt
	i.ItemTaxType = ItemTaxType
}

func (c *Client) InvoiceIssue(postData *InvoiceIssueRequestPostData_) *InvoiceIssueCall {
	postData.TimeStamp = strconv.Itoa(int(time.Now().Unix()))
	postData.Version = InvoiceIssueVersion
	postData.RespondType = InvoiceIssueRespondType
	paramsStr := StructToParamsMap(postData)
	log.Print(paramsStr)
	postDataStr := ParamsMapToURLEncode(paramsStr)
	log.Print(postDataStr)
	encrypt, err := AesCBCEncrypt([]byte(postDataStr), []byte(c.HashKey), []byte(c.HashIV))
	if err != nil {
		logrus.Error(err)
		return nil
	}
	log.Print(encrypt)
	return &InvoiceIssueCall{
		HashIV:  c.HashIV,
		HashKey: c.HashKey,
		InvoiceIssueRequest: &InvoiceIssueRequest{
			MerchantID_: c.MerchantID,
			PostData_:   encrypt,
		},
	}
}

func (i *InvoiceIssueCall) Do() *InvoiceIssueResponse {
	PostData := make(map[string]string)
	PostData["MerchantID_"] = i.InvoiceIssueRequest.MerchantID_
	PostData["PostData_"] = i.InvoiceIssueRequest.PostData_
	body, err := SendEZPayRequest(&PostData, InvoiceIssueURL)
	if err != nil {
		logrus.Error(err)
		return nil
	}
	Res := new(InvoiceIssueResponse)
	err = json.Unmarshal(body, Res)
	if err != nil {
		logrus.Error(err)
		return nil
	}
	return Res
}
func (i *InvoiceIssueCall) DoTest() *InvoiceIssueResponse {
	PostData := make(map[string]string)
	PostData["MerchantID_"] = i.InvoiceIssueRequest.MerchantID_
	PostData["PostData_"] = i.InvoiceIssueRequest.PostData_
	body, err := SendEZPayRequest(&PostData, TestInvoiceIssueURL)
	if err != nil {
		logrus.Error(err)
		return nil
	}
	print(string(body))
	Res := new(InvoiceIssueResponse)
	err = json.Unmarshal(body, Res)
	if err != nil {
		logrus.Error(err)
		return nil
	}
	return Res
}
