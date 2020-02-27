#file-sword
File sword project solved the problem that how to sort the number in a big file. The big file has the following limits.

- file size is bigger than the OS memory free cache size.
- make sure the programming is as fast as possible.
- make sure the data in the big file is order by ascend.
- make sure the programming is never be crashed.

### Generator Data
Split the big file into small file. How do you decide the piece of each small file size?
It is closed to the memory free cache is not a good choice. Now I will demonstrate the cases.

Firstly, I genrated the random number with the following command:

```
#1.5G data whill be generated. you can try again to generate more data.
go run cmd/file/generator/generator.go

```

All data is in the file `bigLongTypeData.txt` which locates in directory `data/`.

### Run The Application

```shell script
#compile the code
go build cmd/file/main/main.go -o sword

#run the application
./sword
```

Optional input parameters.

- filename : original file path. Default is `data/bigLongTypeData.txt`
- maxChunkFileSize : split the big file into the max chunk file. you can set the value. Default is `free RAM/3`.
- sortedFilename: the output file path. Default is `data\sortedBigLongTypeData.txt`
- temp-output-dir: the temp directory. Default is `data\default`.

### Test Case

The two very import condition is the free RAM size and original file size. Now we have the setting as the following.

```
#one machine with 1G memory cache in Tencet Cloud.
4G File size   
1525 M OS Free RAM
```

Run the app:
```
./sword
```

output the result:
```

os RAM free  1311 M, max chunk file size : 437 M
 sort time: 9.294227 
 sort time: 9.582739 
 sort time: 9.401665 
 sort time: 9.483063 
 sort time: 9.457410 
 sort time: 9.581939 
 sort time: 9.398198 
 sort time: 9.399727 
 sort time: 9.271082 
 sort time: 2.607050 
 >>>total lines(include the empty line): 280666447 
 
 output file path : 'data/sortedBigLongTypeData.txt' , writing time: 219.099338 
 --------------------------------------------
 ----------total time elapsed 498.737140s------------
 --------------------------------------------
```

