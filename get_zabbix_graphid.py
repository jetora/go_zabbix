#!/usr/bin/env python
# -*- coding: utf-8 -*-

import os
import sys
#import socket
import json
import urllib2
#import random
#import linecache
#from multiprocessing.dummy import Pool as ThreadPool

class Get_zabbix_graphid:
    def __init__(self):
        self.url = "http://xxx/api_jsonrpc.php"
        self.header = {"Content-Type": "application/json"}
        self.user = "xxx"
        self.passwd = "xxx"
        self.authID = self.user_login()

    def user_login(self):
        data = json.dumps(
                {
                    "jsonrpc": "2.0",
                    "method": "user.login",
                    "params": {
                        "user": self.user,
                        "password": self.passwd
                        },
                    "id": 0
                    })
        request = urllib2.Request(self.url,data)
        for key in self.header:
            request.add_header(key,self.header[key])
        try:
            result = urllib2.urlopen(request)
        except Exception as e:
            print "Auth Failed, Please Check Your Name And Password:",e.code
        else:
            response = json.loads(result.read())
            result.close()
            authID = response['result']
            return authID

    def get_data(self,data,hostip=""):
        request = urllib2.Request(self.url,data)
        for key in self.header:
            request.add_header(key,self.header[key])
        try:
            result = urllib2.urlopen(request)
        except Exception as e:
            if hasattr(e, 'reason'):
                print 'We failed to reach a server.'
                print 'Reason: ', e.reason
            elif hasattr(e, 'code'):
                print 'The server could not fulfill the request.'
                print 'Error code: ', e.code
            return 0
        else:
            response = json.loads(result.read())
            result.close()
            return response

    def get_hostid(self,hostip):
        data = json.dumps(
                {
                    "jsonrpc": "2.0",
                    "method": "host.get",
                    "params": {
                        "output":["name","status","host","groups"],
                        #"output":"extend",
                        "selectGroups":"extend",
                        "filter": {"ip": [hostip]}
                        },
                    "auth": self.authID,
                    "id": 1
                })
        res = self.get_data(data)['result']
        #print res
        hostid = '0'
        if len(res):
            hostid = res[0]['hostid']
        return hostid

    def get_graphid(self,hostip):
        hostid = self.get_hostid(hostip)
        graphid = '0'
        if hostid != '0':
            data = json.dumps(
                    {
                    "jsonrpc": "2.0",
                    "method": "graph.get",
                    "params": {
                        "output": "extend",
                        "hostids": hostid,
                        "sortfield": "name"
                    },
                    "auth": self.authID,
                    "id": 1
                })
            res = self.get_data(data)['result']
            graph_list = {}
            if len(res):
                for graphid in res:
                    #if graphid['name'] in ['Mysql_RW', 'Network', 'MySQL_Thread', 'Seconds_Behind_Master', 'Network_MySQL', 'Cpu_Load', 'Tcp_conect', 'CPU_Used']:
                    if graphid['name'] in ['Mysql_RW', 'MySQL_Thread', 'Cpu_Load', 'CPU_Used']:
                        graph_list[graphid['name']] = graphid['graphid']
                #return graph_list
                #print graph_list
                return graph_list
            else:
                print 'errrr'

#def main(ip):
#    zabbix = zabbixtools()
#    zabbix.get_graphid(ip)#172.28.252.60
#http://zabbixm.mysql.jddb.com/chart2.php?graphid=105685&period=43200&stime=20160721165952&sid=fbd67a5b4e03858b&width=900&height=300&box=box.jpg
#stime=20160721165952 开始时间
#period=43200  多长时间s
#graphid=105685 图的id
#sid=fbd67a5b4e03858b    self.authID
#width=900&height=300  宽高
#box=box.jpg图片的格式，也可以是png

if __name__ == "__main__":
    try:
        import optparse
        parse=optparse.OptionParser(usage='" usage : %prog [options] arg1 "', version="%prog 1.0")
        parse.add_option('-i', '--ip', dest = 'ip', type = str, help = 'Input host IP')
        parse.add_option('-v', help='version 1.0')
        parse.set_defaults(v = 1.0)
        options,args=parse.parse_args()
        zg = Get_zabbix_graphid()
        print zg.get_graphid(options.ip)
    except Exception as e:
        print str(e)
