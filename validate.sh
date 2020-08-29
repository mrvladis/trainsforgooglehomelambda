#!/bin/zsh
   # Make sure we're in the project directory within our GOPATH
    cd "trainsforgooglehomelambda"
      # Fetch all dependencies
    go get -t ./...
      # Ensure code passes all lint tests
    golint -set_exit_status
      # Check the Go code for common problems with 'go vet'
    go vet .
      # Run all tests included with our application
    go test .
    cd ..
