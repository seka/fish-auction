package testing

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"sort"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// Endpoint はテスト対象のエンドポイント定義
type Endpoint struct {
	Method string
	Path   string
	Weight int
	Body   func() []byte
}

// Metrics はテスト結果のメトリクス
type Metrics struct {
	TotalRequests   int64
	SuccessRequests int64
	FailedRequests  int64
	ResponseTimes   []time.Duration
	Errors          map[string]int64
	mu              sync.Mutex
}

// エンドポイント定義（重み付き）
var endpoints = []Endpoint{
	{Method: "GET", Path: "/api/health", Weight: 50, Body: nil},
	{Method: "GET", Path: "/api/items", Weight: 20, Body: nil},
	{Method: "POST", Path: "/api/items", Weight: 5, Body: generateItemBody},
	{Method: "GET", Path: "/api/fishermen", Weight: 10, Body: nil},
	{Method: "POST", Path: "/api/fishermen", Weight: 5, Body: generateFishermanBody},
	{Method: "GET", Path: "/api/buyers", Weight: 5, Body: nil},
	{Method: "POST", Path: "/api/buyers", Weight: 5, Body: generateBuyerBody},
}

func TestLoadTest(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping load test in short mode")
	}

	// 設定を読み込む
	concurrency := getEnvInt("LOAD_TEST_CONCURRENCY", 1000)
	totalRequests := getEnvInt("LOAD_TEST_REQUESTS", 10000)
	duration := getEnvInt("LOAD_TEST_DURATION", 0) // 0 = 無制限
	targetURL := getEnv("LOAD_TEST_TARGET_URL", "http://localhost:8080")

	t.Logf("Starting load test: concurrency=%d, requests=%d, duration=%ds, target=%s",
		concurrency, totalRequests, duration, targetURL)

	// メトリクスを初期化
	metrics := &Metrics{
		Errors:        make(map[string]int64),
		ResponseTimes: make([]time.Duration, 0, totalRequests),
	}

	// 累積重みを計算
	totalWeight := 0
	for _, ep := range endpoints {
		totalWeight += ep.Weight
	}

	// ワーカーを起動
	var wg sync.WaitGroup
	requestChan := make(chan struct{}, totalRequests)
	stopChan := make(chan struct{})

	startTime := time.Now()

	// ワーカー起動
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go worker(&wg, requestChan, stopChan, targetURL, totalWeight, metrics)
	}

	// 期間指定の場合はタイマーを設定
	if duration > 0 {
		go func() {
			time.Sleep(time.Duration(duration) * time.Second)
			close(stopChan)
		}()
	}

	// リクエストを送信
	sendRequests(duration, totalRequests, requestChan, stopChan)

	wg.Wait()
	elapsed := time.Since(startTime)

	// 結果を表示
	printMetrics(t, metrics, elapsed)
}

// sendRequests はリクエストを送信する
func sendRequests(duration, totalRequests int, requestChan chan<- struct{}, stopChan <-chan struct{}) {
	if duration > 0 {
		// 期間指定の場合は無限にリクエストを送る
		for {
			select {
			case <-stopChan:
				close(requestChan)
				return
			default:
				requestChan <- struct{}{}
			}
		}
	} else {
		// リクエスト数指定の場合
		for i := 0; i < totalRequests; i++ {
			requestChan <- struct{}{}
		}
		close(requestChan)
	}
}

// worker は並行してリクエストを送信するワーカー
func worker(wg *sync.WaitGroup, requestChan <-chan struct{}, stopChan <-chan struct{}, targetURL string, totalWeight int, metrics *Metrics) {
	defer wg.Done()

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	for {
		select {
		case <-stopChan:
			return
		case _, ok := <-requestChan:
			if !ok {
				return
			}

			// ランダムにエンドポイントを選択
			endpoint := selectEndpoint(totalWeight)

			// リクエストを送信
			start := time.Now()
			err := sendRequest(client, targetURL, endpoint)
			duration := time.Since(start)

			// メトリクスを更新
			atomic.AddInt64(&metrics.TotalRequests, 1)
			if err != nil {
				atomic.AddInt64(&metrics.FailedRequests, 1)
				metrics.mu.Lock()
				metrics.Errors[err.Error()]++
				metrics.mu.Unlock()
			} else {
				atomic.AddInt64(&metrics.SuccessRequests, 1)
			}

			metrics.mu.Lock()
			metrics.ResponseTimes = append(metrics.ResponseTimes, duration)
			metrics.mu.Unlock()
		}
	}
}

// selectEndpoint は重みに基づいてランダムにエンドポイントを選択
func selectEndpoint(totalWeight int) Endpoint {
	r := rand.Intn(totalWeight)
	cumulative := 0

	for _, ep := range endpoints {
		cumulative += ep.Weight
		if r < cumulative {
			return ep
		}
	}

	return endpoints[0]
}

