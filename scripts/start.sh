#!/bin/bash

cd /home/ssm-user
env GIN_MODE=release nohup ./go-rest-sample >> log.txt 2>&1 &