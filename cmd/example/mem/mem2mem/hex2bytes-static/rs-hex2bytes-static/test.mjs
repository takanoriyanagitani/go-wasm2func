import { readFile } from "node:fs/promises";

(() => {
	return Promise.resolve("./rs_hex2bytes_static.wasm")
	.then(name => readFile(name))
	.then(bytes => WebAssembly.instantiate(bytes))
	.then(pair => {
		const {
			module,
			instance,
		} = pair

		const {
			memory,

			i_ptr,
			o_ptr,

			hex_string2bytes,
		} = instance?.exports || {}

		const iptr = i_ptr()
		const iview = new Uint8Array(memory?.buffer, iptr, 4)
		iview[0] = 0x33
		iview[1] = 0x34
		iview[2] = 0x33
		iview[3] = 0x32

		hex_string2bytes(4)

		const optr = o_ptr()
		const oview = new Uint8Array(memory?.buffer, optr, 2)
		const dec = new TextDecoder()
		const decoded = dec.decode(oview)

		return { decoded }
	})
	.then(console.info)
	.catch(console.warn)
})()
