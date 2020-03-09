

## 初始化go mod

```
go modules init ${/path/to/directory}
```

## 格式化 go.mod文件

```
go mod -fmt
```

## 删除不需要的依赖、新增需要的依赖

```
go mod tidy
```

## 添加依赖或修改依赖版本

```
go mod -require=path@version
```

## 生成 vendor 文件夹

```
go mod -vendor
```

其他的自行 go help mod查看。