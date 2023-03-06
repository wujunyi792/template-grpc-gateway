## 简单示例

并行求平方和（不要嫌弃示例简单，只是模拟并发）

```
package main

import (
    "fmt"
    "log"
	"pinnacle-primary-be/core/errorx"
)

func main() {
    val, err := mr.MapReduce(func(source chan<- int) {
        // generator
        for i := 0; i < 10; i++ {
            source <- i
        }
    }, func(i int, writer mr.Writer[int], cancel func(error)) {
        // mapper
        writer.Write(i * i)
    }, func(pipe <-chan int, writer mr.Writer[int], cancel func(error)) {
        // reducer
        var sum int
        for i := range pipe {
            sum += i
        }
        writer.Write(sum)
    })
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("result:", val)
}
```