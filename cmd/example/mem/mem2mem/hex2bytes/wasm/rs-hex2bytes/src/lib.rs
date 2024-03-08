use std::sync::RwLock;

static INPUT: RwLock<Option<Vec<u8>>> = RwLock::new(None);
static OUTPUT: RwLock<Option<Vec<u8>>> = RwLock::new(None);

static EMSG: RwLock<[u8; 256]> = RwLock::new([0; 256]);
static EMSG_SZ: RwLock<u8> = RwLock::new(0);

static mut EMSG_SIZE: u8 = 0;

fn emsg(msg: &str) -> Result<usize, &'static str> {
    let mb: &[u8] = msg.as_bytes();
    let isz: usize = mb.len().min(256);
    let mt: &[u8] = &mb[..isz];
    let nsz: u8 = mt.len().try_into().map_err(|_| "too long msg")?;

    let mut guard = EMSG.try_write().map_err(|_| "unable to write emsg")?;
    let ma: &mut _ = &mut guard;
    let ms: &mut [u8] = &mut ma[..isz];
    ms.copy_from_slice(mt);

    let mut sguard = EMSG_SZ.try_write().map_err(|_| "unable to write lock")?;
    let mu: &mut u8 = &mut sguard;
    *mu = nsz;

    #[allow(unsafe_code)]
    unsafe {
        EMSG_SIZE = 245;
    };

    Ok(ms.len())
}

#[allow(unsafe_code)]
#[no_mangle]
pub extern "C" fn emsg_size() -> i32 {
    unsafe { EMSG_SIZE.into() }
}

#[allow(unsafe_code)]
#[no_mangle]
pub extern "C" fn emsg_ptr() -> *const u8 {
    EMSG.try_read()
        .ok()
        .map(|g| {
            let a = &g;
            let s: &[u8] = &a[..];
            s.as_ptr()
        })
        .unwrap_or_else(std::ptr::null)
}

#[allow(unsafe_code)]
#[no_mangle]
pub extern "C" fn emsg_sz() -> i32 {
    EMSG_SZ
        .try_read()
        .ok()
        .map(|g| {
            let u: &u8 = &g;
            *u
        })
        .map(|u| u.into())
        .unwrap_or(-1)
}

fn write_lock<T, U, F, I>(l: &RwLock<Option<T>>, f: F, init: I) -> Result<U, &'static str>
where
    F: Fn(&mut T) -> Result<U, &'static str>,
    I: Fn() -> T,
{
    let mut guard = l.try_write().map_err(|_| "unable to write lock")?;
    let mo: &mut Option<T> = &mut guard;
    match mo {
        None => {
            let mut t: T = init();
            let u: U = f(&mut t)?;
            mo.replace(t);
            Ok(u)
        }
        Some(mt) => f(mt),
    }
}

fn read_lock<T, U, F>(l: &RwLock<Option<T>>, f: F) -> Result<U, &'static str>
where
    F: Fn(&T) -> Result<U, &'static str>,
{
    let guard = l.try_read().map_err(|_| "unable to read lock")?;
    let o: &Option<T> = &guard;
    let t: &T = o.as_ref().ok_or("no value")?;
    f(t)
}

fn resize<T>(l: &RwLock<Option<Vec<T>>>, sz: usize, init: T) -> Result<usize, &'static str>
where
    T: Copy,
{
    write_lock(
        l,
        |mv: &mut Vec<T>| {
            mv.resize(sz, init);
            Ok(mv.capacity())
        },
        || Vec::with_capacity(sz),
    )
}

#[allow(unsafe_code)]
#[no_mangle]
pub extern "C" fn input_resize(sz: i32) -> i32 {
    sz.try_into()
        .ok()
        .and_then(|u: usize| resize(&INPUT, u, 0).ok())
        .and_then(|u: usize| u.try_into().ok())
        .unwrap_or(-1)
}

#[allow(unsafe_code)]
#[no_mangle]
pub extern "C" fn output_resize(sz: i32) -> i32 {
    sz.try_into()
        .ok()
        .and_then(|u: usize| resize(&OUTPUT, u, 0).ok())
        .and_then(|u: usize| u.try_into().ok())
        .unwrap_or(-1)
}

#[allow(unsafe_code)]
#[no_mangle]
pub extern "C" fn output_ptr() -> *const u8 {
    read_lock(&OUTPUT, |v: &Vec<u8>| Ok(v.as_ptr()))
        .ok()
        .unwrap_or_else(std::ptr::null)
}

#[allow(unsafe_code)]
#[no_mangle]
pub extern "C" fn input_ptr() -> *const u8 {
    read_lock(&INPUT, |v: &Vec<u8>| Ok(v.as_ptr()))
        .ok()
        .unwrap_or_else(std::ptr::null)
}

fn hex2v(h: &[u8], out: &mut Vec<u8>) -> Result<i32, &'static str> {
    out.clear();
    let mut chunks = h.chunks_exact(2);
    chunks.try_fold(0, |state, next| {
        let s: &str = std::str::from_utf8(next).map_err(|_| "invalid str")?;
        let u: u8 = u8::from_str_radix(s, 16).map_err(|_| "invalid u8")?;
        out.push(u);
        Ok(state + 1)
    })
}

fn hex2vec() -> Result<i32, &'static str> {
    let iguard = INPUT.try_read().map_err(|_| "unable to lock input")?;
    let oi: &Option<_> = &iguard;
    let iv: &Vec<u8> = oi.as_ref().ok_or("no input")?;

    let mut oguard = OUTPUT.try_write().map_err(|_| "unable to lock output")?;
    let oo: &mut Option<_> = &mut oguard;
    let ov: &mut Vec<u8> = oo.as_mut().ok_or("no output")?;
    hex2v(iv, ov)
}

#[allow(unsafe_code)]
#[no_mangle]
pub extern "C" fn hex2bytes() -> i32 {
    hex2vec()
        .map_err(|e| {
            let sz: usize = e.len();
            emsg(e).ok();
            unsafe { EMSG_SIZE = sz as u8 };
        })
        .ok()
        .unwrap_or(-1)
}
