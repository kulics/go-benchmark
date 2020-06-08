package benchmark

import (
	"encoding/json"
	"fmt"
	"testing"
)

// 测试
func TestBenchmark(t *testing.T) {
	count := RunBenchmark(10, 3, 1000, func(x int) error {
		fmt.Println(x)
		return nil
	})
	bts, err := json.Marshal(count)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(bts))
}
