#!/bin/bash

docker run -d -p 8091:8080 -p50000:50000 --name jenkins -v $(pwd)/data:/var/jenkins_home colynn/jenkins:2.277.1-lts-alpine
