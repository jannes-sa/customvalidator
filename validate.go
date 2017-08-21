/*
Version : 1.0
Author  : Jannes Santoso
Noted   : Use it only for validation request external data
*/

package customvalidator

import (
	"reflect"
	"strconv"
	"strings"

	"log"
)

// TypeStructAfterScan ...
type TypeStructAfterScan struct {
	Type        string
	NameField   string
	AssignField string
	Code        string
	Value       interface{}
	Validate    string
}

// Validate Custom Validating
func Validate(st interface{}, overflowStruct interface{}) []string {
	var codeError []string

	v := reflect.ValueOf(st)
	vt := v.Type()
	ve := reflect.ValueOf(overflowStruct).Elem()

	var scanDataRequest []TypeStructAfterScan

	for i, n := 0, v.NumField(); i < n; i++ {
		f := v.Field(i)
		ft := vt.Field(i)

		var realVal interface{}
		var realType string
		stateType := true
		getTypeAndVal(f, ft, &stateType, &realVal, &realType)
		runningValidate(f, ft, stateType, realVal, realType, &scanDataRequest, &codeError)

		if len(codeError) == 0 {
			if realType == "string" {
				ve.Field(i).SetString(realVal.(string))
			} else if realType == "int" {
				ve.Field(i).SetInt(int64(realVal.(int)))
			} else if realType == "float64" {
				ve.Field(i).SetFloat(realVal.(float64))
			}
		}
	}

	// Validate After Scan //
	if len(codeError) == 0 {
		validateAfterScan(scanDataRequest, &codeError)
	}
	// Validate After Scan //

	return codeError
}

func runningValidate(f reflect.Value, ft reflect.StructField, stateType bool, realVal interface{},
	realType string, scanDataRequest *[]TypeStructAfterScan, extractCodeError *[]string) {
	validateStr := ft.Tag.Get("validate")
	validateArr := strings.Split(validateStr, ",")

	for _, val := range validateArr {
		valArr := strings.Split(val, "=")
		if valArr[0] == "type" && !stateType {
			*extractCodeError = append(*extractCodeError, valArr[1])
		} else if stateType {
			if valArr[0] == "required" {
				if len(valArr) == 2 {
					requiredValidate(realType, realVal, extractCodeError, valArr[1])
				}
			} else if valArr[0] == "stringnumericonly" {
				if len(valArr) == 2 {
					stringnumericonlyValidate(realType, realVal, extractCodeError, valArr[1])
				}
			} else if valArr[0] == "gte" || valArr[0] == "lte" || valArr[0] == "len" {
				if len(valArr) == 3 {
					gteLteLenValidate(realType, realVal, extractCodeError, valArr)
				}
			} else if valArr[0] == "email" {
				if len(valArr) == 2 {
					emailValidate(realType, realVal, extractCodeError, valArr[1])
				}
			} else if valArr[0] == "identicField" {
				*scanDataRequest = append(*scanDataRequest, TypeStructAfterScan{
					Type:        realType,
					NameField:   ft.Name,
					AssignField: valArr[1],
					Code:        valArr[2],
					Value:       realVal,
					Validate:    valArr[0],
				})
			}
		}
	}
}

// Create Validation Here //
func requiredValidate(realType string, realVal interface{}, extractCodeError *[]string,
	code string) {
	if realType == "string" {
		if realVal.(string) == "" {
			*extractCodeError = append(*extractCodeError, code)
		}
	} else if realType == "int" {
		if realVal.(int) == 0 {
			*extractCodeError = append(*extractCodeError, code)
		}
	} else if realType == "float64" {
		if realVal.(float64) == 0 {
			*extractCodeError = append(*extractCodeError, code)
		}
	}
}

func stringnumericonlyValidate(realType string, realVal interface{}, extractCodeError *[]string,
	code string) {
	if realType == "string" && realVal.(string) != "" {
		_, errConv := strconv.Atoi(realVal.(string))
		if errConv != nil {
			*extractCodeError = append(*extractCodeError, code)
		}
	}
}

