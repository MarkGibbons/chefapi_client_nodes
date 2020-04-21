#
# Manage the chefapi nodes application services
#

directory '/usr/local/node/certs' do
  recursive true
end

remote_directory '/usr/local/nodes' do
  source 'nodes'
  mode '0755'
end

file '/usr/local/nodes/bin/start_login.sh' do
  mode '0755'
end

file '/usr/local/nodes/bin/start_node_auth.sh' do
  mode '0755'
end

file '/usr/local/nodes/bin/start_nodes.sh' do
  mode '0755'
end

file '/usr/local/nodes/bin/start_organizations.sh' do
  mode '0755'
end

link '/etc/init.d/chefapi-login' do
  to '/usr/local/nodes/bin/start_login.sh'
end

link '/etc/init.d/chefapi-node-auth' do
  to '/usr/local/nodes/bin/start_node_auth.sh'
end

link '/etc/init.d/chefapi-nodes' do
  to '/usr/local/nodes/bin/start_nodes.sh'
end

link '/etc/init.d/chefapi-organizations' do
  to '/usr/local/nodes/bin/start_organizations.sh'
end

# start login
service 'chefapi-login' do
  provider Chef::Provider::Service::Init
  supports status: true
  action [:restart]
end

# start organizations
service 'chefapi-organizations' do
  provider Chef::Provider::Service::Init
  supports status: true
  action [:restart]
end

# start node auth
service 'chefapi-node-auth' do
  provider Chef::Provider::Service::Init
  supports status: true
  action [:restart]
end

# start nodes
service 'chefapi-nodes' do
  provider Chef::Provider::Service::Init
  supports status: true
  action [:restart]
end
