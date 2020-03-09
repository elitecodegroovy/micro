
# Rust

- Variable and Mutability. For example, in cases where you’re using large data structures, mutating an instance in place may be faster than copying and returning newly allocated instances. With smaller data structures, creating new instances and writing in a more functional programming style may be easier to think through, so lower performance might be a worthwhile penalty for gaining that clarity.

- Expression. Expressions do not include ending semicolons. If you add a semicolon to the end of an expression, you turn it into a statement, which will then not return a value.

- The Stack and the Heap. Accessing data in the heap is slower than accessing data on the stack because you have to follow a pointer to get there. When your code calls a function, the values passed into the function (including, potentially, pointers to data on the heap) and the function’s local variables get pushed onto the stack. When the function is over, those values get popped off the stack.

- Shadow copy. concept of copying the pointer, length, and capacity without copying the data probably sounds like making a shallow copy.

- Reference. Just as variables are immutable by default, so are references. We’re not allowed to modify something(immutable ) we have a reference to. Mutable references have one big restriction: you can have only one mutable reference to a particular piece of data in a particular scope. 

- Rust Feature. The concepts of ownership, borrowing, and slices ensure memory safety in Rust programs at compile time. 

- Index String slice. These characters was 2 bytes.You should use ranges to create string slices (let s = &hello[0..4];) with caution, because doing so can crash your program.

-  Most Rust programmers prefer to use the iterator style.

-  Boxes provide only the indirection and heap allocation; 

- Treating Smart Pointers Like Regular References with the Deref Trait.

- Some properties of code are impossible to detect by analyzing the code: the most famous example is the Halting Problem.

- Rc<T> enables multiple owners of the same data; Box<T> and RefCell<T> have single owners.

- Box<T> allows immutable or mutable borrows checked at compile time; Rc<T> allows only immutable borrows checked at compile time; RefCell<T> allows immutable or mutable borrows checked at runtime.

- Because RefCell<T> allows mutable borrows checked at runtime, you can mutate the value inside the RefCell<T> even when the RefCell<T> is immutable.


- Packages: A Cargo feature that lets you build, test, and share crates
Crates: A tree of modules that produces a library or executable
Modules and use: Let you control the organization, scope, and privacy of paths
Paths: A way of naming an item, such as a struct, function, or module

- Specifically when it comes to questions about the difference between &Trait, Box<Trait>, impl Trait, and dyn Trait.