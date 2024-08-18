#!/bin/bash
# Script for updating turdgl to the latest commit

cd ~/repo/turdgl
latest_commit=$(git rev-parse HEAD)
cd - 

go get "github.com/z-riley/turdgl"@$latest_commit
