#!/bin/bash
rsync -avr --exclude .git --exclude output/ --exclude .idea/ . johnsmith@112.74.188.186:/home/johnsmith/gopath/src/poseidon