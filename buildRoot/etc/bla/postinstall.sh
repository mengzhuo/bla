#!/bin/sh

set -e

[ -f /etc/default/bla ] && . /etc/default/bla

if [ -x /bin/systemctl ]; then
	systemctl daemon-reload
	systemctl restart bla
fi
