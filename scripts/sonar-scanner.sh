#!/bin/sh

ls

#if [ ! -d "/tmp/sonarscanner" ]; then
    wget -O /tmp/sonarscanner.zip https://binaries.sonarsource.com/Distribution/sonar-scanner-cli/sonar-scanner-cli-5.0.1.3006-linux.zip
    unzip /tmp/sonarscanner.zip -d /tmp/sonarscanner
#fi


for dir in "/tmp/sonarscanner/*"; do
    sonar_dir=$dir
    break
done


if [ -n "$CODEBUILD_SOURCE_VERSION" ] && [ $(echo "$CODEBUILD_SOURCE_VERSION" | cut -c 1-2) = "pr" ]; then
    PR_NUMBER=$(echo "$CODEBUILD_SOURCE_VERSION" | cut -d "/" -f 2)
    $sonar_dir/bin/sonar-scanner \
        -Dsonar.pullrequest.key=$PR_NUMBER \
        -Dsonar.go.coverage.reportPaths=coverage.out
else
    $sonar_dir/bin/sonar-scanner
fi