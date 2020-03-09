### memory allocation

There are a few things you can do to make the raw memory allocation routines ( ), calloc( ) and realloc( ) safer.

- template template func



## C++ STL Container

A container class describes an object that holds other objects. , you can do the same thing with an array, but there’s more. This new type of object, which is typically referred to in C++ as a container (also called a collection in some languages), will expand itself whenever necessary to accommodate everything you place inside it.


## Baisc Knowledge


## Note

- 注意当数组作为函数的参数进行传递时，该数组自动退化为同类型的指针。
- 编译器总是要为函数的每个参数制作临时副本，指针参数中。编译器总是要为函数的每个参数制作临时副本，指针参数 p的副本是 _p ，编译器使 _p = p。如果函数体内的程序修改了_p 的内容，就导致参数的内容，就导致参数 p的内容作相应的修改。