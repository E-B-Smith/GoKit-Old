

//----------------------------------------------------------------------------------------
//
//                                                                         AppleReceipt.go
//                                                  ServerUtil: Basic API server utilities
//
//                                                                E.B.Smith, December 2015
//                        -©- Copyright © 2015-2016 Edward Smith, all rights reserved. -©-
//
//----------------------------------------------------------------------------------------


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


func errorFromValidationStatus(status int64) error {
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


//  All dates are in RFC 3339 date format


type AppleInAppReceipt struct {
    Quantity                string  `json:"quantity"`
    ProductID               string  `json:"product_id"`
    TransactionID           string  `json:"transaction_id"`
    OriginalTransactionID   string  `json:"original_transaction_id"`
    PurchaseDate            string  `json:"purchase_date"`
    OriginalPurchaseDate    string  `json:"original_purchase_date"`
    ExpirationDate          string  `json:"expires_date"`
    CancellationDate        string  `json:"cancellation_date"`
    AppItemID               string  `json:"app_item_id"`
    VersionExternalID       string  `json:"version_external_identifier"`
    WebOrderLineItemID      string  `json:"web_order_line_item_id"`
}


type AppleReceipt struct {
    BundleID                    string  `json:"bundle_id"`
    ApplicationVersion          string  `json:"application_version"`
    OriginalApplicationVersion  string  `json:"original_application_version"`
    OriginalPurchaseDateMS      string  `json:"original_purchase_date_ms"`
    InAppReceipts               []AppleInAppReceipt  `json:"in_app"`
}


type AppleReceiptResponse struct {
    Status              int64
    Receipt             AppleReceipt
}


func ValidateAppleReceiptTransaction(receiptData []byte, transactionID string) (*AppleInAppReceipt, error) {

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
    appleJSONResponse, error := ioutil.ReadAll(response.Body)
    Log.Debugf("Response: \n%s\n.", string(appleJSONResponse))

    if response.StatusCode != 200 {
        error = fmt.Errorf("Error HTTP Status %s.", response.Status)
        Log.LogError(error)
        return nil, error
    }

    var receiptResponse AppleReceiptResponse
    decoder := json.NewDecoder(bytes.NewReader(appleJSONResponse))
    error =  decoder.Decode(&receiptResponse)
    if error != nil {
        Log.LogError(error)
        return nil, error
    }
    if receiptResponse.Status != 0 {
        error = errorFromValidationStatus(receiptResponse.Status)
        Log.LogError(error)
        return nil, error
    }

    for _, inAppReceipt := range receiptResponse.Receipt.InAppReceipts {
        if inAppReceipt.TransactionID == transactionID {
            return &inAppReceipt, nil
        }
    }

    return nil, nil
}

