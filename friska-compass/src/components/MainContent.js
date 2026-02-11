import React, { useState } from 'react';
import * as Tabs from '@radix-ui/react-tabs';
import { Code, Table as TableIcon } from 'lucide-react';
import QueryEditor from './QueryEditor';
import TableViewer from './TableViewer';
import './MainContent.css';

const MainContent = ({ connected, selectedTable, executeQuery }) => {
    const [activeTab, setActiveTab] = useState('query');

    if (!connected) {
        return (
            <div className="main-content">
                <div className="welcome-screen">
                    <div className="welcome-icon">
                        <Code size={64} strokeWidth={1} />
                    </div>
                    <h1>Welcome to Friska Compass</h1>
                    <p>Connect to your FriskaDB server to get started</p>
                </div>
            </div>
        );
    }

    return (
        <div className="main-content">
            <Tabs.Root value={activeTab} onValueChange={setActiveTab} className="tabs-root">
                <Tabs.List className="tabs-list">
                    <Tabs.Trigger value="query" className="tabs-trigger">
                        <Code size={16} />
                        Query Editor
                    </Tabs.Trigger>
                    <Tabs.Trigger value="table" className="tabs-trigger">
                        <TableIcon size={16} />
                        Browse Data
                    </Tabs.Trigger>
                </Tabs.List>

                <Tabs.Content value="query" className="tabs-content">
                    <QueryEditor executeQuery={executeQuery} />
                </Tabs.Content>

                <Tabs.Content value="table" className="tabs-content">
                    <TableViewer
                        selectedTable={selectedTable}
                        executeQuery={executeQuery}
                    />
                </Tabs.Content>
            </Tabs.Root>
        </div>
    );
};

export default MainContent;
