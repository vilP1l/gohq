package gohq

import (
	"net/http"
	"bytes"
	"encoding/json"
	"log"
	"io/ioutil"
	"errors"
	"strconv"
	"github.com/smartystreets/go-aws-auth"
	"net/url"
	"fmt"
)

// Request makes (GET/POST/PUT/PATCH/etc..) requests to the HQ API
func (a *Account) Request(method string, urlStr string, data interface{}, auth bool) (response []byte, err error) {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return
	}

	var req *http.Request
	if data == nil {
		req, err = http.NewRequest(method, urlStr, nil)
	} else {
		req, err = http.NewRequest(method, urlStr, bytes.NewBuffer(dataBytes))
	}

	if err != nil {
		return
	}

	if auth {
		req.Header.Set("Authorization", "Bearer "+a.AccessToken)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "hq-viewer/1.4.14 (iPhone; iOS 11.2.2; Scale/3.00)")
	req.Header.Set("Content-Length", strconv.Itoa(len(dataBytes)))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}

	defer func() {
		if resp.Body.Close() != nil {
			log.Println("error closing resp body")
		}
	}()

	response, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	var hqerr HQError
	if err = json.Unmarshal(response, &hqerr); err == nil && hqerr.Error != "" {
		err = errors.New(hqerr.Error)
		if err.Error() == "not authorized" && urlStr != EndpointTokens {
			if a.LoginToken != "" {
				tokens, err_two := a.Tokens()
				if err_two == nil {
					a.AccessToken = tokens.AccessToken
					a.LoginToken = tokens.LoginToken
					a.AuthToken = tokens.AuthToken

					response, err = a.Request(method, urlStr, data, auth)
				}
			}
		}
	}

	// TODO: Add a check to see if HQ ever goes down

	return
}

// Verify a phone number via (SMS or CALL?)
func Verify(number, method string) (verification *Verification, err error) {
	type Data struct {
		Method string `json:"method"`
		Phone  string `json:"phone"`
	}

	a := Account{}
	resp, err := a.Request("POST", EndpointVerifications, Data{Method: method, Phone: number}, false)

	if err != nil {
		return
	}

	err = json.Unmarshal(resp, &verification)
	return
}

// Confirm a verification session
func (v *Verification) Confirm(code string) (auth *Auth, err error) {
	type Data struct {
		Code string `json:"code"`
	}

	a := Account{}
	resp, err := a.Request("POST", EndpointVerifications+v.VerificationID, Data{Code: code}, false)

	if err != nil {
		return
	}

	err = json.Unmarshal(resp, &auth)
	return
}

// Create an account if user Confirm().auth == nil (account doesn't exist)
func (v *Verification) Create(username, referrer, region string) (account *Account, err error) {
	type Data struct {
		Country           string `json:"country"`
		Language          string `json:"language"`
		ReferringUsername string `json:"referringUsername"`
		Username          string `json:"username"`
		VerificationID    string `json:"verificationId"`
	}

	a := Account{}
	resp, err := a.Request("POST", EndpointVerifications+v.VerificationID, Data{Country: region, VerificationID: v.VerificationID, Username: username, Language: "en", ReferringUsername: referrer}, false)

	if err != nil {
		return
	}

	err = json.Unmarshal(resp, &account)
	return
}

// Tokens refreshes new tokens based on the login token
func (a *Account) Tokens() (t *Tokens, err error) {
	type Data struct {
		Token string `json:"token"`
	}

	resp, err := a.Request("POST", EndpointTokens, Data{Token: a.LoginToken}, false)
	if err != nil {
		return
	}

	err = json.Unmarshal(resp, &t)

	return
}

// Me gets updated profile information
func (a *Account) Me() (t *Me, err error) {
	type Data struct {
	}

	resp, err := a.Request("GET", EndpointMe, Data{}, true)
	if err != nil {
		return
	}

	err = json.Unmarshal(resp, &t)

	return
}

// Cashout sends a cashout request to HQ
func (a *Account) Cashout(email string) (cd *CashoutData, err error) {
	type Data struct {
		Email string `json:"email"`
	}

	resp, err := a.Request("POST", EndpointPayouts, Data{Email: email}, true)
	if err != nil {
		return
	}

	err = json.Unmarshal(resp, &cd)

	return
}

