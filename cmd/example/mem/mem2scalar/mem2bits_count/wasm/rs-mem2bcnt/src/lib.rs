static mut INPUT_BYTES: [u8; 1048576] = [0; 1048576];

#[allow(unsafe_code)]
#[no_mangle]
pub extern "C" fn input_offset() -> *mut u8 {
    let ms: &mut [u8] = unsafe { &mut INPUT_BYTES[..] };
    ms.as_mut_ptr()
}

fn cnt_ones(i: &[u8], sz: usize) -> u64 {
    i.iter()
        .take(sz)
        .map(|u| u.count_ones())
        .map(|u: u32| u64::from(u))
        .sum()
}

#[allow(unsafe_code)]
#[no_mangle]
pub extern "C" fn count_ones(sz: i32) -> u64 {
    let i: &[u8] = unsafe { &INPUT_BYTES[..] };
    cnt_ones(i, sz as usize)
}
