const { app, BrowserWindow, ipcMain } = require('electron');
const { spawn } = require('child_process');
const path = require('path');
const net = require('net');
const fs = require('fs');

const isDev = process.env.NODE_ENV === 'development' || process.env.ELECTRON_START_URL;

let mainWindow;
let serverProcess = null;

// Start RayhanDB server
function startServer() {
  // Determine server path based on environment
  let serverPath;
  let workingDir;

  if (isDev) {
    // Development: relative to project root
    serverPath = path.join(__dirname, '..', '..', 'bin', 'RayhanDB-server.exe');
    workingDir = path.join(__dirname, '..', '..');
  } else {
    // Production: bundled in resources
    serverPath = path.join(process.resourcesPath, 'bin', 'RayhanDB-server.exe');
    workingDir = process.resourcesPath;
  }

  // Check if server executable exists
  if (!fs.existsSync(serverPath)) {
    console.error('[Server] RayhanDB-server.exe not found at:', serverPath);
    return;
  }

  console.log('[Server] Starting RayhanDB server...');

  serverProcess = spawn(serverPath, [], {
    cwd: workingDir,
    detached: false,
    windowsHide: true
  });

  serverProcess.stdout.on('data', (data) => {
    console.log(`[Server] ${data.toString().trim()}`);
  });

  serverProcess.stderr.on('data', (data) => {
    console.error(`[Server Error] ${data.toString().trim()}`);
  });

  serverProcess.on('error', (err) => {
    console.error('[Server] Failed to start:', err);
  });

  serverProcess.on('exit', (code, signal) => {
    if (code !== null) {
      console.log(`[Server] Exited with code ${code}`);
    } else if (signal) {
      console.log(`[Server] Killed with signal ${signal}`);
    }
    serverProcess = null;
  });
}

// Stop RayhanDB server
function stopServer() {
  if (serverProcess) {
    console.log('[Server] Stopping RayhanDB server...');
    serverProcess.kill();
    serverProcess = null;
  }
}

function createWindow() {
  mainWindow = new BrowserWindow({
    width: 1400,
    height: 900,
    minWidth: 1000,
    minHeight: 600,
    backgroundColor: '#FFFFFF',
    webPreferences: {
      nodeIntegration: false,
      contextIsolation: true,
      preload: path.join(__dirname, 'preload.js')
    },
    titleBarStyle: 'hidden',
    frame: false
  });

  const startUrl = process.env.ELECTRON_START_URL || `file://${path.join(__dirname, '../build/index.html')}`;
  mainWindow.loadURL(startUrl);

  // DevTools can be toggled with Ctrl+Shift+I

  // Add Ctrl+Shift+I to toggle DevTools
  mainWindow.webContents.on('before-input-event', (event, input) => {
    if (input.control && input.shift && input.key.toLowerCase() === 'i') {
      if (mainWindow.webContents.isDevToolsOpened()) {
        mainWindow.webContents.closeDevTools();
      } else {
        mainWindow.webContents.openDevTools();
      }
    }
  });

  mainWindow.on('closed', () => {
    mainWindow = null;
  });
}

app.whenReady().then(() => {
  startServer();

  // Wait a bit for server to start before opening window
  setTimeout(() => {
    createWindow();
  }, 1000);
});

app.on('window-all-closed', () => {
  stopServer();
  if (process.platform !== 'darwin') {
    app.quit();
  }
});

app.on('before-quit', () => {
  stopServer();
});

app.on('activate', () => {
  if (BrowserWindow.getAllWindows().length === 0) {
    createWindow();
  }
});

// Window controls
ipcMain.on('minimize-window', () => {
  mainWindow.minimize();
});

ipcMain.on('maximize-window', () => {
  if (mainWindow.isMaximized()) {
    mainWindow.unmaximize();
  } else {
    mainWindow.maximize();
  }
});

ipcMain.on('close-window', () => {
  mainWindow.close();
});

// RayhanDB connection handler
ipcMain.handle('RayhanDB-connect', async (event, config) => {
  const { host, port, username, password } = config;

  return new Promise((resolve, reject) => {
    const client = new net.Socket();

    client.connect(port, host, () => {
      // Send auth request
      const authRequest = {
        id: Date.now().toString(),
        query: 'RAYSHOW RAYTABLES;',
        auth: { username, password }
      };

      client.write(JSON.stringify(authRequest) + '\n');
    });

    client.on('data', (data) => {
      try {
        const response = JSON.parse(data.toString());
        if (response.success) {
          resolve({ success: true, data: response });
        } else {
          reject(new Error(response.error || 'Connection failed'));
        }
      } catch (err) {
        reject(err);
      }
      client.destroy();
    });

    client.on('error', (err) => {
      reject(err);
    });

    client.setTimeout(5000, () => {
      client.destroy();
      reject(new Error('Connection timeout'));
    });
  });
});

// User registration handler
ipcMain.handle('RayhanDB-register', async (event, config) => {
  const { host, port, username, password } = config;

  return new Promise((resolve, reject) => {
    const client = new net.Socket();

    client.connect(port, host, () => {
      // Send registration request as query
      const registerQuery = `FRISREGISTER user:${username} pass:${password}`;
      const registerRequest = {
        id: Date.now().toString(),
        query: registerQuery,
        auth: { username: 'admin', password: 'rayhan123' } // Use admin credentials for registration
      };

      client.write(JSON.stringify(registerRequest) + '\n');
    });

    client.on('data', (data) => {
      try {
        const response = JSON.parse(data.toString());
        if (response.success) {
          resolve({ success: true, message: response.message });
        } else {
          reject(new Error(response.error || 'Registration failed'));
        }
      } catch (err) {
        reject(err);
      }
      client.destroy();
    });

    client.on('error', (err) => {
      reject(err);
    });

    client.setTimeout(5000, () => {
      client.destroy();
      reject(new Error('Registration timeout'));
    });
  });
});

// Query execution
ipcMain.handle('RayhanDB-query', async (event, config, query) => {
  console.log('[ELECTRON] RayhanDB-query CALLED');
  console.log('[ELECTRON] config:', config);
  console.log('[ELECTRON] query:', query);

  const { host, port, username } = config;

  return new Promise((resolve, reject) => {
    const client = new net.Socket();
    console.log(`[ELECTRON] Connecting to ${host}:${port}...`);

    client.connect(port, host, () => {
      const request = {
        id: Date.now().toString(),
        query: query,
        auth: { username, password: config.password }
      };

      client.write(JSON.stringify(request) + '\n');
    });

    client.on('data', (data) => {
      try {
        const response = JSON.parse(data.toString());
        resolve(response);
      } catch (err) {
        reject(err);
      }
      client.destroy();
    });

    client.on('error', (err) => {
      reject(err);
    });

    client.setTimeout(5000, () => {
      client.destroy();
      reject(new Error('Query timeout'));
    });
  });
});
