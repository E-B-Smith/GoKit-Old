#!/bin/bash


function percentEscapeString()
	{
	perl -lpe 's/([^A-Za-z0-9])/sprintf("%%%02X", ord($1))/seg' <<<"$1"
	}

twiliokey=AC6d456d2e0a49aad0d6ec1a92ddbd925f
twiliosecret=3191a48e7faa540b4c7b6e84c9a82402
twiliofrom=16072694306
#twilioto=6503916685   # Rohit
#twilioto=2404856540   # Sap
#twilioto=4086379251  # Hemant
twilioto=4156152570  # Edward
encodedauth=$(echo "$twiliokey:$twiliosecret" | base64)
body="McD is still better."
body=$(percentEscapeString "$body")

curl -v -X POST \
    --insecure --retry 3 --silent --show-error \
    -H "Authorization: Basic $encodedauth" \
    --data "From=$twiliofrom&To=$twilioto&Body=$body" \
        "https://api.twilio.com/2010-04-01/Accounts/$twiliokey/Messages"

echo ""
