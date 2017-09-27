package framework

import "context"

type AppContext struct {
	cnt context.Context
}
type Resource struct {
	SystemName string
	Resource interface{}
}
type resourceList []Resource

var resList resourceList

func AddResource(res Resource){
	resList = append(resList, res)
}

func GetContext() (AppContext, error) {
	var cont context.Context
	for _, i := range resList {
		cont = context.WithValue(cont, i.SystemName, i.Resource)
	}
	return AppContext{cnt: cont}, nil
}

func (cont AppContext) GetValue(name string) interface{}{
	return cont.cnt.Value(name)
}

