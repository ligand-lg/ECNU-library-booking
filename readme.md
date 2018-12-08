# 华东师大中北图书馆ic管理系统自动预订

## 功能
预约图书馆小黑屋。每天9点开放预约使用浏览器很难预约上，写了一个 python 脚本，配置好之后，20：59点运行，比浏览器成功率大。

## 使用说明
* clone 项目到本地
* 确保装有 python3，并且安装了requests库，没有使用pip3 install requests安装
* 修改 booking.py 中 get_config，每一个字段都有写具体说明
* 20：59 在当前目录下运行 python3 booking.py，即可预订成功
* 配合 crontab 可以实现每天自动预订某一个房间
ps： 某些房间也存在大神的自动脚本预订，此脚本干不过，遇到这种情况请更换预订房间。

## 实现方法
使用 python 中的 requests 库来模拟浏览器登录，之后计算当前时间到21：00点（学校系统开放时间）还要多久，使用 sleep 进行等待，时间一到发出请求，省去了浏览器中选择时间段的时间，因而比浏览器快。

## 声明
图书馆小黑屋属于公共资源，请合理利用。预订成功后，来不急使用，请及时删除预订，以供别人使用。
