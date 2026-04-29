package proxy

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// RefererProxy 反防盗链代理
// 通过修改请求头绕过源站的防盗链限制
type RefererProxy struct {
	client *http.Client
}

// NewRefererProxy 创建代理实例
func NewRefererProxy() *RefererProxy {
	return &RefererProxy{
		client: &http.Client{
			// 不自动跟随重定向
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		},
	}
}

// ProxyRequest 代理请求
// 修改 Referer 和 User-Agent，绕过防盗链
func (p *RefererProxy) ProxyRequest(targetURL string, w http.ResponseWriter, r *http.Request) error {
	// 解析目标 URL
	parsedURL, err := url.Parse(targetURL)
	if err != nil {
		return fmt.Errorf("无效的URL: %v", err)
	}

	// 创建新请求
	proxyReq, err := http.NewRequest(r.Method, targetURL, r.Body)
	if err != nil {
		return fmt.Errorf("创建请求失败: %v", err)
	}

	// 设置请求头
	p.setHeaders(proxyReq, parsedURL, r)

	// 发送请求
	resp, err := p.client.Do(proxyReq)
	if err != nil {
		return fmt.Errorf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 处理 M3U8 文件
	if strings.HasSuffix(targetURL, ".m3u8") {
		return p.handleM3U8Response(resp, w, parsedURL)
	}

	// 转发普通响应
	return p.forwardResponse(resp, w)
}

// setHeaders 设置请求头
func (p *RefererProxy) setHeaders(req *http.Request, targetURL *url.URL, originalReq *http.Request) {
	// 1. 设置 Referer（伪装成从源站访问）
	req.Header.Set("Referer", fmt.Sprintf("%s://%s/", targetURL.Scheme, targetURL.Host))

	// 2. 设置 User-Agent（伪装成浏览器）
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")

	// 3. 设置 Origin
	req.Header.Set("Origin", fmt.Sprintf("%s://%s", targetURL.Scheme, targetURL.Host))

	// 4. 传递 Range 头（支持断点续传）
	if rangeHeader := originalReq.Header.Get("Range"); rangeHeader != "" {
		req.Header.Set("Range", rangeHeader)
	}

	// 5. 设置 Accept
	if strings.HasSuffix(targetURL.Path, ".m3u8") {
		req.Header.Set("Accept", "application/vnd.apple.mpegurl,*/*")
	} else {
		req.Header.Set("Accept", "*/*")
	}
}

// handleM3U8Response 处理 M3U8 响应
// 需要重写其中的 URL，让所有分片也通过代理
func (p *RefererProxy) handleM3U8Response(resp *http.Response, w http.ResponseWriter, baseURL *url.URL) error {
	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// 重写 M3U8 内容
	content := string(body)
	rewrittenContent := p.rewriteM3U8URLs(content, baseURL)

	// 设置响应头
	w.Header().Set("Content-Type", "application/vnd.apple.mpegurl")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(resp.StatusCode)

	// 写入重写后的内容
	_, err = w.Write([]byte(rewrittenContent))
	return err
}

// rewriteM3U8URLs 重写 M3U8 中的 URL（含 URI= 属性）
func (p *RefererProxy) rewriteM3U8URLs(content string, baseURL *url.URL) string {
	// 当前文件的目录（带尾部 /）
	dir := baseURL.Path[:strings.LastIndex(baseURL.Path, "/")+1]

	toProxy := func(raw string) string {
		raw = strings.TrimSpace(raw)
		var full string
		switch {
		case strings.HasPrefix(raw, "http://") || strings.HasPrefix(raw, "https://"):
			full = raw
		case strings.HasPrefix(raw, "/"):
			// 绝对路径：直接拼 scheme + host
			full = fmt.Sprintf("%s://%s%s", baseURL.Scheme, baseURL.Host, raw)
		default:
			// 相对路径：拼目录
			full = fmt.Sprintf("%s://%s%s%s", baseURL.Scheme, baseURL.Host, dir, raw)
		}
		return "/api/storage/proxy?url=" + url.QueryEscape(full)
	}

	lines := strings.Split(content, "\n")
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}

		// 重写 #EXT-X-MEDIA / #EXT-X-KEY 等标签里的 URI="..." 属性
		if strings.HasPrefix(trimmed, "#") {
			if strings.Contains(trimmed, `URI="`) {
				lines[i] = rewriteURIAttr(trimmed, toProxy)
			}
			continue
		}

		// 普通 URL 行（分片地址、子 M3U8 地址）
		lines[i] = toProxy(trimmed)
	}

	return strings.Join(lines, "\n")
}

// rewriteURIAttr 将标签行里 URI="xxx" 替换为代理地址
func rewriteURIAttr(line string, toProxy func(string) string) string {
	const prefix = `URI="`
	start := strings.Index(line, prefix)
	if start < 0 {
		return line
	}
	start += len(prefix)
	end := strings.Index(line[start:], `"`)
	if end < 0 {
		return line
	}
	end += start
	original := line[start:end]
	return line[:start] + toProxy(original) + line[end:]
}

// forwardResponse 转发普通响应（TS 分片等）
func (p *RefererProxy) forwardResponse(resp *http.Response, w http.ResponseWriter) error {
	// 复制响应头
	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	// 添加 CORS 头
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// 设置状态码
	w.WriteHeader(resp.StatusCode)

	// 复制响应体
	_, err := io.Copy(w, resp.Body)
	return err
}
