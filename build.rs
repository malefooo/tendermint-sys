use std::env;
use std::process::Command;
// use std::path::Path;

fn main() {
    println!("cargo:rerun-if-changed=tmgo/app.go");
    println!("cargo:rerun-if-changed=tmgo/main.go");
    println!("cargo:rerun-if-changed=build.rs");
    let out_dir = env::var("OUT_DIR").unwrap();
    let code_dir = env::current_dir().unwrap();
    let go_dir = code_dir.join("tmgo");

    Command::new("go")
        .args(&["build", "-buildmode=c-archive", "-o"])
        .arg(&format!("{}/libtmgo.a", out_dir))
        .current_dir(go_dir)
        .status()
        .unwrap();

    println!("cargo:rustc-link-search=native={}", out_dir);
    println!("cargo:rustc-link-lib=static=tmgo");
}
