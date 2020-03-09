
## In Linux System

Compile and Run Java Program along with External JARs.
```
javac -cp </path/jar1>:<path/jar2>:<path/jar3> MainClass.java

```
If Compiler throws deprecation warning. You can recompile with -Xlint:deprecation argument.
```
javac -Xlint:deprecation -cp </path/jar1>:<path/jar2>:<path/jar3> MainClass.java

```
Finally, run the Java Program:

```
java -cp </path/jar1>:<path/jar2>:<path/jar3>:. MainClass

```
If you want to run Java Process in Background. You can use nohup.
```
nohup java -cp </path/jar1>:<path/jar2>:<path/jar3>:. MainClass &
```