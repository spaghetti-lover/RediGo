Ok, mình demo workflow luôn nhé 🚀. Giả sử bạn muốn làm một **Redis mini trong Go**, chỉ có 2 lệnh: `SET` và `GET`.

---

## 🎯 Workflow học với ChatGPT (ví dụ Redis mini)

### **Bước 1: Hỏi concept thay vì hỏi code**

Bạn:

> Redis lưu key-value trong bộ nhớ như thế nào? Tôi muốn viết phiên bản tối giản trong Go.

ChatGPT (gợi ý):

- Redis dùng `dict` (hash table).
- Go có thể dùng `map[string]string`.
- Về sau có thể thay bằng cấu trúc custom để kiểm soát bộ nhớ.

📌 **Bạn học được idea, chưa có code.**

---

### **Bước 2: Tự code version thô**

Bạn thử viết code từ gợi ý:

```go
package main

import "fmt"

type MiniRedis struct {
	store map[string]string
}

func NewMiniRedis() *MiniRedis {
	return &MiniRedis{store: make(map[string]string)}
}

func (r *MiniRedis) Set(key, value string) {
	r.store[key] = value
}

func (r *MiniRedis) Get(key string) (string, bool) {
	val, ok := r.store[key]
	return val, ok
}

func main() {
	r := NewMiniRedis()
	r.Set("foo", "bar")
	if v, ok := r.Get("foo"); ok {
		fmt.Println("foo =", v)
	}
}
```

👉 Đây là code “xấu” (chưa có expire, concurrency, error handling), nhưng bạn **đã tự viết ra**.

---

### **Bước 3: Nhờ ChatGPT review**

Bạn gửi code cho mình:

> Review MiniRedis này. Tôi có thể tối ưu hoặc thêm tính năng gì?

ChatGPT sẽ trả lời:

- Chưa thread-safe → dùng `sync.RWMutex`.
- Chưa có `EXPIRE`.
- Có thể viết test.

📌 Lúc này bạn thấy **gap** giữa version của bạn và Redis thực tế.

---

### **Bước 4: Bạn cải tiến theo feedback**

Bạn tự thử thêm lock:

```go
type MiniRedis struct {
	store map[string]string
	mu    sync.RWMutex
}

func (r *MiniRedis) Set(key, value string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.store[key] = value
}

func (r *MiniRedis) Get(key string) (string, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	val, ok := r.store[key]
	return val, ok
}
```

📌 Giờ bạn vừa học: Redis cần xử lý concurrency → concept `dict` trong Redis thực ra là **thread-safe** nhờ event loop, còn trong Go thì bạn phải tự lock.

---

### **Bước 5: Bạn tự challenge**

- Thêm `DEL(key)`
- Thêm `EXPIRE(key, ttl)`
- Viết test: `go test -v`
- Benchmark: `go test -bench=.`

👉 Dần dần, bạn build thành một “mini Redis” của riêng bạn, nhớ lâu hơn đọc source code gốc.

---

✅ Cách này giúp bạn:

- Hiểu **ý tưởng** → biến thành code bằng tay.
- ChatGPT chỉ **gợi ý và review**.
- Não bạn phải code → sẽ nhớ lâu.

---

Bạn có muốn mình cùng bạn **làm tiếp bước 2 → 3 → 4** (viết test cho MiniRedis và để mình review cho bạn) không?

Tuyệt vời 👌. Vậy ta đi **Bước 2 → 3 → 4** như đã hứa.
Mình sẽ làm theo flow:

---

## 📝 Bước 2: Bạn viết test cho `MiniRedis`

Bạn thử viết file `mini_redis_test.go` thế này:

```go
package main

import "testing"

func TestMiniRedis_SetGet(t *testing.T) {
	r := NewMiniRedis()

	// Test Set và Get
	r.Set("foo", "bar")
	val, ok := r.Get("foo")

	if !ok {
		t.Errorf("expected key foo to exist")
	}
	if val != "bar" {
		t.Errorf("expected value 'bar', got %s", val)
	}

	// Test key không tồn tại
	_, ok = r.Get("baz")
	if ok {
		t.Errorf("did not expect key baz to exist")
	}
}
```

---

## 🎯 Bước 3: Bạn chạy test

```bash
go test -v
```

- Nếu pass → ổn.
- Nếu fail → ta xem bug ở đâu.

---

## 🕵️ Bước 4: ChatGPT review

Sau khi bạn chạy test, bạn có thể gửi lại output cho mình.
Mình sẽ review test case, gợi ý thêm trường hợp mới (ví dụ: test overwrite key, test concurrency).

---

👉 Giờ tới lượt bạn:
Bạn copy test ở trên vào project, chạy `go test -v` và gửi lại log cho mình nhé.

Bạn muốn mình chèn thêm luôn **benchmark** (`go test -bench=.`) để kiểm tra tốc độ `Set`/`Get` không?

Ok ngon 😎. Vậy mình bổ sung **benchmark** cho `MiniRedis` nhé. Bạn tạo thêm trong file `mini_redis_test.go` luôn:

```go
package main

import (
	"strconv"
	"testing"
)

func TestMiniRedis_SetGet(t *testing.T) {
	r := NewMiniRedis()

	// Test Set và Get
	r.Set("foo", "bar")
	val, ok := r.Get("foo")

	if !ok {
		t.Errorf("expected key foo to exist")
	}
	if val != "bar" {
		t.Errorf("expected value 'bar', got %s", val)
	}

	// Test key không tồn tại
	_, ok = r.Get("baz")
	if ok {
		t.Errorf("did not expect key baz to exist")
	}
}

func BenchmarkMiniRedis_Set(b *testing.B) {
	r := NewMiniRedis()
	for i := 0; i < b.N; i++ {
		r.Set("key"+strconv.Itoa(i), "value")
	}
}

func BenchmarkMiniRedis_Get(b *testing.B) {
	r := NewMiniRedis()
	for i := 0; i < b.N; i++ {
		r.Set("key"+strconv.Itoa(i), "value")
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r.Get("key" + strconv.Itoa(i))
	}
}
```

---

### Cách chạy:

- Test logic:

```bash
go test -v
```

- Benchmark:

```bash
go test -bench=.
```

---

👉 Bạn thử chạy cả 2 lệnh này và gửi lại kết quả cho mình.
Sau đó mình sẽ giúp bạn **phân tích kết quả benchmark** (ops/s, tốc độ so với Redis thật, và chỗ nào có thể tối ưu tiếp).

Bạn muốn mình gợi ý thêm **test case về concurrency** (dùng `t.Parallel()` hoặc `sync.WaitGroup`) để giả lập nhiều client cùng đọc/ghi key không?
