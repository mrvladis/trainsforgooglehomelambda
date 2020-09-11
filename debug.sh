#!/bin/zsh
cd trainsforgooglehomelambda 
GOOS=linux go build -o main 
cd ..
sam local start-api -d 5986 --debugger-path /Users/vlaned/go/delv/  --debug-args "-delveAPI=2" --env-vars env.json
