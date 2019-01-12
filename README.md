![logo](http://pic.yupoo.com/jiananshi/281a0f0e/small.png)

七牛云图片备份 + Yupoo 云图片上传

## Usage

1. 根据系统[下载](https://developer.qiniu.com/kodo/tools/1300/qrsctl?ref=support.qiniu.com)七牛辅助工具至项目根目录
2. `chmod +x qrsctl`
3. `./qrsctl login {{username}} {{password}}`
3. 七牛云图片备份 `go run main.go -m pull -q image -d dist`
4. Yupoo 云图片上传 `go run main.go -m push -d dist`（替换代码中的 Auth Token、Album ID 和 Cookie）

## MIT
