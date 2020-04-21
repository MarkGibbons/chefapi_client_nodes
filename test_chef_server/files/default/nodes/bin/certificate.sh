#!/bin/bash
#openssl req -x509 -newkey rsa:2048 -keyout server.key -out server.crt -config config.txt -subj "/C=US/ST=Washington/L=Seattle/O=Development/OU=Dev/CN=testhost.com"
# chmod 400 server.*
openssl genrsa -out server.key 2048
openssl req -new -key server.key -out server.csr -config config.cnf -subj "/C=US/ST=Washington/L=Seattle/O=Development/OU=Dev/CN=localhost"
openssl x509 -req -days 3650 -in server.csr -signkey server.key -out server.crt
