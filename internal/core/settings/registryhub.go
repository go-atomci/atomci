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
	hostUrl := fmt.Sprintf("%s://%s", schema, strings.TrimRight(basicUrl, "/"))
	resp, err := http.Get(hostUrl)
	if err != nil {
		return errors.New(fmt.Sprintf("%s访问异常:%s", hostUrl, err.Error()))
	}
	url := fmt.Sprintf("%s/v2/", hostUrl)
	resp, err = http.Get(url)
	if err != nil {
		return errors.New(fmt.Sprintf("%s访问异常:%s", url, err.Error()))
	}
	if resp.StatusCode != 401 {
		return errors.New(fmt.Sprintf("%s不支持V2版本访问", url))
	}

	var authBaseUrl, authFullUrl string
	var queryParams []string

	//get Auth Info
	auth := resp.Header.Get("Www-Authenticate")
	if strings.HasPrefix(auth, "Basic") {
		authBaseUrl = url
		queryParams = nil
	} else if strings.HasPrefix(auth, "Bearer") {
		//Bearer realm="https://dockerauth.cn-hangzhou.aliyuncs.com/auth",service="registry.aliyuncs.com:cn-hangzhou:26842"
		kvArr := strings.Split(strings.TrimPrefix(auth, "Bearer "), ",")

		for _, i2 := range kvArr {
			temp := strings.Split(i2, "=")
			if strings.HasPrefix(i2, "realm") {
				authBaseUrl = strings.Trim(temp[1], "\"")
			} else {
				queryParams = append(queryParams, temp[0]+"="+strings.Trim(temp[1], "\""))
			}
		}
	} else {
		return errors.New("basicUrl is incorrect")
	}
	if len(queryParams) > 0 {
		authFullUrl = authBaseUrl + "?" + strings.Join(queryParams, "&")
	} else {
		authFullUrl = authBaseUrl
	}
	req, err := http.NewRequest("GET", authFullUrl, nil)
	if err != nil {
		return errors.New("账号或密码不正确")
	}
	req.SetBasicAuth(username, password)
	client := http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return errors.New("账号或密码不正确")
	}
	//body, err := ioutil.ReadAll(resp.Body)
	//fmt.Println(string(body))
	return nil
}
