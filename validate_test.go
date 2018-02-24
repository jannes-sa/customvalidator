package customvalidator

import (
	"encoding/json"
	"log"
	"testing"
)

func TestValidate(t *testing.T) {
	type TesInterface struct {
		BeneficiaryAccountNumber interface{} `json:"beneficiary_account_number" type:"string" convert:"removezero=E02000401" validate:"required=E02030009,stringnumericonly=E02030008,lte=20=E02030007,type=E02030001,identicField=="`
		BeneficiaryBankCode      interface{} `json:"beneficiary_bank_code" type:"string" validate:"stringnumericonly=E02030018,lte=3=E02030017,required=E02030019"`
		SourceAccountNumber      interface{} `json:"source_account_number" type:"string" convert:"removezero=E02000401" validate:"stringnumericonly=E02030013,lte=20=E02030014,required=E02030015,identicField=BeneficiaryAccountNumber=E02030016"`
		SourceBankCode           interface{} `json:"source_bank_code" type:"string" validate:"stringnumericonly=E02030021,lte=3=E02030020,required=E02030022"`
		TransactionAmount        interface{} `json:"transaction_amount" type:"float64" validate:",gte=0=E0209898984,lte=999999999999999=E02030010,type=E02030011,required=E02030012"`
	}

	type Tes struct {
		BeneficiaryAccountNumber string  `json:"beneficiary_account_number"`
		BeneficiaryBankCode      string  `json:"beneficiary_bank_code"`
		SourceAccountNumber      string  `json:"source_account_number"`
		SourceBankCode           string  `json:"source_bank_code"`
		TransactionAmount        float64 `json:"transaction_amount"`
	}

	bodyJSON := []byte(`
		{
			"beneficiary_account_number":"123456789123",
			"beneficiary_bank_code":"009",
			"source_account_number":"97837641248321",
			"source_bank_code":"021",
			"transaction_amount":100000
		}
	`)

	var tesInter TesInterface
	err := json.Unmarshal(bodyJSON, &tesInter)
	if err != nil {
		t.Error("Error ", err)
		return
	}
	log.Println(tesInter)

	var tes Tes
	codeError := Validate(tesInter, &tes)

	log.Println(tes, codeError)
}

func TestValidateArray(t *testing.T) {
	type ArrTesInterface struct {
		AccountNumber interface{} `json:"account_number" type:"string" validate:"required=E02030009,stringnumericonly=E02030008,lte=20=E02030007,type=E02030001"`
	}

	type TesInterface2 struct {
		TransactionAmount interface{}       `json:"transaction_amount" type:"float64" validate:"gte=0=E0209898984,lte=999999999999999=E02030010,type=E02030011,required=E02030012"`
		Data              []ArrTesInterface `json:"data"`
	}

	type ArrTes struct {
		AccountNumber string `json:"account_number"`
	}

	type Tes2 struct {
		TransactionAmount float64
		Data              []ArrTes
	}

	bodyJSON := []byte(`
		{
			"transaction_amount":100000,
			"data":[
				{"account_number":"123456789123"},
				{"account_number":"985654637462"}
			]
		}
	`)

	var tesInter TesInterface2
	err := json.Unmarshal(bodyJSON, &tesInter)
	if err != nil {
		t.Error("Error ", err)
		return
	}

	var tes Tes2
	errCode := Validate(tesInter, &tes)

	log.Println(errCode, tes)

}