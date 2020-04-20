#!/bin/sh


### BEGIN INIT INFO
# Provides:             chefapi-nodes
# Required-Start:       $local_fs $remote_fs
# Required-Stop:        $local_fs $remote_fs
# Should-Start:
# Should-Stop:
# Default-Start:        S
# Default-Stop:
# Short-Description:    Chefapi nodes rest service
### END INIT INFO

PATH="/sbin:/bin/:usr/sbin:/usr/bin:/usr/local/nodes/bin"
NAME="chefapi-nodes"
DESC="Chefapi nodes rest service"

# chef client settings
. /usr/local/nodes/bin/settings.sh
running=`ps -ef|grep main|grep ${CHEFAPINODEPORT}|grep -v grep|awk '{print $2}'`

start()
{
	cd /root/go/src/github.com/MarkGibbons/chefapi_client_nodes
	go run main.go -restcert ${CHEFAPINODECERT} -restkey ${CHEFAPINODEKEY} -restport ${CHEFAPINODEPORT} -authurl ${CHEFAPIAUTHURL} </dev/null 1>/tmp/chefapi_nodes.log 2>&1 &
}
stop()
{
	echo "${running}" | xargs -L1 kill -9
}

case "${1}" in
	restart)
		stop
		start
		;;
	start)
		# start the nodes rest service
		if [ -z "${running}" ]; then
			start()
			exit $?
		fi
		;;
	stop)
		if [ -n "${running}" ]; then
			stop
		fi
		;;
	status)
		if [ -n "${running}" ]; then
			echo "chefapi nodes is running"
		else
			echo "chefapi nodes is not running"
			exit 1
		fi
		;;
esac
exit 0
