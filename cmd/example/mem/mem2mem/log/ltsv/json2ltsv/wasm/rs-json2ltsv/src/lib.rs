use std::io::Write;

static mut INPUT_JSON: Vec<u8> = vec![];
static mut OUTPUT_LTSV: Vec<u8> = vec![];

#[allow(unsafe_code)]
#[no_mangle]
pub extern "C" fn input_resize(sz: i32) -> i32 {
    let i: &mut Vec<u8> = unsafe { &mut INPUT_JSON };
    i.resize(sz as usize, 0);
    i.capacity().try_into().ok().unwrap_or(-1)
}

#[allow(unsafe_code)]
#[no_mangle]
pub extern "C" fn output_reset(sz: i32) -> i32 {
    let o: &mut Vec<u8> = unsafe { &mut OUTPUT_LTSV };
    let cap: usize = o.capacity();
    let add: usize = (sz as usize).saturating_sub(cap);
    o.try_reserve(add)
        .ok()
        .and_then(|_| {
            o.clear();
            o.capacity().try_into().ok()
        })
        .unwrap_or(-1)
}

#[allow(unsafe_code)]
#[no_mangle]
pub extern "C" fn input_offset() -> *mut u8 {
    let i: &mut Vec<u8> = unsafe { &mut INPUT_JSON };
    i.as_mut_ptr()
}

#[allow(unsafe_code)]
#[no_mangle]
pub extern "C" fn output_offset() -> *mut u8 {
    let o: &mut Vec<u8> = unsafe { &mut OUTPUT_LTSV };
    o.as_mut_ptr()
}

#[derive(serde::Deserialize)]
pub struct Pair {
    pub key: String,
    pub val: String,
}

impl Pair {
    pub fn append2string(&self, target: &mut String, prefix: &str, sep: &str) {
        target.push_str(prefix);
        target.push_str(self.key.as_str());
        target.push_str(sep);
        target.push_str(self.val.as_str());
    }

    pub fn append2string_default(&self, target: &mut String, prefix: &str) {
        self.append2string(target, prefix, ":")
    }

    pub fn append2buf(
        &self,
        target: &mut Vec<u8>,
        prefix: &str,
        sep: &str,
    ) -> Result<(), &'static str> {
        let prefix: &[u8] = prefix.as_bytes();
        let key: &[u8] = self.key.as_bytes();
        let sep: &[u8] = sep.as_bytes();
        let val: &[u8] = self.val.as_bytes();

        target.write(prefix).map_err(|_| "unable to write")?;
        target.write(key).map_err(|_| "unable to write")?;
        target.write(sep).map_err(|_| "unable to write")?;
        target.write(val).map_err(|_| "unable to write")?;
        Ok(())
    }

    pub fn append2buf_default(
        &self,
        target: &mut Vec<u8>,
        prefix: &str,
    ) -> Result<(), &'static str> {
        self.append2buf(target, prefix, ":")
    }
}

#[derive(serde::Deserialize)]
pub struct LogRecord {
    pub pairs: Vec<Pair>,
}

impl LogRecord {
    pub fn write2string(&self, target: &mut String) -> u16 {
        target.clear();
        let mut i = self.pairs.iter();
        let first: Option<_> = i.next();
        let init: u16 = first
            .map(|p| p.append2string_default(target, ""))
            .map(|_| 1)
            .unwrap_or_default();
        i.fold(init, |state, next| {
            next.append2string_default(target, "	");
            state + 1
        })
    }

    pub fn write2buf(&self, target: &mut Vec<u8>) -> Result<u16, &'static str> {
        target.clear();
        let mut i = self.pairs.iter();
        let first: Option<_> = i.next();
        let init: u16 = first
            .and_then(|p| p.append2buf_default(target, "").ok())
            .map(|_| 1)
            .unwrap_or_default();
        i.try_fold(init, |state, next| {
            next.append2buf_default(target, "	").map(|_| state + 1)
        })
    }
}

fn slice2record(s: &[u8]) -> Result<LogRecord, &'static str> {
    let v: Vec<Pair> =
        serde_json::from_slice(s).map_err(|_| "unable to parse a json log record")?;
    Ok(LogRecord { pairs: v })
}

fn json_bytes2ltsv_bytes(j: &[u8], l: &mut Vec<u8>) -> Result<u16, &'static str> {
    let parsed: LogRecord = slice2record(j)?;
    parsed.write2buf(l)?;
    l.len().try_into().map_err(|_| "too big log")
}

#[allow(unsafe_code)]
#[no_mangle]
pub extern "C" fn json2ltsv(isz: i32) -> i32 {
    let isz: usize = isz as usize;
    let i: &[u8] = unsafe { &INPUT_JSON };
    let limited: &[u8] = &i[..isz];
    let o: &mut Vec<u8> = unsafe { &mut OUTPUT_LTSV };
    json_bytes2ltsv_bytes(limited, o)
        .ok()
        .map(|u: u16| u.into())
        .unwrap_or(-1)
}
