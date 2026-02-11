const { contextBridge, ipcRenderer } = require('electron');

contextBridge.exposeInMainWorld('electron', {
    // Window controls
    minimizeWindow: () => ipcRenderer.send('minimize-window'),
    maximizeWindow: () => ipcRenderer.send('maximize-window'),
    closeWindow: () => ipcRenderer.send('close-window'),

    // FriskaDB operations
    connect: (config) => ipcRenderer.invoke('friskadb-connect', config),
    register: (config) => ipcRenderer.invoke('friskadb-register', config),
    query: (config, query) => ipcRenderer.invoke('friskadb-query', config, query)
});
