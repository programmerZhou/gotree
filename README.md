1. md5sum的计算参考开源代码
	```url 
	https://gist.github.com/josephspurrier/90e957f1277964f26852
	```

2. 文件列表保存在的output.log文件中

3. 文件列表中打印目录和文件的信息，目录没有计算md5值

4. 使用方法：
	
	* 列出当前目录下的所有文件
	```shell
	go run tree.go ./ 
	```

	* 列出当前目录下的所有文件，并过滤output.log文件
	```shell
	go run tree.go ./ output.log
	```

	* 列出当前目录下的所有文件，并过滤所有的txt文件
	```shell 
	go run tree.go ./  *.txt
	```

	* 列出当前目录下的所有文件，并过滤.git目录和output.log文件
	```shell 
	go run tree.go ./ output.log .git 
	```


