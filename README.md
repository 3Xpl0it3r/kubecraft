## 介绍
&emsp;`kubecraft`最简单的`kubernetes`模拟器,不需要依赖`cri runtime`,不依赖`etcd`,也不依赖已经存在的`k8s`集群;
- 为用户测试调度器和控制器提供了最小环境模拟器, 在任何操作系统上都可以运行
- 同时也为学习k8s代码提供了最小学习环境,用户可以在这个基础上一步步构建出一个功能完备的k8s,也可以魔改成任何可能性


## 项目Layout
```text
➜  kubecraft git:(main) ✗ tree -L  2
.
├── README.md
├── hack
│   ├── add_types.sh // 提供脚本,用来添加资源
├── main.go
├── pkg
│   ├── api    // api相关辅函数
│   ├── apis   // 资源定义pkg/apis/core是k8s核心资源(拷贝k8s), 其他资源组自定义(example)
│   ├── apiserver // apiserver
│   ├── certs  //证书相关，用来创建必要的证书和配置文件
│   ├── components //apiserver相关组建 (认证，鉴权)
│   ├── fieldpath // copy from k8s.io/kubernetes/pkg/util
│   ├── generated // openapi/client 由code-generater自动生成创建
│   ├── install // 注册资源
│   ├── registry // 具体资源存储
│   ├── scheduler // 调度器
│   ├── storage // 底层实际存储
│   └── util //copy from k8s.io/kuberntes/pkg/uitil
```

## 开发流程

#### 添加资源
&emsp;本项目提供了一个简单脚本`add_types.sh`方便用户来添加一个资源, 例如想要在资源组为`example`,版本是`v1`下添加一个`Bar`的资源
```bash
# 执行add_types.sh 来创建资源

➜  kubecraft git:(main) ✗ bash hack/add_types.sh example v1 bar
Begin generate some necessary code file
Generate hack/..//pkg/apis/example/v1/bar.go
Generate hack/..//pkg/apis/example/bar.go

# 执行hack/codegen-update.sh 来自动生成相关的必要代码
➜  kubecraft git:(main) ✗ bash hack/codegen-update.sh
Generating deepcopy code for 4 targets
Generating defaulter code for 4 targets
Generating conversion code for 4 targets
Generating openapi code for 1 targets
--- /dev/null   2024-06-22 16:56:06
+++ /var/folders/7s/606mx9pd39l1fbghgslk29000000gn/T/codegen-update.sh.api_violations.XXXXXX.BYVTVRlaBx 2024-06-22 16:56:06
@@ -0,0 +1,9 @@
+API rule violation: names_match,k8s.io/apimachinery/pkg/apis/meta/v1,APIResourceList,APIResources
+API rule violation: names_match,k8s.io/apimachinery/pkg/apis/meta/v1,Duration,Duration
+API rule violation: names_match,k8s.io/apimachinery/pkg/apis/meta/v1,InternalEvent,Object
+API rule violation: names_match,k8s.io/apimachinery/pkg/apis/meta/v1,InternalEvent,Type
+API rule violation: names_match,k8s.io/apimachinery/pkg/apis/meta/v1,MicroTime,Time
+API rule violation: names_match,k8s.io/apimachinery/pkg/apis/meta/v1,StatusCause,Type
+API rule violation: names_match,k8s.io/apimachinery/pkg/apis/meta/v1,Time,Time
+API rule violation: names_match,k8s.io/apimachinery/pkg/runtime,Unknown,ContentEncoding
+API rule violation: names_match,k8s.io/apimachinery/pkg/runtime,Unknown,ContentType
ERROR:
        API rule check failed for /dev/null: new reported violations
        Please read api/api-rules/README.md
➜  kubecraft git:(main) ✗

```

#### 注册资源类型
&emsp;首先将添加的资源注册到`scheme`里面,`internal version`和`external version`都需要注册下
- 针对`internal version`资源,通过编辑,`vim pkg/apis/example/register.go` 文件来修改
```go
func addKnowTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SchemeGroupVersion,
		new(Foo),
		new(FooList),
        new(Bar),  // 新添加资源类型
        new(BarList), // 新添加资源类型
    )
	return nil
}
```
- 针对`external version`资源,通过编辑,`vim pkg/apis/example/v1/register.go` 文件来修改
```go
// addKnownTypes adds the set of types defined in this package to the supplied scheme.
func addKnowTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SchemeGroupVersion,
		new(Foo),
		new(FooList),
		new(Bar),  // 新添加资源类型
		new(BarList), // 新添加资源类型
	)
	metav1.AddToGroupVersion(scheme, SchemeGroupVersion)
	return nil
}
```

#### 注册资源存储
&emsp;仓库已经提供了一个示例`pkg/registry/example/foo`
```bash
~ cp -rf pkg/registry/example/foo pkg/registry/example/bar # 修改资源名称 将foo改成bar
```
&emsp; 将`bar storage` 注册到`apiGroup`里面, `vim pkg/registry/example/store.go`
```
func NewRESTStorage(optsGetter generic.RESTOptionsGetter) *genericapiserver.APIGroupInfo {
	apiGroupInfo := genericapiserver.NewDefaultAPIGroupInfo(exampleapi.GroupName, legacyscheme.Scheme, legacyscheme.ParameterCodec, legacyscheme.Codecs)

	storage := map[string]rest.Storage{}
	if store, err := foostorage.NewStorage(optsGetter); err == nil {
		storage["foos"] = store
	}
    // 将bars storage注册到apiGroupinfo里面
	if store, err := barstorage.NewStorage(optsGetter); err == nil {
		storage["bars"] = store
	}

	apiGroupInfo.VersionedResourcesStorageMap[exampleapiv1.SchemeGroupVersion.Version] = storage

	return &apiGroupInfo
}

```

&emsp; 将`storage`注册到`apiserver`, `vim pkg/apiserver/install.go`
```go 

func RegisterApiGroups(scheme *runtime.Scheme, parameterCodec runtime.ParameterCodec, codec serializer.CodecFactory, optsGetter genericregistry.RESTOptionsGetter) []*genericapiserver.APIGroupInfo {
	apiGroupInfos := []*genericapiserver.APIGroupInfo{}

	// register foo
	/* apiGroupInfos = append(apiGroupInfos, foostorage.NewRESTStorage(optsGetter)) */
	apiGroupInfos = append(apiGroupInfos, corestorage.NewRESTStorage(optsGetter))

    // 注册example storage
	apiGroupInfos = append(apiGroupInfos, examplestorage.NewRESTStorage(optsGetter))

	return apiGroupInfos
}

```

#### 测试
&emsp;
```bash
➜  kubecraft git:(main) ✗ kubectl --kubeconfig pki/kubeconfig api-resources
NAME    SHORTNAMES   APIVERSION   NAMESPACED   KIND
nodes   no           v1           false        Node
bars                 example/v1   true         Bar
foos                 example/v1   true         Foo
➜  kubecraft git:(main) ✗
```

## 拓展自定义Authn
&emsp;参考`pkg/components/authentication/fake.go`


## 拓展自定义Authz
&emsp;参考`pkg/components/authorization/fake.go`


## 开发自定义调度器


## 开发自定义控制器


## Agent(Kubelet)



## TODO
- [ ] use sqlite to persist data
- [ ] authn example
- [ ] authz example
- [ ] scheduler example
- [ ] controller example
