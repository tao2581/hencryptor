#### build  for windows
require: golang > = 1.13 , docker 
install fnye-cross
> go get github.com/fyne-io/fyne-cross
build
> fyne-cross windows -env GOPROXY=https://goproxy.cn# hencryptor
