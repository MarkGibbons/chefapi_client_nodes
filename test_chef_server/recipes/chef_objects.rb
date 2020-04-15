# Add chef objects to the server for testing

execute 'get the ssl certificate for the chef server' do
  command 'knife ssl fetch -s https://localhost -u pivotal -k /etc/opscode/pivotal.pem'
end

execute 'Get go modules' do
	command 'export GOPATH=/root/go && cd /root/go/src/github.com/MarkGibbons && go get ./...'
end

execute 'Create organizations' do
  command '/root/go/src/testapi/bin/organization'
end

execute 'Create nodes' do
  command '/root/go/src/testapi/bin/node'
end
