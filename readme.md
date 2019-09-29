# 华东师大图书馆ic管理系统自动预订(中北+闵行)

## 需求

预约中北图书馆4楼单间（小黑屋）和闵行图书馆裙楼五楼单间。每天21点开放后天预约，写了个脚本实现自动预约。

## 使用说明（推荐Golang版本）
### Python版本
* clone/下载 项目到本地
* python3运行环境，没有google下安装
* 安装requests库，pip3 install requests
* 修改conf.json为自己的配置
  * sid 自己的学号
  * password 公共数据库密码
  * roomNo 想要预订的房间编号，例如C421
  * startTime 开始时间
  * endTime 结束时间（持续时间不能大于4小时）
* 20:59:00 在当前目录下运行 python3 booking.py，即可开始预订

配合 crontab 定时任务可以实现每天自动预订某一个房间。例如:

  `59 20 * * * /home/admin/anaconda3/bin/python /home/admin/ECNU-library-booking/booking.py &>> /home/admin/booking.log`

ps： 某些房间也存在大神的自动脚本预订，此脚本干不过，遇到这种情况请更换预订房间。

### Golang 版本
* 直接去下载页面，下载对应的操作系统版本。
* 修改conf.json文件
* 21:59:00 运行程序，即可预定后天的房间

## 实现方法

使用 python 中的 requests 库来模拟浏览器登录，之后计算当前时间到21：00点（学校系统开放时间）的时间间隔，使用 sleep 进行等待，时间一到发出请求，省去了手动预订时在浏览器中选择时间段的时间，因而比浏览器快。

## 声明

图书馆小黑屋属于公共资源，请合理利用。若预订成功后，不能前去，请及时删除预订，以供别人使用。
