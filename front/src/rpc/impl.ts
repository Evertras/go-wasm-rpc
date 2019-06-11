declare var gowasm: any;

export const wasmImpl = (method: any, requestData: any, callback: any) => {
  // Avoid Zalgo, always async!
  setImmediate(() => {
    const result = gowasm[method.name](requestData, callback);

    callback(null, result);
  });
};
