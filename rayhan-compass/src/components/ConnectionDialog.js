import React, { useState } from 'react';
import * as Dialog from '@radix-ui/react-dialog';
import { Server, Lock, User, AlertCircle, Loader2 } from 'lucide-react';
import './ConnectionDialog.css';

const ConnectionDialog = ({ open, onConnect }) => {
    const [host, setHost] = useState('localhost');
    const [port, setPort] = useState('7171');
    const [username, setUsername] = useState('admin');
    const [password, setPassword] = useState('rayhan123');
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState('');

    const handleSubmit = async (e) => {
        e.preventDefault();
        setLoading(true);
        setError('');

        try {
            await onConnect({
                host,
                port: parseInt(port),
                username,
                password
            });
        } catch (err) {
            setError(err.message || 'Connection failed');
        } finally {
            setLoading(false);
        }
    };

    return (
        <Dialog.Root open={open}>
            <Dialog.Portal>
                <Dialog.Overlay className="dialog-overlay" />
                <Dialog.Content className="dialog-content">
                    <Dialog.Title className="dialog-title">
                        Connect to RayhanDB
                    </Dialog.Title>
                    <Dialog.Description className="dialog-description">
                        Enter your server connection details
                    </Dialog.Description>

                    <form onSubmit={handleSubmit} className="connection-form">
                        <div className="form-group">
                            <label htmlFor="host">
                                <Server size={16} />
                                Host
                            </label>
                            <input
                                id="host"
                                type="text"
                                value={host}
                                onChange={(e) => setHost(e.target.value)}
                                placeholder="localhost"
                                required
                            />
                        </div>

                        <div className="form-group">
                            <label htmlFor="port">
                                <Server size={16} />
                                Port
                            </label>
                            <input
                                id="port"
                                type="number"
                                value={port}
                                onChange={(e) => setPort(e.target.value)}
                                placeholder="7171"
                                required
                            />
                        </div>

                        <div className="form-group">
                            <label htmlFor="username">
                                <User size={16} />
                                Username
                            </label>
                            <input
                                id="username"
                                type="text"
                                value={username}
                                onChange={(e) => setUsername(e.target.value)}
                                placeholder="admin"
                                required
                            />
                        </div>

                        <div className="form-group">
                            <label htmlFor="password">
                                <Lock size={16} />
                                Password
                            </label>
                            <input
                                id="password"
                                type="password"
                                value={password}
                                onChange={(e) => setPassword(e.target.value)}
                                placeholder="Enter password"
                                required
                            />
                        </div>

                        {error && (
                            <div className="error-message">
                                <AlertCircle size={16} />
                                {error}
                            </div>
                        )}

                        <button
                            type="submit"
                            className="btn-primary"
                            disabled={loading}
                        >
                            {loading ? (
                                <>
                                    <Loader2 size={16} className="animate-spin" />
                                    Connecting...
                                </>
                            ) : (
                                'Connect'
                            )}
                        </button>
                    </form>
                </Dialog.Content>
            </Dialog.Portal>
        </Dialog.Root>
    );
};

export default ConnectionDialog;
