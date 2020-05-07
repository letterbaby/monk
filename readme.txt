1、单pong测试
start.bat

2、多pong测试
开启第一个pong
pongsvr.exe -cfg pong.json -log ponglog.json

开启第二个pong
pongsvr.exe -cfg pong2.json -log pong2log.json

开启ping
pingsvr.exe -cfg ping.json -log pinglog.json

3、DB测试
pong.cfg修改对应的mysql、redis配置

4、讨论群
QQ:336598527
