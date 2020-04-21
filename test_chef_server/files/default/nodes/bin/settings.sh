# TODO: Make this a template
export GOPATH=/root/go

# api to the chef server
export CHEFAPICHEFUSER=pivotal
export CHEFAPIKEYFILE=/etc/opscode/pivotal.pem
export CHEFAPICHRURL=https://testhost
export CHEFAPICERTFILE=/home/vagrant/.chef/trusted_certs/testhost.crt

# Login token service
export CHEFAPILOGINPORT=8113
export CHEFAPILOGINCERT=/usr/local/nodes/certs/server.crt
export CHEFAPILOGINKEY=/usr/local/nodes/certs/server.key

# Node authorization service
export CHEFAPIAUTHPORT=9001
export CHEFAPIAUTHURL=http://testhost:${CHEFAPIAUTHPORT} 

# Nodes rest service
export CHEFAPINODECERT=/usr/local/nodes/certs/server.crt
export CHEFAPINODEKEY=/usr/local/nodes/certs/server.key
export CHEFAPINODEPORT=8111

# Orgs rest service
export CHEFAPIORGCERT=/usr/local/nodes/certs/server.crt
export CHEFAPIORGKEY=/usr/local/nodes/certs/server.key
export CHEFAPIORGPORT=8112

# Web rest service
export CHEFAPIWEBCERT=/usr/local/nodes/certs/server.crt
export CHEFAPIWEBKEY=/usr/local/nodes/certs/server.key
export CHEFAPIWEBPORT=8443
