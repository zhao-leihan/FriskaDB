import React, { useState } from 'react';
import './App.css';
import TitleBar from './components/TitleBar';
import Sidebar from './components/Sidebar';
import MainContent from './components/MainContent';
import PaintSplash from './components/PaintSplash';
import LoginScreen from './components/LoginScreen';

function App() {
    const [showSplash, setShowSplash] = useState(true);
    const [connected, setConnected] = useState(false);
    const [connectionConfig, setConnectionConfig] = useState(null);
    const [tables, setTables] = useState([]);
    const [selectedTable, setSelectedTable] = useState(null);
    const [showConnectionDialog, setShowConnectionDialog] = useState(false);

    // Define all functions FIRST before using them
    const loadTables = async (config) => {
        try {
            const result = await window.electron.query(config, 'FRISSHOW FRISKABLES;');
            if (result.success && result.data) {
                setTables(result.data);
            }
        } catch (error) {
            console.error('Failed to load tables:', error);
        }
    };

    const executeQuery = async (query) => {
        if (!connectionConfig) {
            return null;
        }

        try {
            const result = await window.electron.query(connectionConfig, query);

            // Auto-refresh sidebar if CREATE/DROP table query
            const trimmedQuery = query.trim().toUpperCase();
            if (result.success && (trimmedQuery.startsWith('FRISRATE') || trimmedQuery.startsWith('FRISDROP'))) {
                setTimeout(() => loadTables(connectionConfig), 500);
            }

            return result;
        } catch (error) {
            throw error;
        }
    };

    const handleConnect = async (config) => {
        try {
            const result = await window.electron.connect(config);
            if (result.success) {
                setConnectionConfig(config);
                setConnected(true);
                setShowConnectionDialog(false);
                // Load tables after connection
                await loadTables(config);
            }
        } catch (error) {
            throw error;
        }
    };

    const handleDisconnect = () => {
        setConnected(false);
        setConnectionConfig(null);
        setTables([]);
        setSelectedTable(null);
        setShowConnectionDialog(false);
    };

    // Show splash screen
    if (showSplash) {
        return <PaintSplash onComplete={() => setShowSplash(false)} />;
    }

    // Show login screen if not connected
    if (!connected) {
        return (
            <div className="app">
                <TitleBar connected={false} />
                <LoginScreen onConnect={handleConnect} />
            </div>
        );
    }

    // Main application view
    return (
        <div className="app">
            <TitleBar connected={connected} />
            <div className="app-body">
                <Sidebar
                    tables={tables}
                    selectedTable={selectedTable}
                    onSelectTable={setSelectedTable}
                    onDisconnect={handleDisconnect}
                    onRefresh={() => loadTables(connectionConfig)}
                    connected={connected}
                />
                <MainContent
                    connected={connected}
                    selectedTable={selectedTable}
                    executeQuery={executeQuery}
                />
            </div>
        </div>
    );
}

export default App;
