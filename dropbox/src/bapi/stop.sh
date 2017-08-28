#!/bin/bash
/bin/kill -9 `ps -ef|grep bapi |grep -v grep |awk -F ' ' '{print $2}'`
