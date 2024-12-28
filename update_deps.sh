#!/bin/bash
# Script for updating dependencies to the latest commit

cd ~/r/gogl
latest_commit=$(git rev-parse HEAD)
cd -
echo "Getting gogl"
echo $latest_commit
go get "github.com/z-riley/gogl"@$latest_commit

exit 1

cd ~/r/servesyouright
latest_commit=$(git rev-parse HEAD)
cd -
echo "Getting servesyouright"
go get "github.com/z-riley/servesyouright"@$latest_commit
