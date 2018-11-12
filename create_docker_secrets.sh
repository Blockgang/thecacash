#!/bin/bash
source .env

declare -A secrets

secrets=( 
    ['mysql_root_password']="8drRNG8RWw9FjzeJuavbY6f9" \
    ['database_password']="6rNhNAPY7yXf" \
)

function setSecretFile(){
    secretName=$1
    secret=$2
    echo ${secret} > secrets/${secretName}.txt
}


for secretName in ${!secrets[@]}; do
    secret=${secrets[${secretName}]}
    docker secret rm ${secretName} 2> /dev/null
    setSecretFile ${secretName} ${secret}
done
ls -al secrets/


