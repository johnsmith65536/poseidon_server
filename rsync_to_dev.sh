#!/bin/bash
rsync -avr --exclude .git --exclude output/ --exclude .idea/ . root@112.74.188.186:/root/gopath/src/poseidon