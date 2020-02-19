#!/bin/bash

if [ $1 ];then
	DATA=`cat $1`
else
	DATA=`cat -`
fi

len=${#DATA}
for((i=$len-1;i>=0;i--)); do rev="$rev${DATA:$i:1}"; done
echo $rev
