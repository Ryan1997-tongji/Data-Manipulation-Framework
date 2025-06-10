// package Data_Manipulation_Framework
// @author: chenzhewei.97
// @create date: 2025/5/28
package main

import (
	"context"
	"fmt"
	"github.com/Ryan1997-tongji/Data-Manipulation-Framework/impl"
	"github.com/Ryan1997-tongji/Data-Manipulation-Framework/service"
)

func main() {
	fmt.Println("Hello World")
	service.DoRefresh(context.Background(), "Test Refresh", "abc@gmail.com", &impl.TestRefresherInput{
		FieldX: "test",
		FieldY: 1,
	}, &impl.TestRefresher{})
}
