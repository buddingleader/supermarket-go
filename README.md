# supermarket-go
It's completely designed for my own supermarket

#操作手册
s               # 查看当天账本
s1              # 查看全部账簿
s2              # 显示指定日期的数据库
d               # 删除上一条记录
d1              # 删除当天账本
d2              # 指定条形码删除当天记录，注意只能按时间正序删除一条


#两种记账模式
1.直接输入条形码 
2.直接输入价格


#更新脚本
cd $GOPATH/src/supermarket-go
git pull
go build main.go
./main.exe