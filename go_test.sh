#!/bin/bash
export CHEFAPICHEFUSER=testuser
export CHEFAPIKEYFILE=test/key.pem
export CHEFAPICHRURL=https://localhost

go test
