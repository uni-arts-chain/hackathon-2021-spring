package token

import (
	"github.com/dgrijalva/jwt-go"
	"nftmart/app/global/consts"
	"nftmart/app/global/my_errors"
	"nftmart/app/global/variable"
	"nftmart/app/http/middleware/my_jwt"
	"nftmart/app/model"
	"time"
)

// 创建 userToken 工厂

func CreateUserFactory() *userToken {
	return &userToken{
		userJwt: my_jwt.CreateMyJWT(consts.JwtTokenSignKey),
	}
}

type userToken struct {
	userJwt *my_jwt.JwtSign
}

//生成token
func (u *userToken) GenerateToken(userid int64, username string, phone string, expireAt int64) (tokens string, err error) {

	// 根据实际业务自定义token需要包含的参数，生成token，注意：用户密码请勿包含在token
	customClaims := my_jwt.CustomClaims{
		UserId: userid,
		Name:   username,
		Phone:  phone,
		// 特别注意，针对前文的匿名结构体，初始化的时候必须指定键名，并且不带 jwt. 否则报错：Mixture of field: value and value initializers
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix() - 10,       // 生效开始时间
			ExpiresAt: time.Now().Unix() + expireAt, // 失效截止时间
		},
	}
	return u.userJwt.CreateToken(customClaims)
}
// 生成pub_key_token
func (u *userToken) GeneratePubKeyToken(pubkey string,expireAt int64) (tokens string, err error) {

	// 根据实际业务自定义token需要包含的参数，生成token，注意：用户密码请勿包含在token
	customClaims := my_jwt.CustomPubKeyClaims{
		PubKey: pubkey,
		// 特别注意，针对前文的匿名结构体，初始化的时候必须指定键名，并且不带 jwt. 否则报错：Mixture of field: value and value initializers
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix() - 10,       // 生效开始时间
			ExpiresAt: time.Now().Unix() + expireAt, // 失效截止时间
		},
	}
	return u.userJwt.CreatePubKeyToken(customClaims)
}

// 用户login成功，记录用户token
func (u *userToken) RecordLoginToken(userToken, clientIp string) bool {
	if customClaims, err := u.userJwt.ParseToken(userToken); err == nil {
		userId := customClaims.UserId
		expiresAt := customClaims.ExpiresAt
		return model.CreateUserFactory("").OauthLoginToken(userId, userToken, expiresAt, clientIp)
	} else {
		return false
	}
}
// wallet login success record token info
func (u *userToken) RecordPKeyToken(userToken, clientIp string) bool {
	if customPubKeyClaims, err := u.userJwt.ParsePubKeyToken(userToken); err == nil {
		pubKey := customPubKeyClaims.PubKey
		expiresAt := customPubKeyClaims.ExpiresAt
		return model.CreateUserFactory("").OauthPKeyLoginToken(pubKey, userToken, expiresAt, clientIp)
	} else {
		return false
	}
}
// 刷新token的有效期（默认+3600秒，参见常量配置项）
func (u *userToken) RefreshToken(oldToken, clientIp string) (newToken string, res bool) {

	// 解析用户token的数据信息
	_, code := u.isNotExpired(oldToken)
	switch code {
	case consts.JwtTokenOK, consts.JwtTokenExpired:
		//如果token已经过期，那么执行更新
		if newToken, err := u.userJwt.RefreshToken(oldToken, variable.ConfigYml.GetInt64("Token.JwtTokenRefreshExpireAt")); err == nil {
			if customClaims, err := u.userJwt.ParseToken(newToken); err == nil {
				userId := customClaims.UserId
				expiresAt := customClaims.ExpiresAt
				if model.CreateUserFactory("").OauthRefreshToken(userId, expiresAt, oldToken, newToken, clientIp) {
					return newToken, true
				}
			}
		}
	case consts.JwtTokenInvalid:
		variable.ZapLog.Error(my_errors.ErrorsTokenInvalid)
	}

	return "", false
}

// 销毁token，基本用不到，因为一个网站的用户退出都是直接关闭浏览器窗口，极少有户会点击“注销、退出”等按钮，销毁token其实无多大意义
func (u *userToken) DestroyToken() {

}

// 判断token是否未过期
func (u *userToken) isNotExpired(token string) (*my_jwt.CustomClaims, int) {
	if customClaims, err := u.userJwt.ParseToken(token); err == nil {

		if time.Now().Unix()-customClaims.ExpiresAt < 0 {
			// token有效
			return customClaims, consts.JwtTokenOK
		} else {
			// 过期的token
			return customClaims, consts.JwtTokenExpired
		}
	} else {
		// 无效的token
		return nil, consts.JwtTokenInvalid
	}
}
// 判断token是否未过期
func (u *userToken) isPubKeyNotExpired(token string) (*my_jwt.CustomPubKeyClaims, int) {
	if customPubKeyClaims, err := u.userJwt.ParsePubKeyToken(token); err == nil {

		if time.Now().Unix()-customPubKeyClaims.ExpiresAt < 0 {
			// token有效
			return customPubKeyClaims, consts.JwtTokenOK
		} else {
			// 过期的token
			return customPubKeyClaims, consts.JwtTokenExpired
		}
	} else {
		// 无效的token
		return nil, consts.JwtTokenInvalid
	}
}
// 判断token是否有效（未过期+数据库用户信息正常）
func (u *userToken) IsEffective(token string) bool {
	customClaims, code := u.isNotExpired(token)
	if consts.JwtTokenOK == code {
		//if user_item := Model.CreateUserFactory("").ShowOneItem(customClaims.UserId); user_item != nil {
		if model.CreateUserFactory("").OauthCheckTokenIsOk(customClaims.UserId, token) {
			return true
		}
	}
	return false
}
// 判断token是否有效（未过期+pubkey正常）
func (u *userToken) IsPubKeyEffective(token string) (bool ,string){
	customPubKeyClaims, code := u.isPubKeyNotExpired(token)
	if consts.JwtTokenOK == code {
		//if user_item := Model.CreateUserFactory("").ShowOneItem(customClaims.UserId); user_item != nil {
		if model.CreateUserFactory("").OauthCheckPubKeyTokenIsOk(customPubKeyClaims.PubKey, token) {
			return true,customPubKeyClaims.PubKey
		}
	}
	return false,""
}
