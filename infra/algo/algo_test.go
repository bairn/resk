package algo

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestDoubleAverage(t *testing.T) {
	ForTest("二倍随机算法", t, DoubleAverage)
}
func ForTest(message string, t *testing.T, fn func(count, amount int64) int64) {
	count, amount := int64(10), int64(10000)
	remain := amount
	sum := int64(0)
	for i := int64(0); i < count; i++ {
		x := fn(count-i, remain)
		remain -= x
		sum += x
	}
	Convey(message, t, func() {
		So(sum, ShouldEqual, amount)
	})
}
