#!/bin/sh


### BEGIN INIT INFO
# Provides:             chefapi-node-auth
# Required-Start:       $local_fs $remote_fs
# Required-Stop:        $local_fs $remote_fs
# Should-Start:
# Should-Stop:
# Default-Start:        S
# Default-Stop:
# Short-Description:    Chefapi node auth rest service
### END INIT INFO

PATH="/sbin:/bin/:usr/sbin:/usr/bin:/usr/local/nodes/bin"
NAME="chefapi-node-auth"
DESC="Chefapi node auth rest service"

# chef client settings
. /usr/local/nodes/bin/settings.sh
search="restport ${CHEFAPIAUTHPORT}"
echo -E "ps -ef|grep main|grep ""${search}""|grep -v grep|awk '{print $2}'"
running=`ps -ef|grep main|grep ""${search}""|grep -v grep|awk '{print $2}'`

case "${1}" in
	start)
		# start the node auth rest service
		if [ -z "${running}" ]; then
			cd /root/go/src/github.com/MarkGibbons/chefapi_node_auth
			go run main.go -restcert ${CHEFAPIAUTHCERT} -restkey ${CHEFAPIAUTHKEY} -restport ${CHEFAPIAUTHPORT} </dev/null 1>/tmp/chefapi_node_auth.log 2>&1 &
			exit $?
		fi
		;;
	stop)
		if [ -n "${running}" ]; then
			echo "${running}" | xargs -L1 kill -9 
		fi
		;;
	status)
		if [ -n "${running}" ]; then
			echo "chefapi node auth is running"
		else
			echo "chefapi node auth is not running"
			exit 1
		fi
		;;
esac
exit 0
