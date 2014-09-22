package main

import (
	"time"
	"net/http"
	"net/url"
    "crypto/hmac"
    "crypto/sha1"
    "encoding/base64"
	)


func uploadFileOfSizeToAWS(pathkey string,  filesize int64) error {
	return nil
	}
	
/*
function uploadFileOfSizeToAWS()
	{
	local filename="$1"
	local filesize="$2"
	local filepath="${awsBackupDrive}"/"${filename}"
	local resource=/"${awsS3Bucket}"/"${filename}"
	local contentType="binary/octet-stream"
	local md5Signature=$(md5 -q "${filepath}")
	local dateValue=$(date)
	local stringToSign="PUT\n${md5Signature}\n${contentType}\n${dateValue}\n${resource}"
	local signature=$(printf "${stringToSign}" | openssl sha1 -hmac "${AWSSecretAccessKey}" -binary | base64)

	curl -X PUT -T "${filepath}" \
		-H "Host: ${awsS3Bucket}.s3.amazonaws.com" \
		-H "Date: ${dateValue}" \
		-H "Content-Type: ${contentType}" \
		-H "Content-Length: ${filesize}" \
		-H "Content-MD5: ${md5Signature}" \
		-H "Authorization: AWS ${AWSAccessKeyId}:${signature}" \
			https://"${awsS3Bucket}".s3.amazonaws.com/"${filename}"
	}
*/


func AWSSignatureForString(string string) string {
    key := []byte(globalAWSAccessSecret)
    h := hmac.New(sha1.New, key)
    h.Write([]byte(string))
    return base64.StdEncoding.EncodeToString(h.Sum(nil))
	}


func listAWSObjectsWithPrefixAndMarkerToFile(prefix string, marker string, filepath string) {
	var query string = ""

	if len(prefix) > 0 {
		query += url.QueryEscape(prefix)
		}

	if len(marker) > 0 {
		if len(query) > 0 {
			query += "&"
			}
		query += "marker=" + url.QueryEscape(marker)
		}

	query = "?" + query
//	var resource string = "/" + globalAWSBackupBucket + "/" + query;
	dateString   := time.Now().Format(time.RFC822Z);
	stringToSign := "GET\n\n\n" + dateString + "\n/" + globalAWSBackupBucket + "/"
	signature    := AWSSignatureForString(stringToSign);
	urlString    := "https://" + globalAWSBackupBucket + ".s3.amazonaws.com/" + query

	//	ToDo: Add retries --

	request, error := http.NewRequest("GET", urlString, nil)
	request.Header.Add("Host", globalAWSBackupBucket + ".s3.amazonaws.com")
	request.Header.Add("Date", dateString)
	request.Header.Add("Authorization", "AWS " + globalAWSAccessKeyID + ":" + signature)

	client := &http.Client{ Timeout:30.0 }
	response, error := client.Do(request)
	response.Body.Close()

//	response=$(
//	curl -X GET --insecure \
//		--output "${filepath}" \
//		--retry 3 --silent --show-error \
//		-H "Host: ${awsBackupBucket}.s3.amazonaws.com" \
//		-H "Date: ${dateValue}" \
//		-H "Authorization: AWS ${AWSAccessKeyID}:${signature}" \
//		-D - https://"${awsBackupBucket}".s3.amazonaws.com/"${query}"
//		)


//	statusCode=$?
//	httpStatus=$(stringDelimitedByStrings "${response}" 'HTTP/1.1 ' 'x-amz')
//	echo "Status: $httpStatus" >&2 
//	
//	result="OK"
//	if (( $statusCode != 0 )); then
//		result="Curl:$statusCode"
//	elif [[ "$httpStatus" != "200 OK" ]]; then
//		result="HTTPError: ${httpStatus}"
//		fi
//	echo -e "\nResult: $result\n" >&2
//	echo "$result"
	}

