#!/usr/bin/python
# -*- coding:utf-8 -*-
import sys
import datetime
import cookielib, urllib2,urllib
#from get_zabbix_graphid import Get_zabbix_graphid 
import get_zabbix_graphid


class Down_graph(object):
    def __init__(self):
        self.url = 'http://xxx/index.php'
        self.name = 'xxx'
        self.password = 'xxx'
        #初始化的时候生成cookies
        cookiejar = cookielib.CookieJar()
        urlOpener = urllib2.build_opener(urllib2.HTTPCookieProcessor(cookiejar))
        values = {"name":self.name,'password':self.password,'autologin':1,"enter":'Sign in'}
        data = urllib.urlencode(values)
        request = urllib2.Request(self.url, data)
        try:
            uo = urlOpener.open(request,timeout=10)
            self.urlOpener=urlOpener
        except Excepiton, e:
            print e
    def GetGraph(self, url, img_name, values, image_dir, ip):
        key=values.keys()
        if "graphid"  not in key:
            print "print input graphid"
            sys.exit(1)
        #以下if 是给定默认值
        if  "period" not in key :
            #默认获取一天的数据，单位为秒
            values["period"]=86400
        if "stime" not in key:
            #默认为当前时间开始
            values["stime"]=datetime.datetime.now().strftime('%Y%m%d%H%M%S')
        if "width" not in key:
            values["width"]=800
        if "height" not in key:
            values["height"]=200
        data=urllib.urlencode(values)
        request = urllib2.Request(url,data)
        url = self.urlOpener.open(request)
        image = url.read()
        #print len(image)
        #imagename="%s%s.png" % (image_dir, values["graphid"])
        imagename="%s%s-%s.png" % (image_dir, ip, img_name)
        with open(imagename, 'wb') as f:
            f.write(image)
        return imagename
    
    def main(self, ip, period, stime):
        #http://zabbixm.mysql.jddb.com/chart2.php?graphid=161483&period=3600&width=916&height=362
        #此url是获取图片是的，请注意饼图的URL 和此URL不一样，请仔细观察！
        gr_url="http://zabbixm.mysql.jddb.com/chart2.php?"
        #登陆URL
        #用于图片存放的目录
        image_dir='/export/zhaochen/script/zabbix/'
        zg = get_zabbix_graphid.Get_zabbix_graphid()
        graphid_dic = zg.get_graphid(ip)
        file_path = ''
        for img_name, graphid in graphid_dic.items():
        #图片的参数，该字典至少传入graphid。
            values={"graphid": graphid, "period": period, "stime":stime, "width": 800, "height": 200}
            file_path += (self.GetGraph(gr_url, img_name, values, image_dir, ip)) + ',' 
        return file_path
        
