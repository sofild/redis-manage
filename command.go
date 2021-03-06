package main

import(
	"fmt"
	"errors"
	"net"
	"encoding/json"
	"strconv"
	"time"
	"crypto/md5"
	"strings"
)

/*
	keys xxx
*/
func getKeys(k string, client net.Conn) ([]string, error){
    if k == "*" {
        k = getAZ()
    }
	cmd := fmt.Sprintf("KEYS %s", k)
	keys, err := exec(cmd, client)
    if err != nil {
        return keys, err
    }
    return filterKeys(keys), err
}

/*
	keys type
*/
func getType(k string, client net.Conn) (string, error){
	cmd := fmt.Sprintf("TYPE %s", k)
	data, err := exec(cmd, client)
	return data[0], err
}

/*
	get xxx
*/
func get(k string, client net.Conn) ([]string, error){
	cmd := fmt.Sprintf("GET %s", k)
	data, err := exec(cmd, client)
	return data, err
}


/*
    SET key value [EX seconds] [PX milliseconds] [NX|XX]
*/
func set(k string, value string, client net.Conn) (string, error){
	cmd := fmt.Sprintf("SET %s %s", k, value)
	data, err := exec(cmd, client)
	return data[0], err
}

/*
	hget xxx
*/
func hGet(k string, client net.Conn) ([]string, error){
	cmd := fmt.Sprintf("HGETALL %s", k)
	data, err := exec(cmd, client)
	json_data := make(map[string]string)
	var key string
	for k, v := range data {
		if k % 2 == 0 {
			key = v
		}else{
			json_data[key] = v
		}
	}
	json_byte, _ := json.Marshal(json_data)
	json_str := string(json_byte)
	json_arr := []string{json_str}
	return json_arr, err
}

/*
	hset k field value
*/
func hSet(k string, field string, value string, client net.Conn) (int, error) {
	cmd := fmt.Sprintf("HSET %s %s %s", k, field, value)
	fmt.Println(cmd)
	r, err := exec(cmd, client)
	if err != nil {
		return -1, err
	}
	d := 0
	d, err = strconv.Atoi(r[0])
	if err != nil {
		return -1, err
	}
	return d, nil
}

/*
	HDEL key field [field ...]
*/
func hDel(k string, field string, client net.Conn) (int, error) {
	cmd := fmt.Sprintf("HDEL %s %s", k, field)
	r, err := exec(cmd, client)
	if err != nil {
		return -1, err
	}
	d := 0
	d, err = strconv.Atoi(r[0])
	if err != nil {
		return -1, err
	}
	return d, nil
}

/*
	DEL key [key ...]
*/
func del(k string, client net.Conn) (int, error) {
	cmd := fmt.Sprintf("DEL %s", k)
	r, err := exec(cmd, client)
	if err != nil {
		return -1, err
	}
	d := 0
	d, err = strconv.Atoi(r[0])
	if err != nil {
		return -1, err
	}
	return d, nil
}

/*
	LLEN key
*/
func lLen(k string, client net.Conn) (int, error) {
	cmd := fmt.Sprintf("LLEN %s", k)
	r, err := exec(cmd, client)
	if err != nil {
		return 0, err
	}
	d := 0
	d, err = strconv.Atoi(r[0])
	if err != nil {
		return d, err
	}
	return d, nil
}

/*
	LRANGE key start stop
*/
func lRange(k string, s int, e int, client net.Conn) ([]string, error){
	cmd := fmt.Sprintf("LRANGE %s %d %d", k, s, e)
	data, err := exec(cmd, client)
	return data, err
}

/*
	LPUSH key value [value ...]
*/
func lPush(k string, value string, client net.Conn) (int, error) {
	cmd := fmt.Sprintf("LPUSH %s %s", k, value)
	r, err := exec(cmd, client)
	if err != nil {
		return 0, err
	}
	d := 0
	d, err = strconv.Atoi(r[0])
	if err != nil {
		return d, err
	}
	return d, nil
}

/*
	LSET key index value
*/
func lSet(k string, index int, value string, client net.Conn) (string, error){
	cmd := fmt.Sprintf("LSET %s %d %s", k, index, value)
	data, err := exec(cmd, client)
	return data[0], err
}

/*
	LSET key index value
	LREM key count value
*/
func lDel(k string, index int, client net.Conn) (int, error){
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	h := md5.Sum([]byte(timestamp))
	value := fmt.Sprintf("%x", h)
	_, err := lSet(k, index, value, client)
	if err != nil {
		return 0, err
	}
	cmd := fmt.Sprintf("LREM %s %d %s", k, 1, value)
	r, err := exec(cmd, client)
	if err != nil {
		return 0, err
	}
	d := 0
	d, err = strconv.Atoi(r[0])
	if err != nil {
		return d, err
	}
	return d, nil
}

/*
    SMEMBERS key
*/
func Smembers(k string, client net.Conn) ([]string, error){
	cmd := fmt.Sprintf("SMEMBERS %s", k)
	data, err := exec(cmd, client)
	return data, err
}

/*
    SADD key member [member ...]
*/
func sAdd(k string, value string, client net.Conn) (int, error) {
	cmd := fmt.Sprintf("SADD %s %s", k, value)
	r, err := exec(cmd, client)
	if err != nil {
		return 0, err
	}
	d := 0
	d, err = strconv.Atoi(r[0])
	if err != nil {
		return d, err
	}
	return d, nil
}

/*
    SREM key member [member ...]
*/
func sRem(k string, value string, client net.Conn) (int, error) {
    cmd := fmt.Sprintf("SREM %s %s", k, value)
    r, err := exec(cmd, client)
    if err != nil {
        return 0, err
    }
    d := 0
    d, err = strconv.Atoi(r[0])
    if err != nil {
        return d, err
    }
    return d, nil
}

/*
    Smod ...
    SREM key member [member ...]
    SADD key member [member ...]
*/
func sMod(k string, old_value string, value string, client net.Conn) (int, error) {
    _, err := sRem(k, old_value, client)
    if err != nil {
        return 0, err
    }
    return sAdd(k, value, client)
}

/*
	value
*/
func getValue(k string, client net.Conn) ([]string, error){
	var data []string
	kType, err := getType(k, client)
	if err != nil {
		return data, err
	}
	switch kType {
	case "none":
		err = errors.New(k + " is not exists.")
		return data, err
	case "hash":
		return hGet(k, client)
	case "list":
		var len int = 0
		len, err = lLen(k, client)
		if err != nil {
			return data, err
		}
		return lRange(k, 0, len, client)
	case "set":
        return Smembers(k, client)
	case "string","zset":
		return get(k, client)
	}
	return data, nil
}

/*
  判断是否是二进制
 */
func isBinary(key string) bool {
    for _, char := range []rune(key) {
        if char == 0 {
            return true
        }
        if char > 126 || (char < 32 && char != 8 && char != 9 && char != 10 && char != 13) {
            return true
        }
    }
    return false
}

/*
 filter keys
 */
func filterKeys(keys []string) []string {
    var ks []string
    for _, value := range keys {
        if !isBinary(value) {
            ks = append(ks, value)
        }
    }
    return ks
}

/*
    get a-z
*/
func getAZ() string {
    var p []string
    for i := 97; i <= 122; i++ {
        p = append(p, string(i))
    }
    str := strings.Join(p, "|")
    str = fmt.Sprintf("[%s]*", str)
    return str
}

