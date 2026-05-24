package utils

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/Sna-ken/videoweb/internal/dao"
	"github.com/pquerna/otp/totp"
	"github.com/skip2/go-qrcode"
)

/*关于MFA部分：因为在过程中个人觉得文档接口中直接返回secret似乎不太好也没必要，所以我这里进行了一些修改
思路为获取MFAqr->暂存到redis->BindMFA->存入MySQL*/

func GenerateMFA(userID string, username string) (string, error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "videoweb",
		AccountName: username,
		Digits:      6,
		Period:      30,
	})
	if err != nil {
		return "", err
	}

	url := key.URL()
	qrcode, err := qrcode.Encode(url, qrcode.Medium, 256)
	if err != nil {
		return "", err
	}

	base64code := base64.StdEncoding.EncodeToString(qrcode)
	if err := dao.StoreMFATemp(context.Background(), userID, key.Secret()); err != nil {
		return "", err
	}

	return fmt.Sprintf("data:image/png;base64,%s", base64code), nil
}

func ValidateMFA(code string, secret string) bool {
	return totp.Validate(code, secret)
}
