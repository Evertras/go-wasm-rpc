// Cheaty but works for now since this is all just localized here
declare var WebAssembly: any;
declare var Go: any;
declare var gowasm: any;

export function loadWASM(cb: (err?: Error) => void) {
  if (!WebAssembly.instantiateStreaming) { // polyfill
    WebAssembly.instantiateStreaming = async (resp: any, importObject: any) => {
      const source = await (await resp).arrayBuffer();
      return await WebAssembly.instantiate(source, importObject);
    };
  }
  const go = new Go();
  // let mod: any;
  let inst: any;
  WebAssembly.instantiateStreaming(fetch('lib.wasm'), go.importObject).then((result: any) => {
    // mod = result.module;
    inst = result.instance;

    go.run(inst);

    const waitForReady = () => {
      if (!(gowasm && gowasm.ready)) {
        setTimeout(waitForReady, 100);
        return;
      }

      cb();
    };

    setTimeout(waitForReady, 10);
  }).catch((err: Error) => {
    cb(err);
  });
}
