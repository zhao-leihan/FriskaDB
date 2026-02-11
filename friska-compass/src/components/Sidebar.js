import React from 'react';
import * as ScrollArea from '@radix-ui/react-scroll-area';
import * as Separator from '@radix-ui/react-separator';
import { Database, Table, RefreshCw, LogOut } from 'lucide-react';
import './Sidebar.css';

const Sidebar = ({ tables, selectedTable, onSelectTable, onDisconnect, onRefresh, connected }) => {
    if (!connected) {
        return (
            <div className="sidebar">
                <div className="sidebar-empty">
                    <Database size={48} strokeWidth={1} />
                    <p>Not Connected</p>
                </div>
            </div>
        );
    }

    return (
        <div className="sidebar">
            <div className="sidebar-header">
                <h3>Tables</h3>
                <div className="sidebar-actions">
                    <button
                        className="icon-btn"
                        onClick={onRefresh}
                        title="Refresh Tables"
                    >
                        <RefreshCw size={16} />
                    </button>
                    <button
                        className="icon-btn danger"
                        onClick={onDisconnect}
                        title="Disconnect"
                    >
                        <LogOut size={16} />
                    </button>
                </div>
            </div>

            <Separator.Root className="separator" />

            <ScrollArea.Root className="scroll-area">
                <ScrollArea.Viewport className="scroll-viewport">
                    <div className="tables-list">
                        {tables.length === 0 ? (
                            <div className="empty-state">
                                <Table size={32} strokeWidth={1} />
                                <p>No tables found</p>
                            </div>
                        ) : (
                            tables.map((tableName) => (
                                <button
                                    key={tableName}
                                    className={`table-item ${selectedTable === tableName ? 'active' : ''}`}
                                    onClick={() => onSelectTable(tableName)}
                                >
                                    <Table size={16} />
                                    <span>{tableName}</span>
                                </button>
                            ))
                        )}
                    </div>
                </ScrollArea.Viewport>
                <ScrollArea.Scrollbar className="scrollbar" orientation="vertical">
                    <ScrollArea.Thumb className="scrollbar-thumb" />
                </ScrollArea.Scrollbar>
            </ScrollArea.Root>
        </div>
    );
};

export default Sidebar;
