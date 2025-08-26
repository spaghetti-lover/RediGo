package main

import (
	"errors"
	"fmt"
	"math"
)

// readInt64 đọc một số nguyên 64-bit từ một byte slice theo RESP Integer.
// Dữ liệu phải có dạng ":<số>\r\n".
func readInt64(data []byte) (int64, int, error) {
	// Độ dài tối thiểu: ":0\r\n" = 4 byte
	if len(data) < 4 {
		return 0, 0, errors.New("input data too short")
	}

	// Kiểm tra ký tự đầu tiên
	if data[0] != ':' {
		return 0, 0, errors.New("invalid prefix: must start with ':'")
	}

	var res int64 = 0
	var sign int64 = 1
	pos := 1

	// Xử lý dấu
	if data[pos] == '-' {
		sign = -1
		pos++
	} else if data[pos] == '+' {
		pos++
	}

	// Đọc chữ số
	digitCount := 0
	for pos < len(data) {
		ch := data[pos]
		if ch >= '0' && ch <= '9' {
			// Overflow check
			if res > (math.MaxInt64-int64(ch-'0'))/10 {
				return 0, 0, errors.New("integer overflow")
			}
			res = res*10 + int64(ch-'0')
			pos++
			digitCount++
		} else {
			break
		}
	}

	if digitCount == 0 {
		return 0, 0, errors.New("no digits found")
	}

	// Kiểm tra kết thúc \r\n
	if pos+1 >= len(data) || data[pos] != '\r' || data[pos+1] != '\n' {
		return 0, 0, errors.New("invalid terminator")
	}

	return res * sign, pos + 2, nil
}

func main() {
	tests := [][]byte{
		[]byte(":123\r\n"),
		[]byte(":-456\r\n"),
		[]byte(":+789\r\n"),
		[]byte(":\r\n"),                    // invalid
		[]byte(":9223372036854775808\r\n"), // overflow
	}

	for _, input := range tests {
		num, n, err := readInt64(input)
		if err != nil {
			fmt.Printf("Input %q -> Error: %v\n", input, err)
		} else {
			fmt.Printf("Input %q -> Number: %d, Bytes read: %d\n", input, num, n)
		}
	}
}
