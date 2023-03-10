package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"net"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

const (
	PI           float64 = 3.14159265
	EARTH_RADIUS int     = 6378137
	RAD          float64 = PI / 180.0
	InvalidFloat         = math.MaxFloat64
)

func CalEarthDistance(lat1 float64, lng1 float64, lat2 float64, lng2 float64) int {
	radLat1 := lat1 * RAD
	radLat2 := lat2 * RAD
	a := radLat1 - radLat2
	b := (lng1 - lng2) * RAD
	s := 2 * math.Asin(math.Sqrt(math.Pow(math.Sin(a/2), 2)+math.Cos(radLat1)*math.Cos(radLat2)*math.Pow(math.Sin(b/2), 2)))
	s = s * float64(EARTH_RADIUS)
	return int(s)
}

/**
* 获取本机ip
 */
func GetLocalIp() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}
		}
	}

	return "", errors.New("Unable to determine local IP address (non loopback). Exiting.")
}
func GetClientIP(req *http.Request) (client_ip string, err error) {

	if _, ok := req.Header["X-Forwarded-For"]; ok {
		values := req.Header["X-Forwarded-For"]
		client_ip = strings.TrimSpace(values[0])
	} else if value := os.Getenv("X-Forwarded-For"); value != "" {
		values := strings.Split(value, ",")
		client_ip = strings.TrimSpace(values[0])
	} else if req.RemoteAddr != "" {
		client_ip, _, err = net.SplitHostPort(req.RemoteAddr)
	}

	return client_ip, err
}

func GenLogId(req *http.Request) string {
	return GenTraceId(req.Header.Get("TraceId"))
}

func GenTraceId(defaultValue string) string {
	if defaultValue != "" {
		return defaultValue
	}
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	key := fmt.Sprintf("%d", rnd.Int())
	return Md5(key)
}

func Contains(obj interface{}, target interface{}) bool {
	targetTyp := reflect.TypeOf(target)
	if targetTyp == nil {
		return false
	}
	if reflect.TypeOf(obj) == nil {
		return false
	}
	targetValue := reflect.ValueOf(target)
	switch targetTyp.Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < targetValue.Len(); i++ {
			if targetValue.Index(i).Interface() == obj {
				return true
			}
		}
	case reflect.Map:
		if targetValue.MapIndex(reflect.ValueOf(obj)).IsValid() {
			return true
		}
	}
	return false
}

func ContainsInt(sl []int, v int) bool {
	for _, vv := range sl {
		if vv == v {
			return true
		}
	}
	return false
}

func ContainsString(sl []string, v string) bool {
	for _, vv := range sl {
		if vv == v {
			return true
		}
	}
	return false
}

func ContainsBytes(sl [][]byte, v []byte) bool {
	for _, vv := range sl {
		if bytes.Equal(vv, v) {
			return true
		}
	}
	return false
}

// Distinct returns the unique vals of a slice
//
// [1, 1, 2, 3] >> [1, 2, 3]
func Distinct(arr interface{}) (reflect.Value, bool) {
	// create a slice from our input interface
	slice, ok := takeArg(arr, reflect.Slice)
	if !ok {
		return reflect.Value{}, ok
	}

	// put the values of our slice into a map
	// the key's of the map will be the slice's unique values
	c := slice.Len()
	m := make(map[interface{}]bool)
	for i := 0; i < c; i++ {
		m[slice.Index(i).Interface()] = true
	}
	mapLen := len(m)

	// create the output slice and populate it with the map's keys
	out := reflect.MakeSlice(reflect.TypeOf(arr), mapLen, mapLen)
	i := 0
	for k := range m {
		v := reflect.ValueOf(k)
		o := out.Index(i)
		o.Set(v)
		i++
	}

	return out, ok
}

// Intersect returns a slice of values that are present in all of the input slices
//
// [1, 1, 3, 4, 5, 6] & [2, 3, 6] >> [3, 6]
//
// [1, 1, 3, 4, 5, 6] >> [1, 3, 4, 5, 6]
func Intersect(arrs ...interface{}) (reflect.Value, bool) {
	// create a map to count all the instances of the slice elems
	arrLength := len(arrs)
	var kind reflect.Kind
	var kindHasBeenSet bool

	tempMap := make(map[interface{}]int)
	for _, arg := range arrs {
		tempArr, ok := Distinct(arg)
		if !ok {
			return reflect.Value{}, ok
		}

		// check to be sure the type hasn't changed
		if kindHasBeenSet && tempArr.Len() > 0 && tempArr.Index(0).Kind() != kind {
			return reflect.Value{}, false
		}
		if tempArr.Len() > 0 {
			kindHasBeenSet = true
			kind = tempArr.Index(0).Kind()
		}

		c := tempArr.Len()
		for idx := 0; idx < c; idx++ {
			// how many times have we encountered this elem?
			if _, ok := tempMap[tempArr.Index(idx).Interface()]; ok {
				tempMap[tempArr.Index(idx).Interface()]++
			} else {
				tempMap[tempArr.Index(idx).Interface()] = 1
			}
		}
	}

	// find the keys equal to the length of the input args
	numElems := 0
	for _, v := range tempMap {
		if v == arrLength {
			numElems++
		}
	}
	out := reflect.MakeSlice(reflect.TypeOf(arrs[0]), numElems, numElems)
	i := 0
	for key, val := range tempMap {
		if val == arrLength {
			v := reflect.ValueOf(key)
			o := out.Index(i)
			o.Set(v)
			i++
		}
	}

	return out, true
}

