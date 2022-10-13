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

	if strings.EqualFold("Basic", auth[:5]) {
		authBaseUrl = url
		queryParams = nil
	} else if strings.EqualFold("Bearer", auth[:6]) {
		//Bearer realm="https://dockerauth.cn-hangzhou.aliyuncs.com/auth",service="registry.aliyuncs.com:cn-hangzhou:26842"
		//Bearer realm="https://auth.pkg.coding.net/artifacts-auth/docker/jwt?host=leafly-docker.pkg.coding.net",service="docker"
		kvArr := strings.Split(auth[7:], ",")

		for _, i2 := range kvArr {
			index := strings.Index(i2, "=")
			if index == -1 {
				continue
			} else if strings.EqualFold(strings.Trim(i2[:index], " "), "realm") {
				authBaseUrl = strings.Trim(i2[index+1:], "\"")
			} else {
				queryParams = append(queryParams, i2[:index]+"="+strings.Trim(i2[index+1:], "\""))
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
