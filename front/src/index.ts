import { loadWASM } from './wasmloader';
import { sample } from './proto/proto';
import { wasmImpl } from './rpc/impl';

console.log("Hello from Typescript land!")

loadWASM(async (err) => {
  if (err) {
    return console.error(err);
  }

  console.log('WASM loaded!');
  const wasmService = sample.WasmService.create(wasmImpl, false, false);

  const helloRequest = sample.EchoRequest.create({ text: "Testing echo" });
  const helloResponse = await wasmService.echo(helloRequest);

  console.log('Received response from WASM:', helloResponse.text);
});
