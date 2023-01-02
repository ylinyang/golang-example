package rpc_demo_v1

import (
	"log"
	"reflect"
)

func InitClientProxy(service Service) (err error) {
	// 可以做校验确保它传入的时候必须是一个指向结构体的指针
	val := reflect.ValueOf(service).Elem()
	typ := reflect.TypeOf(service).Elem()
	numField := val.NumField()

	for i := 0; i < numField; i++ {
		fieldType := typ.Field(i)
		fieldValue := val.Field(i)

		if !fieldValue.CanSet() {
			//	 可以报错也可以跳掉
			continue
		}
		//	 替换新的方法实现
		fn := reflect.MakeFunc(fieldType.Type,
			func(args []reflect.Value) (results []reflect.Value) {
				//	把调用信息拼凑起来
				arg := args[1].Interface()
				req := &Request{
					ServiceName: service.Name(),
					MethodName:  fieldType.Name,
					Args:        arg,
				}
				numOut := fieldType.Type.NumOut()
				for j := 0; j < numOut; j++ {
					results = append(results, reflect.New(fieldType.Type.Out(j)).Elem())
				}
				log.Println(req)
				return
			})
		fieldValue.Set(fn)
	}
	return nil
}

type Request struct {
	ServiceName string
	MethodName  string
	Args        any
}

type Service interface {
	Name() string
}
