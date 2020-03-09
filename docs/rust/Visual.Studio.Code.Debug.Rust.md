
# Visual Studio Code Debug Rust

## Install Rust

Rust is installed by shell script:
```
curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh
```

Or download the tar file  "https://forge.rust-lang.org/release/platform-support.html" from your platform.


## Download Visual Studio Code 

Download the vsc from the page "https://code.visualstudio.com/".



## Install VSC extentions

- C/C++ by Microsoft, which will add VS Code support for the C/C++ build tools.

And also the Rust Extension Pack, a set of neat extensions which are ideal for Rust development on VS Code. The pack will install the following extensions:

- **Rust (rls)**, by Daniel Griffen. This neat extension adds language support for Rust to Visual Studio Code using the Rust Language Server (RLS).

- **Crates**, by Seray Uzgur. A tool that will help manage crates.io dependencies within the Cargo.toml files.
 
- **Better TOML**, by bungcip. A VS Code extension that adds syntax highlighting, syntax validation and other useful enhancements to handle TOML files in a better way.
- **CodeLLDB** ,Install the plugin if you’re on a Linux OS.


## Create the “Hello World” project


```
> cargo new hello-world
```


Open the Debug menu, choose Add Configuration and add the following configuration:

- C++ (Windows), if you’re on a Windows OS

- LLDB: Custom Launch, if you’re on a Linux OS

This should add a new launch.json file to your project, which you’ll need to edit to change the value of the program key in the following way:

```
"program": "${workspaceFolder}/target/debug/hello-world.exe",


```

Open the project and configure the DEBUG option.

Right after that, go to File > Preferences > Settings, browse through User > Features > Debug and enable the option “Allow setting breakpoints in any file”.


To do that, double click to the Cargo.toml file to open it and add a reference to the actix-web package right below the [dependencies] block:
```
[dependencies]
actix-web = "1.0.3"

```

Open the /src/main.rs file and replace the existing sample with the following a minimalistic code for a HttpRequest handler that will return a “HELLO WORLD!” HttpResponse:

```
use actix_web::{web, App, HttpResponse, HttpServer, Responder};
fn index() -> impl Responder {
    let my_string = String::from("Rust Async");

    // first_word works on slices of `String`s
    let _word = first_word(&my_string[..]);

    let my_string_literal = "Rust Async";

    // first_word works on slices of string literals
    let _word = first_word(&my_string_literal[..]);

    // Because string literals *are* string slices already,
    // this works too, without the slice syntax!
    let _word = first_word(my_string_literal);
    println!(" word: {}", _word);

    HttpResponse::Ok().body("HELLO WORLD!")
}

```

Still on the the /src/main.rs file, right below the previous function, add a new main() function with the following code:

```
fn main() {
    HttpServer::new(|| {
        App::new()
            .route("/", web::get().to(index))
    })
    .bind("127.0.0.1:8080")
    .unwrap()
    .run()
    .unwrap();
}
```

## Build and Run

Here we are! Now we can compile the program by issuing the following command in the Visual Studio Code terminal window:
```
> cargo build

```


And then run the program in debug mode by selecting Debug > Start Debugging (or by hitting F5).

As soon as our web app is running, we can launch our favorite browser and visit the http://localhost:8080/ address to see the results.
