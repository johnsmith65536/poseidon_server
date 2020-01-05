#!/bin/bash

rsync -avr --exclude .git --exclude output/ --exclude .idea/ . johnsmith@192.168.6.128:~/go/src/poseidon
