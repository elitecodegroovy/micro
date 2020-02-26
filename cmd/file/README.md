#file-sword
File sword project solved the problem that how to sort the number in a big file. The big file has the following limits.

- file size is bigger than the OS memory free cache size.
- make sure the programming is as fast as possible.
- make sure the data in the big file is order by ascend.
- make sure the programming is never be crashed.


Split the big file into small file. How do you decide the piece of each small file size?
It is closed to the memory free cache is not a good choice. Now I will demonstrate the cases.

os RAM free  1635 M, max chunk file size : 545 M
total time elapsed 267.271180s

os RAM free  2713 M, max chunk file size : 904 M
total time elapsed 253.649489s