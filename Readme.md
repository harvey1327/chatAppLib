
# ChatAppLib

Contains libraries required for the chat app platform

## database

Library code for interactions to the database

## messagebroker

Library code for the messagebroker setup and pub/sub

## proto

Library code for geerated protobuff files\
To run generation build docker image using `docker build . -t goprotobuff:v1` and run using `docker run --rm -v ${PWD}/generated:/go/generated -v ${PWD}/protofiles:/go/input goprotobuff:v1`
