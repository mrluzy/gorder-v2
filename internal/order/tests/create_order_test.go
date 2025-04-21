package tests

import (
	"context"
	"fmt"
	sw "github.com/mrluzy/gorder-v2/common/client/order"
	_ "github.com/mrluzy/gorder-v2/common/config"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	ctx    = context.Background()
	server = fmt.Sprintf("http://%s/api", viper.GetString("order.http-addr"))
)

func TestMain(m *testing.M) {
	before()
	m.Run()
}

func before() {
	logrus.Infof("server=%s", server)
}

func TestCreateOrder_invalidParams(t *testing.T) {
	response := getResponse(t, "123", sw.PostCustomerCustomerIdOrderJSONRequestBody{
		CustomerId: "123",
		Items:      nil,
	})
	assert.Equal(t, 200, response.StatusCode())
	assert.Equal(t, 2, response.JSON200.Errno)
}

func TestCreateOrder_success(t *testing.T) {
	response := getResponse(t, "123", sw.PostCustomerCustomerIdOrderJSONRequestBody{
		CustomerId: "123",
		Items: []sw.ItemWithQuantity{
			{
				Id:       "prod_S8Kcs4H65OdCxg",
				Quantity: 10,
			},
			{
				Id:       "prod_S8KbxrN4dAEE0y",
				Quantity: 5,
			},
		},
	})
	assert.Equal(t, 200, response.StatusCode())
	assert.Equal(t, 0, response.JSON200.Errno)
}

func getResponse(t *testing.T, customerID string, body sw.PostCustomerCustomerIdOrderJSONRequestBody) *sw.PostCustomerCustomerIdOrderResponse {
	// 加上 t.Helper()，可以让测试框架报告错误时，
	//把锅甩给真正写测试逻辑的地方，
	//让你快速找到“谁”用了“哪个 helper”出问题，而不是去翻工具函数代码。
	t.Helper()
	client, err := sw.NewClientWithResponses(server)
	if err != nil {
		t.Fatal(err)
	}
	response, err := client.PostCustomerCustomerIdOrderWithResponse(ctx, customerID, body)

	if err != nil {
		t.Fatal(err)
	}
	return response
}
