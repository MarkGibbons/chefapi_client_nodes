# TODO: Make this a template
export GOPATH=/root/go

# Node authorization service
export CHEFAPIAUTHPORT=9001
export CHEFAPIAUTHURL=http://testhost:${CHEFAPIAUTHPORT} 

# api to the chef server
export CHEFAPICHEFUSER=pivotal
export CHEFAPIKEYFILE=/etc/opscode/pivotal.pem
export CHEFAPICHRURL=https://testhost
export CHEFAPICERTFILE=/home/vagrant/.chef/trusted_certs/testhost.crt

# Nodes rest service
export CHEFAPINODECERT=/usr/local/node/certs/server.crt
export CHEFAPINODEKEY=/usr/local/node/certs/server.key
export CHEFAPINODEPORT=8111

# Orgs rest service
export CHEFAPIORGCERT=/usr/local/node/certs/server.crt
export CHEFAPIORGKEY=/usr/local/node/certs/server.key
export CHEFAPIORGPORT=8112

# Web rest service
export CHEFAPIWEBCERT=/usr/local/node/certs/server.crt
export CHEFAPIWEBKEY=/usr/local/node/certs/server.key
export CHEFAPIWEBPORT=8443
