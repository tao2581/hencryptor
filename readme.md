#### About  Homomorphic Encryption  
https://baike.baidu.com/item/%E5%90%8C%E6%80%81%E5%8A%A0%E5%AF%86/6380351?fr=aladdin  
https://www.zhihu.com/question/27645858  

### Getting Started  
 go mod vendor 
 go run *.go 

#### Cross complie  
require: golang > = 1.13 , docker 
- install fnye-cross 
> go get github.com/fyne-io/fyne-cross 
- build for windows
> fyne-cross windows -env GOPROXY=https://goproxy.cn# hencryptor 

- build for macos
> fyne-cross darwin -env GOPROXY=https://goproxy.cn# hencryptor