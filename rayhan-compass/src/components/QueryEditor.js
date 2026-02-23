import React, { useState } from 'react';
import { Play, Loader2, CheckCircle, AlertCircle } from 'lucide-react';
import './QueryEditor.css';

const QueryEditor = ({ executeQuery }) => {
    const [query, setQuery] = useState('');
    const [result, setResult] = useState(null);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState('');

    const handleExecute = async () => {
        if (!executeQuery) {
            setError('Query execution not available');
            return;
        }

        if (!query.trim()) {
            return;
        }

        setLoading(true);
        setError('');
        setResult(null);

        try {
            const res = await executeQuery(query);

            if (res && res.success) {
                setResult(res);
            } else {
                setError(res?.error || 'Query failed');
            }
        } catch (err) {
            setError(err.message);
        } finally {
            setLoading(false);
        }
    };

    return (
        <div className="query-editor">
            <div className="editor-section">
                <div className="editor-header">
                    <span>Query Editor</span>
                    <button
                        className="btn-execute"
                        onClick={handleExecute}
                        disabled={loading}
                    >
                        {loading ? (
                            <>
                                <Loader2 size={16} className="animate-spin" />
                                Executing...
                            </>
                        ) : (
                            <>
                                <Play size={16} />
                                Run Query
                            </>
                        )}
                    </button>
                </div>
                <textarea
                    className="query-textarea"
                    value={query}
                    onChange={(e) => setQuery(e.target.value)}
                    placeholder="Enter your Friska query here...&#10;&#10;Example:&#10;RAYLECT * RAYFROM users;"
                    spellCheck={false}
                />
            </div>

            <div className="results-section">
                <div className="results-header">Results</div>
                <div className="results-content">
                    {loading && (
                        <div className="result-state">
                            <Loader2 size={32} className="animate-spin" />
                            <p>Executing query...</p>
                        </div>
                    )}

                    {error && (
                        <div className="result-state error">
                            <AlertCircle size={32} />
                            <p>{error}</p>
                        </div>
                    )}

                    {result && result.success && (
                        <div className="result-success">
                            <div className="result-message">
                                <CheckCircle size={16} />
                                {result.message || 'Query executed successfully'}
                            </div>

                            {result.data && Array.isArray(result.data) && result.data.length > 0 && (
                                <div className="result-table-container">
                                    <table className="result-table">
                                        <thead>
                                            <tr>
                                                {Object.keys(result.data[0]).map((key) => (
                                                    <th key={key}>{key}</th>
                                                ))}
                                            </tr>
                                        </thead>
                                        <tbody>
                                            {result.data.map((row, idx) => (
                                                <tr key={idx}>
                                                    {Object.values(row).map((val, i) => (
                                                        <td key={i}>{String(val)}</td>
                                                    ))}
                                                </tr>
                                            ))}
                                        </tbody>
                                    </table>
                                </div>
                            )}
                        </div>
                    )}

                    {!loading && !error && !result && (
                        <div className="result-state">
                            <Play size={32} strokeWidth={1} />
                            <p>Run a query to see results</p>
                        </div>
                    )}
                </div>
            </div>
        </div>
    );
};

export default QueryEditor;
