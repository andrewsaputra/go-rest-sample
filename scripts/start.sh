#!/bin/bash

APP_DIR=/home/ssm-user

if [ -d "$APP_DIR" ]
then
    cd $APP_DIR
fi

env GIN_MODE=release nohup ./go-rest-sample >> log.txt 2>&1 &