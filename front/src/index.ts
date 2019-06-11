import { loadWASM } from './wasmloader';

console.log("Hello from Typescript land!")

loadWASM(async () => {
  console.log('WASM loaded!');
});
