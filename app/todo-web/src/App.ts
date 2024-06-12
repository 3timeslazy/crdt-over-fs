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
    let doc = automerge.load<AppState>(state);

    // Header
    let app = document.querySelector<HTMLDivElement>('#app') as HTMLDivElement;

    const appTitle = document.createElement("h2");
    appTitle.innerText = "TODO over FS\n"

    // Tasks list
    let tasksList = renderList(doc);

    // Sync button
    let syncBtn = document.createElement("button");
    syncBtn.textContent = "Sync";
    syncBtn.addEventListener('click', async () => {
        const synced = await syncEngine.sync(automerge.save(doc));
        await syncEngine.saveOwnState(synced.state);
        doc = automerge.load(synced.state);

        const newList = renderList(doc);
        app.replaceChild(newList, tasksList);
        tasksList = newList;
    })

    // Add button
    const addTask = document.createElement("div");

    const addTaskInput = document.createElement("input");
    addTaskInput.setAttribute("id", "addTask");
    addTaskInput.setAttribute("type", "text");

    let submitBtn = document.createElement("button");
    submitBtn.textContent = "Add Task";
    submitBtn.addEventListener('click', async () => {
        const input = document.getElementById("addTask") as HTMLInputElement;
        const task: Task = {
            title: input.value,
            author: "3timeslazy (web)"
        };
       
        doc = automerge.change(doc, d => {
            automerge.insertAt(d.tasks, 0, task);
        })

        const newList = renderList(doc);
        app.replaceChild(newList, tasksList);
        tasksList = newList;

        await syncEngine.saveOwnState(automerge.save(doc));
    });
    addTask.appendChild(addTaskInput);
    addTask.appendChild(submitBtn);

    app.appendChild(appTitle);
    app.appendChild(tasksList)
    app.appendChild(syncBtn);
    app.appendChild(addTask);
}

function renderList(doc: automerge.Doc<AppState>) {
    const tasksList = document.createElement("ul");
    
    for (const task of doc.tasks) {
        // Task item
        let item = document.createElement("li");

        let title = document.createElement("b");
        title.innerText = task.title;

        let createdBy = document.createElement("p");
        createdBy.innerText = `by ${task.author}`;

        // Delete button
        let deleteBtn = document.createElement("button", );
        deleteBtn.innerText = "Delete";
        deleteBtn.addEventListener('click', async () => {
            const idx = findTaskByName(tasksList, task.title);
            doc = automerge.change(doc, d => {
                automerge.deleteAt(d.tasks, idx);
            });
            await syncEngine.saveOwnState(automerge.save(doc));
            item.parentElement?.removeChild(item);
        })

        item.appendChild(title);
        item.appendChild(createdBy);
        item.appendChild(deleteBtn);

        tasksList.appendChild(item);
    }
    
    return tasksList;
}

function findTaskByName(tasksList: HTMLUListElement, task: string) {
    const titles = tasksList.getElementsByTagName("p");
    for (let i = 0; i < titles.length; i++) {
        if (titles.item(i)?.textContent === task) {
            return i
        }
    }

    return -1
}

interface AppState {
    tasks: Task[]
}

interface Task {
    title: string,
    author: string,
}

