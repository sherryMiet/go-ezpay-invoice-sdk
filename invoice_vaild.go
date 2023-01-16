package ezpay_invoice

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"strconv"
	"time"
)

const (
	TestInvoiceInvalidURL     = "https://cinv.ezpay.com.tw/Api/invoice_invalid"
	InvoiceInvalidURL         = "https://inv.ezpay.com.tw/Api/invoice_invalid"
	InvoiceInvalidVersion     = "1.0"
	InvoiceInvalidRespondType = "JSON"
)

type InvoiceInvalidCall struct {
	HashKey               string
	HashIV                string
	InvoiceInvalidRequest *InvoiceInvalidRequest
}

type InvoiceInvalidRequest struct {
	MerchantID_ string
	PostData_   string
}

type InvoiceInvalidRequestPostData_ struct {
	//回傳格式
	RespondType string
	//串接程式版本
	Version       string
	TimeStamp     string
	InvoiceNumber string
	InvalidReason string
}

type InvoiceInvalidResponse struct {
	Status  string
	Message string
	Result  string
}

type InvoiceInvalidResponseResult struct {
	//商店代號
	MerchantID string
	//發票號碼
	InvoiceNumber string
	//開立發票時間
	CreateTime string
	//檢查碼
	CheckCode string
}

func (i *InvoiceInvalidRequestPostData_) SetInvoiceData(InvoiceNumber, InvalidReason string) *InvoiceInvalidRequestPostData_ {
	i.InvoiceNumber = InvoiceNumber
	i.InvalidReason = InvalidReason
	return i
}

func (c *Client) InvoiceInvalid(postData *InvoiceInvalidRequestPostData_) *InvoiceInvalidCall {
	postData.TimeStamp = strconv.Itoa(int(time.Now().Unix()))
	postData.Version = InvoiceInvalidVersion
	postData.RespondType = InvoiceInvalidRespondType
	paramsStr := StructToParamsMap(postData)
	postDataStr := ParamsMapToURLEncode(paramsStr)
	encrypt, err := AesCBCEncrypt([]byte(postDataStr), []byte(c.HashKey), []byte(c.HashIV))
	if err != nil {
		logrus.Error(err)
		return nil
	}
	return &InvoiceInvalidCall{
		HashIV:  c.HashIV,
		HashKey: c.HashKey,
		InvoiceInvalidRequest: &InvoiceInvalidRequest{
			MerchantID_: c.MerchantID,
			PostData_:   encrypt,
		},
	}
}

func (i *InvoiceInvalidCall) Do() *InvoiceInvalidResponse {
	PostData := make(map[string]string)
	PostData["MerchantID_"] = i.InvoiceInvalidRequest.MerchantID_
	PostData["PostData_"] = i.InvoiceInvalidRequest.PostData_
	body, err := SendEZPayRequest(&PostData, InvoiceInvalidURL)
	if err != nil {
		logrus.Error(err)
		return nil
	}
	Res := new(InvoiceInvalidResponse)
	err = json.Unmarshal(body, Res)
	if err != nil {
		logrus.Error(err)
		return nil
	}
	return Res
}
func (i *InvoiceInvalidCall) DoTest() *InvoiceInvalidResponse {
	PostData := make(map[string]string)
	PostData["MerchantID_"] = i.InvoiceInvalidRequest.MerchantID_
	PostData["PostData_"] = i.InvoiceInvalidRequest.PostData_
	body, err := SendEZPayRequest(&PostData, TestInvoiceInvalidURL)
	if err != nil {
		logrus.Error(err)
		return nil
	}
	print(string(body))
	Res := new(InvoiceInvalidResponse)
	err = json.Unmarshal(body, Res)
	if err != nil {
		logrus.Error(err)
		return nil
	}
	return Res
}
