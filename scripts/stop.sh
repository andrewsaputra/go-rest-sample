#!/bin/bash

PID=$(pidof go-rest-sample)

if [ -z "$PID" ]
then
    echo "application not running"
else 
    kill -9 $PID
fi

