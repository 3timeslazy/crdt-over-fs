import './wasm/init.ts'

import * as automerge from "@automerge/automerge"

const automergeOverFs = {
  emptyState: () => {
    let doc = automerge.from({});
    return automerge.save(doc).toString()
  },

  merge: (s1: Uint8Array, s2: Uint8Array) => {
    let doc = automerge.load(s1);
    let doc2 = automerge.load(s2);
    const merged = automerge.merge(doc, doc2);

    return {
      state: automerge.save(merged)
    }
  }
};

(<any>window).Automerge = automerge;
(<any>window).automergeOverFs = automergeOverFs;
(<any>window).s3creds = {
  keyId: import.meta.env.VITE_S3_KEY_ID,
  keySecret: import.meta.env.VITE_S3_KEY_SECRET,
  endpoint: import.meta.env.VITE_S3_ENDPOINT,
  region: import.meta.env.VITE_S3_REGION,
  bucket: import.meta.env.VITE_S3_BUCKET
};
(<any>window).sync = {
  stateId: "web.3timeslazy",
  rootDir: "."
}

document.querySelector<HTMLDivElement>('#app')!.innerHTML = `
  <div>
    <p>
    Try in console:
    <br>
    $ const syncEngine = newSyncS3({sync: sync, crdt: automergeOverFs, s3: s3creds});<br>
    $ const state = await syncEngine.loadOwnState();<br>
    </p>
  </div>
`
