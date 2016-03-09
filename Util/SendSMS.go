//  SendSMS  -  Send an SMS to a phone using Twilio.
//
//  E.B.Smith  -  March, 2015


package Util


import (
    "errors"
    "strings"
    "io/ioutil"
    "net/url"
    "net/http"
    "encoding/json"
    "encoding/base64"
    "violent.blue/GoKit/Log"
)


var (
    twilioAccount         = "AC5f879594f852bb0052429f9ac0090ec0"
    twilioAuthToken       = "7de42f76a24d47854c4c909e7789b2bd"
    twilioEncodedAuth     = ""
    twilioFromNumber      = "+14153196030"
    twilioUrlString       = "https://api.twilio.com/2010-04-01/Accounts/"+twilioAccount+"/Messages.json"
)


func SendSMS(toNumber string, message string) error {
    Log.LogFunctionName()

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
    //Log.Debugf("Request:\n%v", *request)

    client := &http.Client{}
    httpResponse, error := client.Do(request)
    if error != nil {
        Log.Errorf("Can't contact Twilio: %v.", error)
        httpResponse.Body.Close()
        return error
    }
    defer httpResponse.Body.Close()
    body, error := ioutil.ReadAll(httpResponse.Body)
    if error != nil {
        Log.Warningf("SMS response error: %v.", error)
        return error
    }

    //Log.Debugf("Response body: %s.", string(body))

    type TwilioResponse struct {
        code float64
        detail string
    }
    var response map[string]interface{}
    error = json.Unmarshal(body, &response)
    //Log.Debugf("%v", response)
    if error != nil {
        Log.Errorf("SMS response error: %v.", error)
        return error
    }

    if  response["code"] != nil {
        Log.Errorf("SMS error %1.0f: %s", response["code"], response["detail"])
        error = errors.New(message)
    } else {
        Log.Debugf("Sent SMS OK.")
    }

    return error
}

