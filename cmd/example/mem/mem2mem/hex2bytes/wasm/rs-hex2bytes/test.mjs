import { readFile } from "node:fs/promises"

(() => {

    return Promise.resolve("./rs_hex2bytes.wasm")
    .then(name => readFile(name))
    .then(bytes => WebAssembly.instantiate(bytes))
    .then(pair => {
        const {
            //module,
            instance,
        } = pair || {}
        const {
            memory,
            emsg_size,
            emsg_ptr,
            emsg_sz,
            input_resize,
            output_resize,
            output_ptr,
            input_ptr,
            hex2bytes,
        } = instance?.exports || {}

        console.info(input_resize(2))
        console.info(output_resize(1))
        const iptr = input_ptr()
        console.info({iptr})
        const ibuf = new Uint8Array(memory?.buffer, iptr, 2)
        console.info({ibuf})
        ibuf[0] = 0x34
        ibuf[1] = 0x32

        console.info(hex2bytes())
        console.info(emsg_sz())

        const optr = output_ptr()
        console.info({optr})
        const obuf = new Uint8Array(memory?.buffer, optr, 1)
        console.info({obuf})
        return memory?.buffer
    })
    .then(console.info)
    .catch(console.warn)

})()
