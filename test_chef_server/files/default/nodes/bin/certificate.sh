#!/bin/bash
openssl req -x509 -newkey rsa:2048 -keyout server.key -out server.crt -config config.txt -subj "/C=US/ST=Washington/L=Seattle/O=Development/OU=Dev/CN=testhost.com"
# chmod 400 server.*
