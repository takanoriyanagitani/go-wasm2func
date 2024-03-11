import { readFile } from "node:fs/promises"

(() => {
    return Promise.resolve("./rs_zcat.wasm")
    .then(name => readFile(name))
    .then(bytes => WebAssembly.instantiate(bytes))
    .then(pair => {
        const {
            module,
            instance,
        } = pair || {}
        const {
            memory,

            input_resize,
            output_reset,

            offset_i,
            offset_o,

            gzip_decode,
        } = instance?.exports || {}

        const gzipBytes = Buffer.from(
            "H4sIAAAAAAAEA8tIzckHAFlRj4UEAAAA",
            "base64",
        )
        input_resize(gzipBytes.length)
        const ioff = offset_i()
        const iview = new Uint8Array(memory?.buffer, ioff, gzipBytes.length)
        iview.set(gzipBytes)
        output_reset(gzipBytes.length * 2)
        const ooff = offset_o()
        const osz = gzip_decode()
        const oview = new Uint8Array(memory?.buffer, ooff, osz)
        const dec = new TextDecoder()
        const decoded = dec.decode(oview)
        return decoded
    })
    .then(console.info)
    .catch(console.warn)
})()
