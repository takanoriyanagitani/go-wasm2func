import { readFile } from "node:fs/promises"

(() => {
    return Promise.resolve("rs_json2ltsv.wasm")
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

            input_offset,
            output_offset,

            json2ltsv,
        } = instance?.exports || {}

        const inputLog = JSON.stringify([
            {key: "timestamp", val: new Date().toISOString()},
            {key: "level", val: "INFO"},
            {key: "msg", val: "wasm parsed."},
        ])
        const len = inputLog.length
        const isz = 3*len
        const icap = input_resize(isz)
        const ocap = output_reset(isz)

        const enc = new TextEncoder()

        const ioff = input_offset()
        const iview = new Uint8Array(memory?.buffer, ioff, isz)
        const { read, written } = enc.encodeInto(inputLog, iview)

        const osz = json2ltsv(written)
        const ooff = output_offset()
        const oview = new Uint8Array(memory?.buffer, ooff, osz)
        const dec = new TextDecoder()
        const ltsv = dec.decode(oview)

        return ltsv
    })
    .then(console.info)
    .catch(console.warn)
})()
