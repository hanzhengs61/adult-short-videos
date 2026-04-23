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

// rewriteM3U8URLs 重写 M3U8 中的 URL
func (p *RefererProxy) rewriteM3U8URLs(content string, baseURL *url.URL) string {
	lines := strings.Split(content, "\n")

	for i, line := range lines {
		line = strings.TrimSpace(line)

		// 跳过注释和空行
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// 处理 URL 行
		var fullURL string
		if strings.HasPrefix(line, "http://") || strings.HasPrefix(line, "https://") {
			// 绝对 URL
			fullURL = line
		} else {
			// 相对 URL，需要拼接
			fullURL = fmt.Sprintf("%s://%s%s/%s",
				baseURL.Scheme,
				baseURL.Host,
				strings.TrimSuffix(baseURL.Path, baseURL.Path[strings.LastIndex(baseURL.Path, "/"):]),
				line,
			)
		}

		// 重写为代理 URL
		lines[i] = fmt.Sprintf("/api/storage/proxy?url=%s", url.QueryEscape(fullURL))
	}

	return strings.Join(lines, "\n")
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
