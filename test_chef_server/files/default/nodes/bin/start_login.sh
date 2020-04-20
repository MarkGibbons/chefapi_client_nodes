#!/bin/sh


### BEGIN INIT INFO
# Provides:             chefapi-login
# Required-Start:       $local_fs $remote_fs
# Required-Stop:        $local_fs $remote_fs
# Should-Start:
# Should-Stop:
# Default-Start:        S
# Default-Stop:
# Short-Description:    Chefapi login rest service
### END INIT INFO

PATH="/sbin:/bin/:usr/sbin:/usr/bin:/usr/local/nodes/bin"
NAME="chefapi-login"
DESC="Chefapi login rest service"

# chef client settings
. /usr/local/nodes/bin/settings.sh
running=$(ps -ef|grep main|grep "restport ${CHEFAPILOGINPORT}"|grep -v grep|awk '{print $2}')

start()
{
	cd /root/go/src/github.com/MarkGibbons/chefapi_login
	go run main.go -restcert ${CHEFAPILOGINCERT} -restkey ${CHEFAPILOGINKEY} -restport ${CHEFAPILOGINPORT} </dev/null 1>/tmp/chefapi_login.log 2>&1 &

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
		# start the login rest service
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
			echo "chefapi login is running"
		else
			echo "chefapi login is not running"
			exit 1
		fi
		;;
esac
exit 0
