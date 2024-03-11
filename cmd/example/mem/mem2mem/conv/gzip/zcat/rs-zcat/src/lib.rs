use std::io::Read;

use flate2::read::GzDecoder;

static mut INPUT: Vec<u8> = vec![];
static mut OUTPUT: Vec<u8> = vec![];

fn reset(v: &mut Vec<u8>, sz: usize) {
    v.clear();
    let cap: usize = v.capacity();
    let add: usize = sz.saturating_sub(cap);
    v.reserve(add);
}

fn resize(v: &mut Vec<u8>, sz: usize, init: u8) {
    v.resize(sz, init);
}

#[allow(unsafe_code)]
#[no_mangle]
pub extern "C" fn input_resize(sz: i32) -> i32 {
    let mv: &mut Vec<u8> = unsafe { &mut INPUT };
    resize(mv, sz as usize, 0);
    mv.capacity().try_into().ok().unwrap_or(-1)
}

#[allow(unsafe_code)]
#[no_mangle]
pub extern "C" fn output_reset(sz: i32) -> i32 {
    let mv: &mut Vec<u8> = unsafe { &mut OUTPUT };
    reset(mv, sz as usize);
    mv.capacity().try_into().ok().unwrap_or(-1)
}

fn mv2ptr(mv: &mut Vec<u8>) -> *mut u8 {
    mv.as_mut_ptr()
}

#[allow(unsafe_code)]
#[no_mangle]
pub extern "C" fn offset_i() -> *mut u8 {
    let mv: &mut Vec<u8> = unsafe { &mut INPUT };
    mv2ptr(mv)
}

#[allow(unsafe_code)]
#[no_mangle]
pub extern "C" fn offset_o() -> *mut u8 {
    let mv: &mut Vec<u8> = unsafe { &mut OUTPUT };
    mv2ptr(mv)
}

fn gzip2bytes<R>(rdr: R, buf: &mut Vec<u8>) -> Result<usize, &'static str>
where
    R: Read,
{
    buf.clear();
    let mut dec: GzDecoder<_> = GzDecoder::new(rdr);
    dec.read_to_end(buf)
        .map_err(|_| "unable to decode gzip bytes")
}

fn slice2decoded(s: &[u8], buf: &mut Vec<u8>) -> Result<usize, &'static str> {
    gzip2bytes(s, buf)
}

#[allow(unsafe_code)]
#[no_mangle]
pub extern "C" fn gzip_decode() -> i32 {
    let i: &Vec<u8> = unsafe { &INPUT };
    let o: &mut Vec<u8> = unsafe { &mut OUTPUT };
    slice2decoded(i, o)
        .ok()
        .and_then(|u| u.try_into().ok())
        .unwrap_or(-1)
}
