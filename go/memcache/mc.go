package mc

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"runtime"
	"strings"
)

type mcStore interface {
	GetWithErr(context.Context, string) (string, error)
	SetWithErr(context.Context, string, string, int) error
	Delete(context.Context, string) error
}

type mc struct {
	Store mcStore
}

var m = mc{}

// Setup 初始化 mc
func Setup(ms mcStore) {
	m.Store = ms
}

// GenCacheKey 自动生成 cacheKey
// 规则：functionName|namespace|args.Join(":")
// 如 getName(context.Background(), 1, "a")，namespace 为 v1
// cacheKey = package.getName|v1|1:a
func GenCacheKey(f interface{}, namespace string, args ...interface{}) string {
	var argsStrArr []string

	fullFuncName := runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
	arr := strings.Split(fullFuncName, "/")
	funcName := arr[len(arr)-1]

	// args 去掉 context
	if len(args) != 0 {
		av := reflect.ValueOf(args[0])
		if strings.HasPrefix(fmt.Sprintf("%v", av), "context.") {
			if len(args) == 1 {
				args = nil
			} else {
				args = args[1:len(args)]
			}
		}
	}

	for _, arg := range args {
		argsStrArr = append(argsStrArr, fmt.Sprintf("%v", arg))
	}
	argsStr := strings.Join(argsStrArr, ":")

	return fmt.Sprintf("%s|%s|%v", funcName, namespace, argsStr)
}

// Cache 缓存函数结果，输出 [][]byte，和 error，可指定 cacheKey
func Cache(cacheKey, namespace string, expire int, f interface{}, args ...interface{}) ([][]byte, error) {
	var rets [][]byte
	var err error
	var retValues, argValues []reflect.Value
	var numResultOut int

	ctx := context.Background()
	if cacheKey == "" {
		cacheKey = GenCacheKey(f, namespace, args...)
	}

	result, err1 := m.Store.GetWithErr(ctx, cacheKey)
	if result != "" && err1 == nil {
		rets = decodeRetJSONs(result)
		return rets, nil
	}

	for _, arg := range args {
		argValues = append(argValues, reflect.ValueOf(arg))
	}

	fValue := reflect.ValueOf(f)
	numOut := fValue.Type().NumOut()

	hasErrorOut := func() bool {
		return fValue.Type().Out(numOut-1).Name() == "error"
	}

	if fValue.Kind() == reflect.Func {
		retValues = fValue.Call(argValues)
	}

	if hasErrorOut() {
		numResultOut = numOut - 1
		errIn := retValues[numOut-1].Interface()
		if errIn != nil {
			err = errIn.(error)
		}
	} else {
		numResultOut = numOut
	}

	for i := 0; i < numResultOut; i++ {
		retJSON, err2 := json.Marshal(retValues[i].Interface())
		if err2 != nil {
			return rets, err2
		}
		rets = append(rets, retJSON)
	}

	data := encodeRetJSONs(rets)

	// set cache
	_ = m.Store.SetWithErr(ctx, cacheKey, data, expire)
	return rets, err
}

// Unmarshal 解析 Cache 结果，本质上就是 json.Unmarshal，所以遵守 json.Unmarshal 规则:
// bool, for JSON booleans
// float64, for JSON numbers
// string, for JSON strings
// []interface{}, for JSON arrays
// map[string]interface{}, for JSON objects
// nil for JSON null
func Unmarshal(src []byte, dst interface{}) error {
	return json.Unmarshal(src, dst)
}

func Invalidate(ctx context.Context, cacheKey string) error {
	return m.Store.Delete(ctx, cacheKey)
}

func encodeRetJSONs(retJSONs [][]byte) string {
	var arr []string
	for _, retJSON := range retJSONs {
		arr = append(arr, string(retJSON))
	}
	if len(arr) == 0 {
		return ""
	}
	return strings.Join(arr, "::")
}

func decodeRetJSONs(retStr string) [][]byte {
	var arr []string
	var retJSONs [][]byte
	arr = strings.Split(retStr, "::")

	for _, item := range arr {
		retJSONs = append(retJSONs, []byte(item))
	}
	return retJSONs
}
