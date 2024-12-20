#!/bin/bash
# Script for updating dependencies to the latest commit

cd ~/r/turdgl
latest_commit=$(git rev-parse HEAD)
cd -
echo "Getting turdgl"
go get "github.com/z-riley/turdgl"@$latest_commit

cd ~/r/turdserve
latest_commit=$(git rev-parse HEAD)
cd -
echo "Getting turdserve"
go get "github.com/z-riley/turdserve"@$latest_commit
