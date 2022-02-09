package settings

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

func TryLoginRegistry(basicUrl, username, password string, insecure bool) error {
	var schema string
	if insecure {
		schema = "http"
	} else {
		schema = "https"
	}
	url := fmt.Sprintf("%s://%s/v2/", schema, strings.TrimRight(basicUrl, "/"))
	resp, err := http.Get(url)
	if err != nil || resp.StatusCode != 401 {
		return errors.New("error basicUrl")
	}
	//get Auth Info
	auth := resp.Header.Get("Www-Authenticate")
	if !strings.HasPrefix(auth, "Bearer") {
		return errors.New("basicUrl is incorrect")
	}
	//Bearer realm="https://dockerauth.cn-hangzhou.aliyuncs.com/auth",service="registry.aliyuncs.com:cn-hangzhou:26842"
	kvArr := strings.Split(strings.TrimPrefix(auth, "Bearer "), ",")

	var authBaseUrl, authFullUrl string
	var queryParams []string
	for _, i2 := range kvArr {
		temp := strings.Split(i2, "=")
		if strings.HasPrefix(i2, "realm") {
			authBaseUrl = strings.Trim(temp[1], "\"")
		} else {
			queryParams = append(queryParams, temp[0]+"="+strings.Trim(temp[1], "\""))
		}
	}
	if len(queryParams) > 0 {
		authFullUrl = authBaseUrl + "?" + strings.Join(queryParams, "&")
	} else {
		authFullUrl = authBaseUrl
	}
	req, err := http.NewRequest("GET", authFullUrl, nil)
	req.SetBasicAuth(username, password)
	if err != nil {
		return errors.New("incorrect username or password")
	}
	client := http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return errors.New("incorrect username or password")
	}
	//body, err := ioutil.ReadAll(resp.Body)
	//fmt.Println(string(body))
	return nil
}
