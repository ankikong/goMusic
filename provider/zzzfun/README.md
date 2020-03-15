# zzzfun

找番时无意间发现的平台，资源还行，所以就做一下

## 加密

目前只测试了客户端的play url加密,网页端有点难搞.反编译一下客户端,很容易就能发现:

``` python
sing=hashlib.md5((str(playID) + "534697").encode()).hexdigest()
```

`videoID` 就是:`http://www.zzzfun.com/vod-detail-id-193.html` 中的`193`

`playID` 要通过`videoID`获取,`http://service-agbhuggw-1259251677.gz.apigw.tencentcs.com/android/video/list_ios?videoId=193&userid=`
`videoID`填进去,访问这个链接,就可以看到每集的playID了

但其实`playID`就是`videoID-{n}-a`,n就是一个数字,表第几集,如`193-1-a`