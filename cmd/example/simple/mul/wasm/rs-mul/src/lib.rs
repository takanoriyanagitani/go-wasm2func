#[allow(unsafe_code)]
#[no_mangle]
pub extern "C" fn mul32i(a: i32, b: i32) -> i32 {
    a * b
}

#[allow(unsafe_code)]
#[no_mangle]
pub extern "C" fn mul64i(a: i64, b: i64) -> i64 {
    a * b
}

#[allow(unsafe_code)]
#[no_mangle]
pub extern "C" fn mul32f(a: f32, b: f32) -> f32 {
    a * b
}

#[allow(unsafe_code)]
#[no_mangle]
pub extern "C" fn mul64f(a: f64, b: f64) -> f64 {
    a * b
}
