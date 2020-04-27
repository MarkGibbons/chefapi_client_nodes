#!/bin/sh


### BEGIN INIT INFO
# Provides:             chefapi-users
# Required-Start:       $local_fs $remote_fs
# Required-Stop:        $local_fs $remote_fs
# Should-Start:
# Should-Stop:
# Default-Start:        S
# Default-Stop:
# Short-Description:    Chefapi users rest service
### END INIT INFO

PATH="/sbin:/bin/:usr/sbin:/usr/bin:/usr/local/nodes/bin"
NAME="chefapi-users"
DESC="Chefapi users rest service"

# chef client settings
. /usr/local/nodes/bin/settings.sh
running=`ps -ef|grep main|grep ${CHEFAPIUSERPORT}|grep -v grep|awk '{print $2}'`

start()
{
	cd /root/go/src/github.com/MarkGibbons/chefapi_client_users
	go run main.go -restcert ${CHEFAPIUSERCERT} -restkey ${CHEFAPIUSERKEY} -restport ${CHEFAPIUSERPORT} -authurl ${CHEFAPIAUTHURL} </dev/null 1>/tmp/chefapi_users.log 2>&1 &
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
		# start the users rest service
		if [ -z "${running}" ]; then
			start
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
			echo "chefapi users is running"
		else
			echo "chefapi users is not running"
			exit 1
		fi
		;;
esac
exit 0
