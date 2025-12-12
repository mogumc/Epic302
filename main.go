package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"sync"
	"syscall"
)

var originalBackup string

var cdnDomains = map[string][]string{
	"Amazon": {
		"download.epicgames.com",
		"download2.epicgames.com",
		"download3.epicgames.com",
		"download4.epicgames.com",
	},
	"Akamai": {
		"epicgames-download1.akamaized.net",
	},
	"Fastly": {
		"fastly-download.epicgames.com",
	},
	"Cloudflare": {
		"cloudflare.epicgamescdn.com",
	},
	"Tencent": {
		"epicgames-download1-1251447533.file.myqcloud.com",
	},
}

func getHostsPath() string {
	if runtime.GOOS == "windows" {
		return `C:\Windows\System32\drivers\etc\hosts`
	}
	return "/etc/hosts"
}

func prepareHostsModifications(selectedCDN string) (newLines []string, addedDomains []string, backupContent string, err error) {
	hostsFile := getHostsPath()
	data, err := os.ReadFile(hostsFile)
	if err != nil {
		err = fmt.Errorf("无法读取 hosts 文件，请以管理员身份运行: %v", err)
		return
	}

	backupContent = string(data)
	var lines []string
	existing := make(map[string]bool)

	scanner := bufio.NewScanner(strings.NewReader(backupContent))
	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)
		lines = append(lines, line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}
		fields := strings.Fields(trimmed)
		if len(fields) >= 2 {
			ip := fields[0]
			host := fields[1]
			if strings.HasPrefix(ip, "127.") || ip == "localhost" {
				existing[host] = true
			}
		}
	}

	for name, domains := range cdnDomains {
		if name == selectedCDN {
			continue
		}
		for _, domain := range domains {
			if !existing[domain] {
				newLine := fmt.Sprintf("127.0.0.1\t%s", domain)
				lines = append(lines, newLine)
				addedDomains = append(addedDomains, domain)
			}
		}
	}

	newLines = lines
	return
}

func writeHosts(content string) error {
	path := getHostsPath()
	return os.WriteFile(path, []byte(content), 0644)
}

func restoreHosts(original string) {
	fmt.Print("\r🔄 正在还原 hosts 文件... ")
	err := writeHosts(original)
	if err != nil {
		fmt.Printf("失败！请手动还原。\n")
		log.Printf("❌ 还原 hosts 失败: %v", err)
	} else {
		fmt.Printf("已完成\n")
	}
}

func main() {
	fmt.Println("🎮 Epic302 - 极简代理模式")
	fmt.Println("⚠️  注意：被选中的 CDN 不会劫持，其余全部指向本地")
	fmt.Println("请选择代理后端 CDN：")
	fmt.Println("----------------------------------------")

	names := []string{"Amazon", "Akamai", "Fastly", "Cloudflare", "Tencent"}
	for i, name := range names {
		fmt.Printf("%d. %s\n", i+1, name)
	}
	fmt.Println()

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("输入编号 (1-5): ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	var choice int
	_, err := fmt.Sscanf(input, "%d", &choice)
	if err != nil || choice < 1 || choice > 5 {
		log.Fatal("❌ 输入无效")
	}

	selectedCDN := names[choice-1]
	fmt.Printf("\n✅ 已选择: %s\n", selectedCDN)

	modifiedLines, addedDomains, originalHosts, err := prepareHostsModifications(selectedCDN)
	if err != nil {
		log.Fatalf("❌ 读取 hosts 失败: %v", err)
	}

	modifiedContent := strings.Join(modifiedLines, "\n")
	originalBackup = originalHosts

	err = writeHosts(modifiedContent)
	if err != nil {
		log.Fatalf("❌ 修改 hosts 失败，请以管理员身份运行: %v", err)
	}

	if len(addedDomains) > 0 {
		fmt.Printf("📝 已添加以下域名到 127.0.0.1：\n")
		for _, d := range addedDomains {
			fmt.Printf("   → %s\n", d)
		}
	} else {
		fmt.Printf("🟢 hosts 已包含所需条目，无需修改\n")
	}

	var once sync.Once
	restoreFunc := func() {
		once.Do(func() {
			restoreHosts(originalBackup)
		})
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		restoreFunc()
		os.Exit(0)
	}()

	backendDomains := cdnDomains[selectedCDN]
	if len(backendDomains) == 0 {
		log.Fatalf("❌ 未找到 %s 的后端域名", selectedCDN)
	}
	targetDomain := backendDomains[0]

	proxy := newReverseProxy(targetDomain)

	fmt.Printf("\n🚀 本地代理已启动！\n")
	fmt.Printf("📍 监听端口: :80\n")
	fmt.Printf("🎯 代理目标: %s (%s)\n", selectedCDN, targetDomain)
	fmt.Printf("🛑 按 Ctrl+C 退出，程序将自动还原 hosts\n\n")

	log.Printf("[Epic302] 服务启动 | Selected=%s | Target=%s", selectedCDN, targetDomain)
	if err := http.ListenAndServe(":80", proxy); err != nil {
		restoreFunc()
		log.Fatalf("❌ 启动失败: %v\n请确认是否以管理员权限运行", err)
	}
}
