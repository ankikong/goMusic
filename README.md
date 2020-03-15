# goMusic

> 网易居然把低音号的版权全丢了!因此怒更此项目

## 用法

``` bash
> goMusic -kw "关键词"  #搜索关键词(搜索目前只支持音乐)
> goMusic -url "url"   #分享链接获取歌曲/视频
```

- url方法还是有点问题,等待修复

## 目前状态

| 平台 | 支持链接 | 128 | 320 | 无损 |
| :---- | :---- | :---- | :---- | :--- |
| 网易云 | `http://music.163.com/song?id=404543135` `http://music.163.com/song/28481683` | 支持 | 支持 | 部分支持|
|酷狗|暂不支持|支持|支持|部分支持|
|QQ|`https://i.y.qq.com/v8/playsong.html?songid=105603683`|支持|不支持|不支持|
|Bilibili|`https://www.bilibili.com/video/av114514` `https://www.bilibili.com/bangumi/media/md1547` `https://www.bilibili.com/bangumi/play/ep28919`|-|-|-|
|zzzfun|`http://www.zzzfun.com/vod-detail-id-193.html`|-|-|-|
## 开发计划

- 哔哩哔哩 ~~视频&番剧&~~ 音频
- ACFUN
- 酷我音乐
- 虾米音乐
- 咪咕音乐

## 状态

- V0.2.2
  - 增加zzzfun
- V0.2.1
  - 增加视频合并
  - 修复番剧分块问题
  - 纪录片应该可以正常下载了
- V0.2.0
  - 增加Bilibili番剧支持
  - 增加进度条
- V0.1.0
  - 修复大量bug
  - 修改结构
  - 增加Bilibili视频
- V0.0.2
  - 增加QQ音乐
- V0.0.1
  - 首次更新

## 参考项目

- [NeteaseCloudMusicApi](https://github.com/Binaryify/NeteaseCloudMusicApi)
- [QQMusicApi](https://github.com/jsososo/QQMusicApi)
