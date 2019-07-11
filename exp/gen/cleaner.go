package gen

import (
	"fmt"
	"reflect"
	"sort"

	"github.com/containous/traefik/v2/pkg/config/dynamic"
	"github.com/containous/traefik/v2/pkg/config/runtime"
)

func cleanServers(element *dynamic.Configuration) {
	for _, svc := range element.HTTP.Services {
		if svc.LoadBalancer != nil {
			server := svc.LoadBalancer.Servers[0]
			svc.LoadBalancer.Servers = nil
			svc.LoadBalancer.Servers = append(svc.LoadBalancer.Servers, server)
		}
	}

	for _, svc := range element.TCP.Services {
		if svc.LoadBalancer != nil {
			server := svc.LoadBalancer.Servers[0]
			svc.LoadBalancer.Servers = nil
			svc.LoadBalancer.Servers = append(svc.LoadBalancer.Servers, server)
		}
	}

	for _, svc := range element.UDP.Services {
		if svc.LoadBalancer != nil {
			server := svc.LoadBalancer.Servers[0]
			svc.LoadBalancer.Servers = nil
			svc.LoadBalancer.Servers = append(svc.LoadBalancer.Servers, server)
		}
	}
}

func cleanHTTPServices(element *dynamic.Configuration) {
	valSvcs := reflect.ValueOf(element.HTTP.Services)

	svcRoot := element.HTTP.Services["Service0"]
	valueSvcRoot := reflect.ValueOf(svcRoot).Elem()

	var svcFieldNames []string
	for i := 0; i < valueSvcRoot.NumField(); i++ {
		svcFieldNames = append(svcFieldNames, valueSvcRoot.Type().Field(i).Name)
	}

	sort.Strings(svcFieldNames)

	for i, fieldName := range svcFieldNames {
		v := reflect.New(reflect.TypeOf(dynamic.Service{}))
		v.Elem().FieldByName(fieldName).Set(valueSvcRoot.FieldByName(fieldName))

		valSvcs.SetMapIndex(reflect.ValueOf(fmt.Sprintf("Service%.2d", i+1)), v)
	}

	delete(element.HTTP.Services, "Service0")
	delete(element.HTTP.Services, "Service1")
}

func cleanTCPServices(element *dynamic.Configuration) {
	valSvcs := reflect.ValueOf(element.TCP.Services)

	svcRoot := element.TCP.Services["TCPService0"]
	valueSvcRoot := reflect.ValueOf(svcRoot).Elem()

	var svcFieldNames []string
	for i := 0; i < valueSvcRoot.NumField(); i++ {
		svcFieldNames = append(svcFieldNames, valueSvcRoot.Type().Field(i).Name)
	}

	sort.Strings(svcFieldNames)

	for i, fieldName := range svcFieldNames {
		v := reflect.New(reflect.TypeOf(dynamic.TCPService{}))
		v.Elem().FieldByName(fieldName).Set(valueSvcRoot.FieldByName(fieldName))

		valSvcs.SetMapIndex(reflect.ValueOf(fmt.Sprintf("TCPService%.2d", i+1)), v)
	}

	delete(element.TCP.Services, "TCPService0")
	delete(element.TCP.Services, "TCPService1")
}

func cleanUDPServices(element *dynamic.Configuration) {
	valSvcs := reflect.ValueOf(element.UDP.Services)

	svcRoot := element.UDP.Services["UDPService0"]
	valueSvcRoot := reflect.ValueOf(svcRoot).Elem()

	var svcFieldNames []string
	for i := 0; i < valueSvcRoot.NumField(); i++ {
		svcFieldNames = append(svcFieldNames, valueSvcRoot.Type().Field(i).Name)
	}

	sort.Strings(svcFieldNames)

	for i, fieldName := range svcFieldNames {
		v := reflect.New(reflect.TypeOf(dynamic.UDPService{}))
		v.Elem().FieldByName(fieldName).Set(valueSvcRoot.FieldByName(fieldName))

		valSvcs.SetMapIndex(reflect.ValueOf(fmt.Sprintf("UDPService%.2d", i+1)), v)
	}

	delete(element.UDP.Services, "UDPService0")
	delete(element.UDP.Services, "UDPService1")
}

func cleanMiddlewares(element *dynamic.Configuration) {
	valMds := reflect.ValueOf(element.HTTP.Middlewares)

	mdRoot := element.HTTP.Middlewares["Middleware0"]
	valueMdRoot := reflect.ValueOf(mdRoot).Elem()

	var mdFieldNames []string
	for i := 0; i < valueMdRoot.NumField(); i++ {
		mdFieldNames = append(mdFieldNames, valueMdRoot.Type().Field(i).Name)
	}

	sort.Strings(mdFieldNames)

	for i, fieldName := range mdFieldNames {
		v := reflect.New(reflect.TypeOf(dynamic.Middleware{}))
		v.Elem().FieldByName(fieldName).Set(valueMdRoot.FieldByName(fieldName))

		valMds.SetMapIndex(reflect.ValueOf(fmt.Sprintf("Middleware%.2d", i)), v)
	}

	delete(element.HTTP.Middlewares, "Middleware0")
	delete(element.HTTP.Middlewares, "Middleware1")
}

func cleanMiddlewaresRuntime(element *runtime.Configuration) {
	valMds := reflect.ValueOf(element.Middlewares)

	mdRoot := element.Middlewares["MiddlewareInfo0"].Middleware

	// cleanPlugins(mdRoot)

	valueMdRoot := reflect.ValueOf(mdRoot).Elem()

	var mdFieldNames []string
	for i := 0; i < valueMdRoot.NumField(); i++ {
		mdFieldNames = append(mdFieldNames, valueMdRoot.Type().Field(i).Name)
	}

	sort.Strings(mdFieldNames)

	mdi := reflect.ValueOf(element.Middlewares["MiddlewareInfo0"]).Elem()
	var mdiFieldNames []string
	for i := 0; i < mdi.NumField(); i++ {
		name := mdi.Type().Field(i).Name
		if name != "Middleware" {
			mdiFieldNames = append(mdiFieldNames, name)
		}
	}

	for i, fieldName := range mdFieldNames {
		mi := reflect.New(reflect.TypeOf(runtime.MiddlewareInfo{}))
		for _, mdiFN := range mdiFieldNames {
			mi.Elem().FieldByName(mdiFN).Set(mdi.FieldByName(mdiFN))
		}

		v := reflect.New(reflect.TypeOf(dynamic.Middleware{}))
		v.Elem().FieldByName(fieldName).Set(valueMdRoot.FieldByName(fieldName))

		mi.Elem().FieldByName("Middleware").Set(v)

		valMds.SetMapIndex(reflect.ValueOf(fmt.Sprintf("middleware%.2d@docker", i)), mi)
	}

	delete(element.Middlewares, "MiddlewareInfo0")
	delete(element.Middlewares, "MiddlewareInfo1")
}

// func cleanPlugins(md *dynamic.Middleware) {
// 	if md.Plugin != nil {
// 		md.Plugin = map[string]dynamic.PluginConf{
// 			"foobar": map[string]interface{}{
// 				"foo": "bar",
// 			},
// 		}
// 	}
// }
