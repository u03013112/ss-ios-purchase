package ios

import (
	"bytes"
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	uuid "github.com/nu7hatch/gouuid"
	"github.com/tidwall/gjson"
	pb "github.com/u03013112/ss-pb/ios"
)

// Srv ：服务
type Srv struct{}

var appleTestURL = "https://sandbox.itunes.apple.com/verifyReceipt"
var appleURL = "https://buy.itunes.apple.com/verifyReceipt"

// Purchase : 支付确认
func (s *Srv) Purchase(ctx context.Context, in *pb.PurchaseRequest) (*pb.PurchaseReply, error) {
	url := appleURL
	value, exsist := os.LookupEnv(strings.ToUpper("SANDBOX"))
	if exsist && value == "TRUE" {
		url = appleTestURL
	}

	postStr := `{"password":"098157441c3e4f71b5c85f62c164b7f3",` +
		`"receipt-data":"` + in.Data +
		`","exclude-old-transactions":true}`
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Post(url, "application/json;charset=utf-8", bytes.NewBufferString(postStr))
	if err != nil {
		return &pb.PurchaseReply{}, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return &pb.PurchaseReply{}, err
	}
	str := string(body)
	// fmt.Print(str)
	status := gjson.Get(str, "status").Int()
	if status != 0 {
		return &pb.PurchaseReply{}, errors.New(gjson.Get(str, "status").String())
	}
	v := gjson.Get(str, "latest_receipt_info")
	if len(v.Array()) > 0 {
		b := v.Array()[0].Raw
		ex := gjson.Get(b, "expires_date_ms").Int()
		t := time.Unix(ex/int64(1000), 0)
		user, err := getUserByToken(in.Token)
		if err != nil {
			return nil, err
		}
		// user.updateExpireDate(t, in.Data)
		user.updateExpireDate(t, "too long")

		recordBills(str,user.UUID)

		return &pb.PurchaseReply{
			ExpiresDate: ex / int64(1000),
		}, nil
	}
	return &pb.PurchaseReply{}, nil
}

// Login ：ios 登录
func (s *Srv) Login(ctx context.Context, in *pb.LoginRequest) (*pb.LoginReply, error) {
	user := getOrCreateUserByUUID(in.Uuid)
	token, _ := uuid.NewV4()
	user.updateToken(token.String())
	return &pb.LoginReply{
		Token:       user.Token,
		ExpiresDate: user.ExpireDate.Unix(),
	}, nil
}

// GetConfig :
func (s *Srv) GetConfig(ctx context.Context, in *pb.GetConfigRequest) (*pb.GetConfigReply, error) {
	user, err := getUserByToken(in.Token)
	if err != nil {
		fmt.Print(err.Error())
		return nil, err
	}
	if user.ExpireDate.Unix() > time.Now().Unix() {
		ret := new(pb.GetConfigReply)
		if config, err := grpcGetConfig(); err != nil {
			fmt.Print(err.Error())
			return ret, err
		} else {
			ret.IP = config.IP
			ret.Port = config.Port
			ret.Method = config.Method
			ret.Passwd = config.Passwd
			ret.ExpiresDate = user.ExpireDate.Unix()
		}
		return ret, nil
	}
	return &pb.GetConfigReply{}, errors.New("ExpireDate")
}
