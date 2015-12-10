//  ValidateAppleReceipt  -  Validate a receipt from an Apple in-app purchase.
//
//  E.B.Smith  -  December, 2015


package ServerUtil


import (
    "fmt"
    "bytes"
    "net/http"
    "io/ioutil"
    "encoding/json"
    "encoding/base64"
    "violent.blue/GoKit/Log"
)


//----------------------------------------------------------------------------------------
//                                                                    ValidateAppleReceipt
//----------------------------------------------------------------------------------------


func errorFromValidationStatus(status int32) error {
    if status == 0 { return nil }

    s := ""
    switch status {
    case 21000: s = "The App Store could not read the JSON object you provided."
    case 21002: s = "The data in the receipt-data property was malformed or missing."
    case 21003: s = "The receipt could not be authenticated."
    case 21004: s = "The shared secret you provided does not match the shared secret on file for your account."
    case 21005: s = "The receipt server is not currently available."
    case 21006: s = "This receipt is valid but the subscription has expired. When this status code is returned to your server, the receipt data is also decoded and returned as part of the response."
    case 21007: s = "This receipt is from the test environment, but it was sent to the production environment for verification. Send it to the test environment instead."
    case 21008: s = "This receipt is from the production environment, but it was sent to the test environment for verification. Send it to the production environment instead."
    default:    s = "Unknown error."
    }

    return fmt.Errorf("%s (%d)", s, status)
}


//  eDebug -- Fix this!
type AppleReceipt struct {
    Status              int32
    Receipt             map[string]string
    LatestReceipt       string
    LatestReceiptInfo   map[string]string
    Amount              int32
}


func ValidateAppleReceiptTransaction(receiptData []byte, transactionID string) (*AppleReceipt, error) {

    if len(receiptData) == 0 {
        return nil, errorFromValidationStatus(21002)
    }

    receipt := make(map[string]string)
    receipt["receipt-data"] = base64.StdEncoding.EncodeToString(receiptData)
    jsonBytes, error := json.Marshal(receipt)
    if error != nil {
        Log.LogError(error)
        return nil, error
    }

    appleURL := "https://buy.itunes.apple.com/verifyReceipt"
    if true {
        appleURL = "https://sandbox.itunes.apple.com/verifyReceipt"
    }

    request, error := http.NewRequest("POST", appleURL, bytes.NewBuffer(jsonBytes))
    request.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    response, error := client.Do(request)
    if error != nil {
        Log.LogError(error)
        return nil, error
    }
    defer response.Body.Close()

    if response.StatusCode != 200 {
        error = fmt.Errorf("Error HTTP Status %s.", response.Status)
        Log.LogError(error)
        return nil, error
    }
    body, _ := ioutil.ReadAll(response.Body)
    decoder := json.NewDecoder(bytes.NewBuffer(body))

    var validatedReceipt AppleReceipt
    error =  decoder.Decode(&validatedReceipt)
    if error == nil  && validatedReceipt.Status != 0 {
        error = errorFromValidationStatus(validatedReceipt.Status)
    }

    return &validatedReceipt, error
}