// Union returns a slice that contains the unique values of all the input slices
//
// [1, 2, 2, 4, 6] & [2, 4, 5] >> [1, 2, 4, 5, 6]
//
// [1, 1, 3, 4, 5, 6] >> [1, 3, 4, 5, 6]
func Union(arrs ...interface{}) (reflect.Value, bool) {
	// create a temporary map to hold the contents of the arrays
	tempMap := make(map[interface{}]uint8)
	var kind reflect.Kind
	var kindHasBeenSet bool

	// write the contents of the arrays as keys to the map. The map values don't matter
	for _, arg := range arrs {
		tempArr, ok := Distinct(arg)
		if !ok {
			return reflect.Value{}, ok
		}

		// check to be sure the type hasn't changed
		if kindHasBeenSet && tempArr.Len() > 0 && tempArr.Index(0).Kind() != kind {
			return reflect.Value{}, false
		}
		if tempArr.Len() > 0 {
			kindHasBeenSet = true
			kind = tempArr.Index(0).Kind()
		}

		c := tempArr.Len()
		for idx := 0; idx < c; idx++ {
			tempMap[tempArr.Index(idx).Interface()] = 0
		}
	}

	// the map keys are now unique instances of all of the array contents
	mapLen := len(tempMap)
	out := reflect.MakeSlice(reflect.TypeOf(arrs[0]), mapLen, mapLen)
	i := 0
	for key := range tempMap {
		v := reflect.ValueOf(key)
		o := out.Index(i)
		o.Set(v)
		i++
	}

	return out, true
}

// Difference returns a slice of values that are only present in one of the input slices
//
// [1, 2, 2, 4, 6] & [2, 4, 5] >> [1, 5, 6]
//
// [1, 1, 3, 4, 5, 6] >> [1, 3, 4, 5, 6]
func Difference(arrs ...interface{}) (reflect.Value, bool) {
	// create a temporary map to hold the contents of the arrays
	tempMap := make(map[interface{}]int)
	var kind reflect.Kind
	var kindHasBeenSet bool

	for _, arg := range arrs {
		tempArr, ok := Distinct(arg)
		if !ok {
			return reflect.Value{}, ok
		}

		// check to be sure the type hasn't changed
		if kindHasBeenSet && tempArr.Len() > 0 && tempArr.Index(0).Kind() != kind {
			return reflect.Value{}, false
		}
		if tempArr.Len() > 0 {
			kindHasBeenSet = true
			kind = tempArr.Index(0).Kind()
		}

		c := tempArr.Len()
		for idx := 0; idx < c; idx++ {
			// how many times have we encountered this elem?
			if _, ok := tempMap[tempArr.Index(idx).Interface()]; ok {
				tempMap[tempArr.Index(idx).Interface()]++
			} else {
				tempMap[tempArr.Index(idx).Interface()] = 1
			}
		}
	}

	// write the final val of the diffMap to an array and return
	numElems := 0
	for _, v := range tempMap {
		if v == 1 {
			numElems++
		}
	}
	out := reflect.MakeSlice(reflect.TypeOf(arrs[0]), numElems, numElems)
	i := 0
	for key, val := range tempMap {
		if val == 1 {
			v := reflect.ValueOf(key)
			o := out.Index(i)
			o.Set(v)
			i++
		}
	}

	return out, true
}

func takeArg(arg interface{}, kind reflect.Kind) (val reflect.Value, ok bool) {
	val = reflect.ValueOf(arg)
	if val.Kind() == kind {
		ok = true
	}
	return
}

