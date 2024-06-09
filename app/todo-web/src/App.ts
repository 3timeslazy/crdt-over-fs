import { crdtOverFs } from "./wasm/init";

import * as automerge from "@automerge/automerge"

const automergeOverFs = {
    emptyState: () => {
      let doc = automerge.init();
      return automerge.save(doc)
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

const syncEngine = crdtOverFs.newSyncS3({
    crdt: automergeOverFs,
    s3: {
        keyId: import.meta.env.VITE_S3_KEY_ID,
        keySecret: import.meta.env.VITE_S3_KEY_SECRET,
        endpoint: import.meta.env.VITE_S3_ENDPOINT,
        region: import.meta.env.VITE_S3_REGION,
        bucket: import.meta.env.VITE_S3_BUCKET
    },
    sync: {
        stateId: 'web.3timeslazy',
        rootDir: '.'
    }
});

export const App = async () => {
    const state = await syncEngine.loadOwnState();
    await syncEngine.saveOwnState(state);

    const doc = automerge.load<AppState>(state);

    let app = document.querySelector<HTMLDivElement>('#app') as HTMLDivElement;

    const appTitle = document.createElement("h2");
    appTitle.innerText = "TODO over FS\n"

    app.appendChild(appTitle);

    for (const task of doc.tasks) {
        let div = document.createElement("div");

        let title = document.createElement("b");
        title.innerText = task.Name;

        let createdBy = document.createElement("p");
        createdBy.innerText = `by ${task.CreatedBy}`;

        div.appendChild(title);
        div.appendChild(createdBy);

        app.appendChild(div);
    }
}

interface AppState {
    tasks: Task[]
}

interface Task {
    Name: string,
    CreatedBy: string,
}

