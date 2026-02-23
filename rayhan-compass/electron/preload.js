const { contextBridge, ipcRenderer } = require('electron');

contextBridge.exposeInMainWorld('electron', {
    // Window controls
    minimizeWindow: () => ipcRenderer.send('minimize-window'),
    maximizeWindow: () => ipcRenderer.send('maximize-window'),
    closeWindow: () => ipcRenderer.send('close-window'),

    // RayhanDB operations
    connect: (config) => ipcRenderer.invoke('RayhanDB-connect', config),
    register: (config) => ipcRenderer.invoke('RayhanDB-register', config),
    query: (config, query) => ipcRenderer.invoke('RayhanDB-query', config, query)
});
