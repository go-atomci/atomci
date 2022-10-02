package notification

import (
	"testing"

	"github.com/astaxie/beego"
)

func init() {
	_ = beego.LoadAppConfig("ini", "./app.unittest.conf")
}

func Test_SEND_SHOULD_NO_ERROR(t *testing.T) {
	// httpmock.Activate()
	// defer httpmock.DeactivateAndReset()

	// httpmock.RegisterResponder("POST", "https://dingtalk.unittest.com",
	// 	func(req *http.Request) (*http.Response, error) {
	// 		resp, err := httpmock.NewJsonResponse(200, "")
	// 		if err != nil {
	// 			return httpmock.NewStringResponse(500, ""), nil
	// 		}
	// 		return resp, nil
	// 	},
	// )

	// var temp notification.INotify

	// mockResult := new(messages.StepCallbackResult)

	// temp = new(notification.Email)

	// assert.NoError(t, temp.Send(*mockResult))

	// temp = new(notification.DingRobot)
	// assert.NoError(t, temp.Send(*mockResult))

	// hit := httpmock.GetTotalCallCount()
	// assert.Equal(t, 1, hit)
}
