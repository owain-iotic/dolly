#!/bin/bash 

TOKEN="$(i-make-prd-jwt.sh)"

echo "london air"
./twingen -host="uk-london-air.iotics.space" -twinid=did:iotics:iotA3LucoeK9yPXq3cjN8SSizU316ah45yiq -jwt=$TOKEN -output="londonair"


echo "tfl bike"
./twingen -host="uk-london-air.iotics.space" -twinid=did:iotics:iotA3LucoeK9yPXq3cjN8SSizU316ah45yiq -jwt=$TOKEN -output="londonair"

