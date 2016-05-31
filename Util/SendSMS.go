//  SendSMS  -  Send an SMS to a phone using Twilio.
//
//  E.B.Smith  -  March, 2015


package Util


import (
    "fmt"
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
    Log.Debugf("Response body: %s.", string(body))

    type TwilioResponse struct {
        Code    float64
        Message string
        Status  string
    }
    var twilioResponse TwilioResponse
    error = json.Unmarshal(body, &twilioResponse)
    if error != nil { return error }

    Log.Debugf("%+v", twilioResponse)
    if twilioResponse.Code != 0 {
        Log.Errorf("SMS response error: %s.", string(body))
        return fmt.Errorf("%f %+v", twilioResponse.Code, twilioResponse.Message)
    }

    Log.Debugf("Sent SMS OK.")
    return error
}

