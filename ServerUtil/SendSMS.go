//  SendSMS  -  Send an SMS to a phone using Twilio.
//
//  E.B.Smith  -  March, 2015


package ServerUtil


import (
    "errors"
    "strings"
    "io/ioutil"
    "net/url"
    "net/http"
    "encoding/json"
    "encoding/base64"
    "../log"
)


var (
    twilioAccount         = "AC5f879594f852bb0052429f9ac0090ec0"
    twilioAuthToken       = "7de42f76a24d47854c4c909e7789b2bd"
    twilioEncodedAuth     = ""
    twilioFromNumber      = "+14153196030"
    twilioUrlString       = "https://api.twilio.com/2010-04-01/Accounts/"+twilioAccount+"/Messages.json"
)


func SendSMS(toNumber string, message string) error {
    //  Send a Twilio message to a phone number.

    if twilioEncodedAuth == "" {
        twilioEncodedAuth = "Basic "+base64.StdEncoding.EncodeToString([]byte(twilioAccount+":"+twilioAuthToken))
    }

    formValues := url.Values{}
    formValues.Set("From", twilioFromNumber)
    formValues.Set("To", toNumber)
    formValues.Set("Body", message)

    request, _ := http.NewRequest("POST", twilioUrlString, strings.NewReader(formValues.Encode()))
    request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
    request.Header.Add("Authorization", twilioEncodedAuth)
    //log.Debugf("Request:\n%v", *request)

    client := &http.Client{}
    httpResponse, error := client.Do(request)
    if error != nil {
        log.Errorf("Can't contact Twilio: %v.", error)
        httpResponse.Body.Close()
        return error
    }
    defer httpResponse.Body.Close()
    body, error := ioutil.ReadAll(httpResponse.Body)
    if error != nil {
        log.Warningf("SMS response error: %v.", error)
        return error
    }

    //log.Debugf("Response body: %s.", string(body))

    type TwilioResponse struct {
        code float64
        detail string
    }
    var response map[string]interface{}
    error = json.Unmarshal(body, &response)
    //log.Debugf("%v", response)
    if error != nil {
        log.Errorf("SMS response error: %v.", error)
        return error
    }

    if  response["code"] != nil {
        log.Errorf("SMS error %1.0f: %s", response["code"], response["detail"])
        error = errors.New(message)
    } else {
        log.Debugf("Sent SMS OK.")
    }

    return error
}