// Cashout sends a cashout request to HQ
func (a *Account) RequestDocuments(email, payout, country string) (err error) {
	type Data struct {
		Email string `json:"email"`
		Payout string `json:"payoutEmail"`
		Country string `json:"country"`
	}

	_, err = a.Request("POST", EndpointDocuments(strconv.Itoa(a.UserID)), Data{Email: email, Country:country, Payout:payout}, true)
	if err != nil {
		return
	}

	return
}

// Cashout sends a cashout request to HQ
// TODO: Fix
func (a *Account) RegisterDeviceToken(token string) (err error) {
	type Data struct {
		Token string `json:"token"`
	}

	_, err = a.Request("POST", EndpointDevices, Data{Token:token}, true)
	if err != nil {
		return
	}

	return
}

// Payouts gets all of the past payout data
func (a *Account) Payouts() (pd *PayoutData, err error) {
	resp, err := a.Request("GET", EndpointPayouts, nil, true)
	if err != nil {
		return
	}

	err = json.Unmarshal(resp, &pd)

	return
}

// Schedule
func (a *Account) Schedule() (sd *ScheduleData, err error) {
	resp, err := a.Request("GET", EndpointSchedule, nil, true)
	if err != nil {
		return
	}

	err = json.Unmarshal(resp, &sd)

	return
}

// Weekly runs the makeItRain easter egg
func (a *Account) Weekly() (err error) {
	if _, err = a.Request("POST", EndpointMakeItRain, nil, true); err != nil {
		return
	}

	return
}

// Logout and invalidate a bearer
func (a *Account) RefreshLogin() (lt *LoginTokenResponse, err error) {
	resp, err := a.Request("GET", EndpointToken, nil, true)
	if err != nil {
		return
	}

	err = json.Unmarshal(resp, &lt)

	return
}

// Claim a gift
func (a *Account) Claim(gID string) (err error) {
	resp, err := a.Request("POST", EndpointClaim(gID), nil, true)
	if err != nil {
		return
	}

	fmt.Println(string(resp))

	return
}

// ChangeUsername changes a users username
func (a *Account) ChangeUsername(username string) (ud *UpdateInfo, err error) {
	type Data struct {
		Username string `json:"username"`
	}

	resp, err := a.Request("PATCH", EndpointMe, Data{Username: username}, true)
	if err != nil {
		return
	}

	err = json.Unmarshal(resp, &ud)

	return
}

// SearchUser searches for a user
func (a *Account) SearchUser(username string) (sd *SearchData, err error) {
	resp, err := a.Request("GET", EndpointSearchUser(username), nil, true)
	if err != nil {
		return
	}

	err = json.Unmarshal(resp, &sd)

	return
}

// AddFriend adds a user by id
func (a *Account) AddFriend(uID string) (err error) {
	_, err = a.Request("POST", EndpointFriendRequest(uID), nil, true)
	return
}

// DeleteFriend removes a user from your friend list
func (a *Account) DeleteFriend(uID string) (err error) {
	_, err = a.Request("DELETE", EndpointFriendRequest(uID), nil, true)
	return
}

// Request an AWS session
func (a *Account) RequestAWS() (aws *AWSSession, err error) {
	resp, err := a.Request("GET", EndpointAWS, nil, true)

	if err != nil {
		return
	}

	err = json.Unmarshal(resp, &aws)
	return
}



// Upload to AWS
func (aws *AWSSession) Upload(filename string, data []byte) (err error) {
	req, _ := http.NewRequest("PUT", "https://hypespace-quiz.s3.amazonaws.com/avatars/"+url.QueryEscape(filename), bytes.NewReader(data))
	req.Header.Add("Content-Type", "image/jpeg")
	req.Header.Add("Host", "hypespace-quiz.s3.amazonaws.com")

	awsauth.Sign(req, awsauth.Credentials{Expiration: aws.Expiration, AccessKeyID: aws.AccessKeyID, SecretAccessKey: aws.SecretKey, SecurityToken: aws.SessionToken})

	_, err = http.DefaultClient.Do(req)

	return
}

// Change the profile picture to a profile picture on the AWS path
func (a *Account) ChangeAvatar(awsPath string) (result *UpdateInfo, err error) {
	type Data struct {
		AvatarURL string `json:"avatarUrl"`
	}

	resp, err := a.Request("PUT", EndpointAvatarURL, Data{AvatarURL: awsPath}, true)

	if err != nil {
		return
	}

	err = json.Unmarshal(resp, &result)
	return
}
