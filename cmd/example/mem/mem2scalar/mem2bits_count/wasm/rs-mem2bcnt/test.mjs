import { readFile } from "node:fs/promises";

(() => {

	return Promise.resolve("./rs_mem2bcnt.wasm")
	.then(name => readFile(name))
	.then(bytes => WebAssembly.instantiate(bytes))
	.then(pair => {
		const {
			module,
			instance,
		} = pair || {}
		const {
			memory,

			input_offset,
			count_ones,
		} = instance?.exports || {};

		const iptr = input_offset()
		const iview = new Uint8Array(memory?.buffer, iptr, 4)
		iview[0] = 3 // 2 ones
		iview[1] = 7 // 3 ones
		iview[2] = 7 // 3 ones
		iview[3] = 6 // 2 ones
		const cnt = count_ones(4)

		return cnt
	})
	.then(console.info)
	.catch(console.warn)
	;

})()
