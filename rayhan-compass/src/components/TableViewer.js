import React, { useState, useEffect } from 'react';
import { Table, AlertCircle } from 'lucide-react';
import './TableViewer.css';

const TableViewer = ({ selectedTable, executeQuery }) => {
    const [tableData, setTableData] = useState(null);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState('');

    useEffect(() => {
        if (selectedTable) {
            loadTableData();
        }
    }, [selectedTable]);

    const loadTableData = async () => {
        if (!selectedTable) return;

        setLoading(true);
        setError('');

        try {
            const result = await executeQuery(`RAYLECT * RAYFROM ${selectedTable};`);
            if (result.success) {
                setTableData(result);
            } else {
                setError(result.error || 'Failed to load table data');
            }
        } catch (err) {
            setError(err.message);
        } finally {
            setLoading(false);
        }
    };

    if (!selectedTable) {
        return (
            <div className="table-viewer">
                <div className="empty-viewer">
                    <Table size={48} strokeWidth={1} />
                    <p>Select a table from the sidebar to browse data</p>
                </div>
            </div>
        );
    }

    if (loading) {
        return (
            <div className="table-viewer">
                <div className="empty-viewer">
                    <div className="animate-spin">
                        <Table size={48} />
                    </div>
                    <p>Loading table data...</p>
                </div>
            </div>
        );
    }

    if (error) {
        return (
            <div className="table-viewer">
                <div className="empty-viewer error">
                    <AlertCircle size={48} />
                    <p>{error}</p>
                </div>
            </div>
        );
    }

    if (!tableData || !tableData.data || tableData.data.length === 0) {
        return (
            <div className="table-viewer">
                <div className="table-header">
                    <h3>{selectedTable}</h3>
                    <span className="row-count">0 rows</span>
                </div>
                <div className="empty-viewer">
                    <Table size={48} strokeWidth={1} />
                    <p>No data in this table</p>
                </div>
            </div>
        );
    }

    return (
        <div className="table-viewer">
            <div className="table-header">
                <h3>{selectedTable}</h3>
                <span className="row-count">{tableData.data.length} rows</span>
            </div>
            <div className="table-data-container">
                <table className="data-table">
                    <thead>
                        <tr>
                            <th className="row-number">#</th>
                            {Object.keys(tableData.data[0]).map((key) => (
                                <th key={key}>{key}</th>
                            ))}
                        </tr>
                    </thead>
                    <tbody>
                        {tableData.data.map((row, idx) => (
                            <tr key={idx}>
                                <td className="row-number">{idx + 1}</td>
                                {Object.values(row).map((val, i) => (
                                    <td key={i}>{String(val)}</td>
                                ))}
                            </tr>
                        ))}
                    </tbody>
                </table>
            </div>
        </div>
    );
};

export default TableViewer;
