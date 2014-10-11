package main

import (
	"io"
	"time"
	"net/http"
	"net/url"
	"crypto/tls"
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


func RFC8222StringFromDate(t time.Time) string {
	return t.Format("Mon, 2 Jan 2006 15:04:05 -0700")
	}


func AWSSignatureForString(string string) string {
    key := []byte(globalAWSAccessSecret)
    h := hmac.New(sha1.New, key)
    h.Write([]byte(string))
    return base64.StdEncoding.EncodeToString(h.Sum(nil))
	}


func listAWSObjectsWithPrefixAndMarker(writer io.Writer, prefix string, marker string) error {
	var query string = ""

	if len(prefix) > 0 {
		query += "prefix=" + url.QueryEscape(prefix)
		}

	if len(marker) > 0 {
		if len(query) > 0 {
			query += "&"
			}
		query += "marker=" + url.QueryEscape(marker)
		}

	query = "?" + query
//	var resource string = "/" + globalAWSBackupBucket + "/" + query;
	dateString   := RFC8222StringFromDate(time.Now())
	stringToSign := "GET\n\n\n" + dateString + "\n/" + globalAWSBackupBucket + "/"
	signature    := AWSSignatureForString(stringToSign);
	urlString    := "https://" + globalAWSBackupBucket + ".s3.amazonaws.com/" + query

	log(AWSLogDebug, "     Date: %v.", dateString)
	log(AWSLogDebug, "   String: %v.", stringToSign)
	log(AWSLogDebug, "Signature: %v.", signature)
	log(AWSLogDebug, "      URL: %v.", urlString)

	//	ToDo: Add retries --

	request, error := http.NewRequest("GET", urlString, nil)
	request.Header.Add("Host", globalAWSBackupBucket + ".s3.amazonaws.com")
	request.Header.Add("Date", dateString)
	request.Header.Add("Authorization", "AWS " + globalAWSAccessKeyID + ":" + signature)

	tr := &http.Transport{ TLSClientConfig: &tls.Config{InsecureSkipVerify: true} }
	client := &http.Client{ Timeout:time.Minute*2.0, Transport: tr }
	response, error := client.Do(request)
	if response.statusCode != 200 {
		error = response.error
	} else
	if error == nil {
		log(AWSLogDebug, "Read %d bytes.", response.ContentLength)
		var n int
		buffer := make([] byte, 1000)
		n, error = response.Body.Read(buffer)
		for n > 0 {
			log(AWSLogDebug, "Writing %d bytes.", n)
			writer.Write(buffer[:n])
			n, error = response.Body.Read(buffer)
			}
		response.Body.Close()
		}
	

	log(AWSLogError, "AWS GET error: %v.", error)
	if error != nil && error == error.EOF {
		error = nil;

	return error;
	}

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
//	}

