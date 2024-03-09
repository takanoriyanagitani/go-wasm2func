static mut I_BUF: [u8; 1048576] = [0; 1048576];
static mut O_BUF: [u8; 1048576] = [0; 1048576];

#[allow(unsafe_code)]
#[no_mangle]
pub extern "C" fn i_ptr() -> *mut u8 {
	unsafe { I_BUF.as_mut_ptr() }
}

#[allow(unsafe_code)]
#[no_mangle]
pub extern "C" fn o_ptr() -> *mut u8 {
	unsafe { O_BUF.as_mut_ptr() }
}

fn hex2bytes(h: &[u8], b: &mut [u8]) -> Result<usize, &'static str>{
	let chunks = h.chunks_exact(2);
	chunks.enumerate().try_fold(0, |state, pair|{
		let (i, next) = pair;
		let s: &str = std::str::from_utf8(next).map_err(|_| "invalid utf8")?;
		let u: u8 = u8::from_str_radix(s, 16).map_err(|_| "invalid u8")?;
		let mu: &mut u8 = b.get_mut(i).ok_or("out of range")?;
		*mu = u;
		Ok(state + 1)
	})
}

#[allow(unsafe_code)]
#[no_mangle]
pub extern "C" fn hex_string2bytes() -> i32 {
	let i: &[u8] = unsafe { &I_BUF };
	let o: &mut [u8] = unsafe { &mut O_BUF };
	hex2bytes(i, o)
	.ok()
	.and_then(|u| u.try_into().ok())
	.unwrap_or(-1)
}