// sendRequest はHTTPリクエストを送信
func sendRequest(client *http.Client, baseURL string, endpoint Endpoint) error {
	var body io.Reader
	if endpoint.Body != nil {
		body = bytes.NewReader(endpoint.Body())
	}

	req, err := http.NewRequest(endpoint.Method, baseURL+endpoint.Path, body)
	if err != nil {
		return err
	}

	if endpoint.Body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// レスポンスボディを読み捨て
	_, _ = io.Copy(io.Discard, resp.Body)

	if resp.StatusCode >= 400 {
		return fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	return nil
}

// printMetrics はメトリクスを表示
func printMetrics(t *testing.T, metrics *Metrics, elapsed time.Duration) {
	metrics.mu.Lock()
	defer metrics.mu.Unlock()

	sort.Slice(metrics.ResponseTimes, func(i, j int) bool {
		return metrics.ResponseTimes[i] < metrics.ResponseTimes[j]
	})

	total := metrics.TotalRequests
	success := metrics.SuccessRequests
	failed := metrics.FailedRequests
	successRate := float64(success) / float64(total) * 100
	throughput := float64(total) / elapsed.Seconds()

	t.Logf("\n=== Load Test Results ===")
	t.Logf("Total Requests:     %d", total)
	t.Logf("Successful:         %d (%.2f%%)", success, successRate)
	t.Logf("Failed:             %d (%.2f%%)", failed, 100-successRate)
	t.Logf("Duration:           %s", elapsed.Round(time.Millisecond))
	t.Logf("Throughput:         %.2f req/s", throughput)

	if len(metrics.ResponseTimes) > 0 {
		min := metrics.ResponseTimes[0]
		max := metrics.ResponseTimes[len(metrics.ResponseTimes)-1]
		mean := calculateMean(metrics.ResponseTimes)
		p50 := percentile(metrics.ResponseTimes, 50)
		p95 := percentile(metrics.ResponseTimes, 95)
		p99 := percentile(metrics.ResponseTimes, 99)

		t.Logf("\nResponse Times:")
		t.Logf("  Min:              %s", min.Round(time.Millisecond))
		t.Logf("  Max:              %s", max.Round(time.Millisecond))
		t.Logf("  Mean:             %s", mean.Round(time.Millisecond))
		t.Logf("  P50:              %s", p50.Round(time.Millisecond))
		t.Logf("  P95:              %s", p95.Round(time.Millisecond))
		t.Logf("  P99:              %s", p99.Round(time.Millisecond))
	}

	if len(metrics.Errors) > 0 {
		t.Logf("\nErrors:")
		for errMsg, count := range metrics.Errors {
			t.Logf("  %s: %d", errMsg, count)
		}
	}
}

// calculateMean は平均を計算
func calculateMean(durations []time.Duration) time.Duration {
	if len(durations) == 0 {
		return 0
	}
	var sum time.Duration
	for _, d := range durations {
		sum += d
	}
	return sum / time.Duration(len(durations))
}

// percentile はパーセンタイルを計算
func percentile(durations []time.Duration, p int) time.Duration {
	if len(durations) == 0 {
		return 0
	}
	index := int(float64(len(durations)) * float64(p) / 100.0)
	if index >= len(durations) {
		index = len(durations) - 1
	}
	return durations[index]
}

// generateItemBody はアイテム作成用のリクエストボディを生成
func generateItemBody() []byte {
	body := map[string]interface{}{
		"fisherman_id": rand.Intn(100) + 1,
		"fish_type":    fmt.Sprintf("Fish-%d", rand.Intn(10)),
		"quantity":     rand.Intn(100) + 1,
		"unit":         "kg",
	}
	data, _ := json.Marshal(body)
	return data
}

// generateFishermanBody は漁師作成用のリクエストボディを生成
func generateFishermanBody() []byte {
	body := map[string]interface{}{
		"name":         fmt.Sprintf("Fisherman-%d", rand.Intn(10000)),
		"phone_number": fmt.Sprintf("090-%04d-%04d", rand.Intn(10000), rand.Intn(10000)),
	}
	data, _ := json.Marshal(body)
	return data
}

// generateBuyerBody は買い手作成用のリクエストボディを生成
func generateBuyerBody() []byte {
	body := map[string]interface{}{
		"name":         fmt.Sprintf("Buyer-%d", rand.Intn(10000)),
		"phone_number": fmt.Sprintf("090-%04d-%04d", rand.Intn(10000), rand.Intn(10000)),
	}
	data, _ := json.Marshal(body)
	return data
}

// getEnv は環境変数を取得（デフォルト値あり）
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvInt は環境変数を整数として取得（デフォルト値あり）
func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
