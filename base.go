package ezpay_invoice

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

func ParamsMapToURLEncode(data map[string]string) string {
	var uri url.URL
	q := uri.Query()
	for key, value := range data {
		q.Set(key, value)
	}
	queryStr := q.Encode()
	return queryStr
}
func StructToParamsMap(data interface{}) map[string]string {
	params := map[string]string{}
	iVal := reflect.ValueOf(data)
	iTyp := reflect.TypeOf(data)
	if iVal.Kind() == reflect.Ptr {
		iVal = iVal.Elem()
	}
	if iTyp.Kind() == reflect.Ptr {
		iTyp = iTyp.Elem()
	}

	//stringerInterface := reflect.TypeOf((*fmt.Stringer)(nil)).Elem()
	marshalerInterface := reflect.TypeOf((*json.Marshaler)(nil)).Elem()
	//typ := iVal.Type()
	for i := 0; i < iVal.NumField(); i++ {

		f := iVal.Field(i)
		ft := iTyp.Field(i)

		if f.Kind() == reflect.Ptr && f.IsNil() {
			continue
		}
		if f.Kind() == reflect.Ptr {
			f = f.Elem()

		}

		if ft.Anonymous {
			nestedParams := StructToParamsMap(f.Interface())
			for key, val := range nestedParams {
				params[key] = val
			}
			continue
		}

		var v string
		//f.Anonymous
		switch realVal := f.Interface().(type) {
		case int, int8, int16, int32, int64:
			tag := ft.Tag.Get("json")
			if strings.Contains(tag, ",omitempty") && f.Int() == 0 {
				continue
			}
			v = strconv.FormatInt(f.Int(), 10)
		case uint, uint8, uint16, uint32, uint64:
			tag := ft.Tag.Get("json")
			if strings.Contains(tag, ",omitempty") && f.Uint() == 0 {
				continue
			}
			v = strconv.FormatUint(f.Uint(), 10)
		case float32:
			tag := ft.Tag.Get("json")
			if strings.Contains(tag, ",omitempty") && f.Float() == 0 {
				continue
			}
			v = strconv.FormatFloat(f.Float(), 'f', 4, 32)
		case float64:
			tag := ft.Tag.Get("json")
			if strings.Contains(tag, ",omitempty") && f.Float() == 0 {
				continue
			}
			v = strconv.FormatFloat(f.Float(), 'f', 4, 64)
		case []byte:
			tag := ft.Tag.Get("json")
			if strings.Contains(tag, ",omitempty") && string(realVal) == "" {
				continue
			}
			v = string(realVal)
		case string:
			tag := ft.Tag.Get("json")
			if strings.Contains(tag, ",omitempty") && string(realVal) == "" {
				continue
			}
			v = realVal
		//case time.Time:
		//	v = realVal.Format(time.RFC3339)
		//case base.ECPayDateTime:
		//	v = realVal.String()
		default:
			switch {
			case f.Type().Implements(marshalerInterface):
				data, err := f.Interface().(json.Marshaler).MarshalJSON()
				if err != nil {
					panic(err)
				}
				unquoteData, err := strconv.Unquote(string(data))
				if err != nil {
					v = string(data)
				}
				v = unquoteData
			//case f.Type().Implements(stringerInterface):
			//	v = f.Interface().(fmt.Stringer).String()
			default:
				switch f.Kind() {
				case reflect.String:
					v = f.String()
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					v = strconv.FormatInt(f.Int(), 10)
				default:
					panic(fmt.Sprintf("Unknown type %T during transfer struct to map!\n", f.Interface()))

				}
			}
		}
		params[ft.Name] = v
	}
	return params
}

func SendEZPayRequest(postData *map[string]string, URL string) ([]byte, error) {
	body := new(bytes.Buffer)
	w := multipart.NewWriter(body)
	for k, v := range *postData {
		w.WriteField(k, v)
	}
	w.Close()
	log.Print(body.String())
	req, _ := http.NewRequest(http.MethodPost, URL, body)
	req.Header.Set("Content-Type", w.FormDataContentType())
	resp, _ := http.DefaultClient.Do(req)
	data, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	log.Print(resp.StatusCode)
	log.Printf("%s", data)
	return data, nil
}
