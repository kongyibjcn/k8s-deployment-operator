#!/bin/bash

while true
do

output=`curl http://192.168.99.100:30417 -s`

if [ "$output" = "Welcome to the Green Docker!" ];then

   echo -e "\033[32m Welcome to the Green Docker! \033[0m"

else
   echo -e "\033[34m Welcome to the Blue Docker! \033[0m"
fi

sleep 2s

done