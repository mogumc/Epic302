# 🚀 Epic302

> ⚡ Epic 游戏优化的本地代理转发，劫持并转发请求至固定 CDN 节点。

![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go)
![Platform](https://img.shields.io/badge/Platform-Windows-blue)
![License](https://img.shields.io/github/license/mogumc/Epic302)

---

## 🔍 项目简介

**Epic302** 是一个轻量级 Go 编写的代理工具，旨在解决 Epic Games 启动器在中国大陆或其他网络环境下下载缓慢、连接失败的问题。

它通过以下方式实现无缝加速：

1. 修改系统 `hosts` 文件，将官方域名（如 `download.epicgames.com`）指向本机（`127.0.0.1`）；
2. 在本地启动 HTTP 代理服务器监听 80 端口；
3. 接收请求后，保持原始路径、查询参数和 Header 不变，将其转发到实际 CDN 地址；
4. 实现零配置加速体验，无需修改客户端设置。

> 名字来源：“Epic” + HTTP 重定向状态码 `302` → **Epic302**

---

## ✨ 使用须知

### 请以管理员模式启动 

**重定向后速度不满意可使用[UsbEAm Hosts Editor](https://www.dogfight360.com/blog/18627/)修改对应hosts来指定IP,指定后暂停再开始下载即可** 
暂不支持https流量处理，经过观察epic大部分时间使用http流量。
仅提供Windows下编译,且仅支持64位Win10以上系统。其他系统请自行编译。

## 🩺 常见问题

**下载长时间为0Mbps？**  
访问启动器目录下``Epic Games\Launcher\Portal\Config``中的``DefaultEngine.ini``修改HTTP组以下部分  
```
[HTTP]
HttpConnectionTimeout=30
```  
将数值减小(建议 10)后重启启动器

**速度不理想？**  
使用[UsbEAm Hosts Editor](https://www.dogfight360.com/blog/18627/)修改对应hosts来指定IP下载

**EPIC返回tls握手失败？**  
访问启动器目录下``Epic Games\Launcher\Engine\Config``中的``DefaultEngine.ini``寻找HTTP组下是否存在  
```
[HTTP]
bUseNullHttp=true
```  
若存在修改该行为``false``后重启启动器

## 📌 AI生成使用说明

本程序部分代码使用AI生成

## 📄 开源许可

[MIT](https://mit-license.org/)