func gteLteLenValidate(realType string, realVal interface{}, extractCodeError *[]string,
	valArr []string) {
	intNil, errAtoi := strconv.Atoi(valArr[1])
	CheckErr("Failed Convert custom validate line 69", errAtoi)

	if realType == "string" {
		stCheck := false
		if valArr[0] == "gte" {
			if len(realVal.(string)) >= intNil {
				stCheckAsgn(&stCheck)
			}
		} else if valArr[0] == "lte" {
			if len(realVal.(string)) <= intNil {
				stCheckAsgn(&stCheck)
			}
		} else if valArr[0] == "len" {
			if len(realVal.(string)) == intNil {
				stCheckAsgn(&stCheck)
			}
		}
		if !stCheck {
			*extractCodeError = append(*extractCodeError, valArr[2])
		}
	} else if realType == "int" {
		stCheck := false
		if valArr[0] == "gte" {
			if realVal.(int) >= intNil {
				stCheckAsgn(&stCheck)
			}
		} else if valArr[0] == "lte" {
			if realVal.(int) <= intNil {
				stCheckAsgn(&stCheck)
			}
		} else if valArr[0] == "len" {
			if realVal.(int) == intNil {
				stCheckAsgn(&stCheck)
			}
		}
		if !stCheck {
			*extractCodeError = append(*extractCodeError, valArr[2])
		}
	} else if realType == "float64" {
		stCheck := false
		if valArr[0] == "gte" {
			if realVal.(float64) >= float64(intNil) {
				stCheckAsgn(&stCheck)
			}
		} else if valArr[0] == "lte" {
			if realVal.(float64) <= float64(intNil) {
				stCheckAsgn(&stCheck)
			}
		} else if valArr[0] == "len" {
			if realVal.(float64) == float64(intNil) {
				stCheckAsgn(&stCheck)
			}
		}
		if !stCheck {
			*extractCodeError = append(*extractCodeError, valArr[2])
		}
	}
}
func stCheckAsgn(check *bool) {
	*check = true
}

func emailValidate(realType string, realVal interface{}, extractCodeError *[]string,
	code string) {
	if realType == "string" && realVal.(string) != "" {
		errMail := ValidateFormatMail(realVal.(string))
		if errMail != nil {
			*extractCodeError = append(*extractCodeError, code)
		}
	}
}

///////////////////////

// Validate After Scan //
func validateAfterScan(scanDataRequest []TypeStructAfterScan, extractCodeError *[]string) {
	for _, val := range scanDataRequest {
		if val.Validate == "identicField" {
			st := validateIdentical(scanDataRequest, val)
			if st == false {
				*extractCodeError = append(*extractCodeError, val.Code)
			}
		}
	}
}

func validateIdentical(scanDataRequest []TypeStructAfterScan, valAfterScan TypeStructAfterScan) bool {
	state := true
	for _, val := range scanDataRequest {
		if val.NameField == valAfterScan.AssignField && val.Type == valAfterScan.Type {
			if val.Type == "string" {
				if val.Value.(string) == valAfterScan.Value.(string) {
					state = false
				}
			} else if val.Type == "int" {
				if val.Value.(int) == valAfterScan.Value.(int) {
					state = false
				}
			} else if val.Type == "float64" {
				if val.Value.(float64) == valAfterScan.Value.(float64) {
					state = false
				}
			}
		}
	}
	return state
}

////////////////////////

func getTypeAndVal(f reflect.Value, ft reflect.StructField, stateType *bool, realVal *interface{},
	realType *string) {
	strType := ft.Tag.Get("type")
	arrType := strings.Split(strType, ",")

	checkType(f.Interface(), stateType, arrType, realVal, realType)
}

func checkType(mpt interface{}, state *bool, status []string, val *interface{},
	typeVal *string) {
	*state = false
	switch v := mpt.(type) {
	case int:
		if contains(status, "int") {
			*state = true
		}
		*val = v
		*typeVal = "int"
	case float64:
		if contains(status, "float64") {
			*state = true
		}
		*val = v
		*typeVal = "float64"
	case string:
		if contains(status, "string") {
			*state = true
		}
		*val = v
		*typeVal = "string"
	default:
		*state = false
	}
}

func contains(slice []string, item string) bool {
	set := make(map[string]struct{}, len(slice))
	for _, s := range slice {
		set[s] = struct{}{}
	}

	_, ok := set[item]
	return ok
}

// CheckErr ...
func CheckErr(msg string, err error) {
	if err != nil {
		log.Println(msg)
		panic(err)
	}
}