func Empty(obj interface{}) bool {
	typ := reflect.TypeOf(obj)
	if typ == nil {
		return true
	}
	value := reflect.ValueOf(obj)

	switch typ.Kind() {
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice, reflect.String:
		if value.Len() == 0 {
			return true
		} else {
			return false
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if value.Int() == 0 {
			return true
		} else {
			return false
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if value.Uint() == 0 {
			return true
		} else {
			return false
		}
	case reflect.Bool:
		if value.Bool() == false {
			return true
		} else {
			return false
		}
	default:
		return false
	}

	return false
}

func Count(target interface{}) int {
	targetTyp := reflect.TypeOf(target)
	if targetTyp == nil {
		return 0
	}
	value := reflect.ValueOf(target)
	switch value.Kind() {
	case reflect.Map, reflect.Slice, reflect.Array, reflect.String:
		return value.Len()
	default:
		return 0

	}
}

func ToString(value interface{}) (s string) {
	if reflect.TypeOf(value) == nil {
		return ""
	}

	switch v := value.(type) {
	case bool:
		s = strconv.FormatBool(v)
	case float32:
		s = strconv.FormatFloat(float64(v), 'f', -1, 32)
	case float64:
		s = strconv.FormatFloat(v, 'f', -1, 64)
	case int:
		s = strconv.FormatInt(int64(v), 10)
	case int8:
		s = strconv.FormatInt(int64(v), 10)
	case int16:
		s = strconv.FormatInt(int64(v), 10)
	case int32:
		s = strconv.FormatInt(int64(v), 10)
	case int64:
		s = strconv.FormatInt(v, 10)
	case uint:
		s = strconv.FormatUint(uint64(v), 10)
	case uint8:
		s = strconv.FormatUint(uint64(v), 10)
	case uint16:
		s = strconv.FormatUint(uint64(v), 10)
	case uint32:
		s = strconv.FormatUint(uint64(v), 10)
	case uint64:
		s = strconv.FormatUint(v, 10)
	case string:
		s = v
	case []byte:
		s = string(v)
	case error:
		s = v.Error()
	default:
		b, _ := json.Marshal(v)
		s = string(b)
	}
	return s
}

func ToStringMap(data interface{}) map[string]interface{} {
	valMap := make(map[string]interface{})
	if data == nil {
		return valMap
	}
	dataStr, ok := data.(string)
	if !ok {
		return valMap
	}
	err := json.Unmarshal([]byte(dataStr), &valMap)
	if err != nil {
		return valMap
	}
	return valMap
}

func ToStringSetMap(data interface{}) map[string]int {
	valMap := make(map[string]int)
	if data == nil {
		return valMap
	}
	dataStr, ok := data.(string)
	if !ok {
		return valMap
	}
	valList := strings.Split(dataStr, ",")
	for _, val := range valList {
		valMap[val] = 1
	}
	return valMap
}

func GetStringAbParams(layerId string, abTest map[string]interface{}) string {
	if info, ok := abTest[layerId]; ok {
		params, ok := info.(map[string]interface{})
		if !ok {
			return ""
		}
		if group, ok := params["group"].(string); ok {
			return group
		}
	}
	return ""
}

func GetFloatMapValue(abTest map[string]interface{}, paths ...string) (float64, bool) {
	val := getInterfaceMapValue(abTest, paths...)
	if ans, ok := val.(float64); ok {
		return ans, ok
	}
	return InvalidFloat, false
}

func GetStringMapValue(mapInfo map[string]interface{}, paths ...string) (string, bool) {
	val := getInterfaceMapValue(mapInfo, paths...)
	if ans, ok := val.(string); ok {
		return ans, ok
	}
	return "", false
}

func GetStringSliceMapValue(mapInfo map[string]interface{}, paths ...string) []string {
	val := getInterfaceMapValue(mapInfo, paths...)
	if ans, ok := val.([]string); ok {
		return ans
	}
	return nil
}

func getInterfaceMapValue(tempMap map[string]interface{}, paths ...string) interface{} {
	if len(paths) == 0 {
		return tempMap
	}
	val := getMapInterface(tempMap, paths[0])
	if val == nil {
		return nil
	}
	paths = paths[1:]
	for _, path := range paths {
		if newVal, ok := val.(map[string]interface{}); ok {
			val = newVal[path]
		}
	}
	return val
}

func getMapInterface(tempMap map[string]interface{}, key string) interface{} {
	if val, ok := tempMap[key]; ok {
		return val
	}
	return nil
}

func UTCTransLocal(utcTime string) string {
	t, _ := time.Parse("2006-01-02T15:04:05.000+08:00", utcTime)
	return t.Local().Format("2006-01-02 15:04:05")
}

func LocalTransUTC(localTime string) string {
	t, _ := time.ParseInLocation("2006-01-02 15:04:05", localTime, time.Local)
	return t.UTC().Format("2006-01-02T15:04:05.000+08:00")
}
