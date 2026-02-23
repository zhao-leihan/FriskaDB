import React, { useState } from 'react';
import { Server, Lock, User, AlertCircle, Loader2, UserPlus, LogIn } from 'lucide-react';
import './LoginScreen.css';

const LoginScreen = ({ onConnect }) => {
    const [mode, setMode] = useState('login'); // 'login' or 'register'
    const [host, setHost] = useState('localhost');
    const [port, setPort] = useState('7171');
    const [username, setUsername] = useState('admin');
    const [password, setPassword] = useState('rayhan123');
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState('');
    const [success, setSuccess] = useState('');

    const handleLogin = async (e) => {
        e.preventDefault();
        setLoading(true);
        setError('');
        setSuccess('');

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

    const handleRegister = async (e) => {
        e.preventDefault();
        setLoading(true);
        setError('');
        setSuccess('');

        try {
            const result = await window.electron.register({
                host,
                port: parseInt(port),
                username,
                password
            });

            if (result.success) {
                setSuccess(`User '${username}' registered successfully! You can now login.`);
                setTimeout(() => {
                    setMode('login');
                    setSuccess('');
                }, 2000);
            }
        } catch (err) {
            setError(err.message || 'Registration failed');
        } finally {
            setLoading(false);
        }
    };

    const logoStyle = {
        width: '100px',
        height: '100px',
        marginBottom: '24px',
        backgroundColor: '#991B1B', // Maroon
        WebkitMask: `url(logo.png) no-repeat center / contain`,
        mask: `url(logo.png) no-repeat center / contain`,
    };

    return (
        <div className="login-screen">
            <div className="login-card animate-slide-in">
                <div className="login-header">
                    <div style={logoStyle} className="login-logo-base"></div>
                    <h1>Rayhan Compass</h1>
                    <p>{mode === 'login' ? 'Connect to Server' : 'Create New Account'}</p>
                </div>

                {/* Tab Switcher */}
                <div className="tab-switcher">
                    <button
                        className={`tab-btn ${mode === 'login' ? 'active' : ''}`}
                        onClick={() => { setMode('login'); setError(''); setSuccess(''); }}
                        type="button"
                    >
                        <LogIn size={16} />
                        Login
                    </button>
                    <button
                        className={`tab-btn ${mode === 'register' ? 'active' : ''}`}
                        onClick={() => { setMode('register'); setError(''); setSuccess(''); setUsername(''); setPassword(''); }}
                        type="button"
                    >
                        <UserPlus size={16} />
                        Register
                    </button>
                </div>

                <form onSubmit={mode === 'login' ? handleLogin : handleRegister} className="login-form">
                    <div className="form-group">
                        <label htmlFor="host">
                            <Server size={16} />
                            Host & Port
                        </label>
                        <div className="input-row">
                            <input
                                id="host"
                                type="text"
                                value={host}
                                onChange={(e) => setHost(e.target.value)}
                                placeholder="localhost"
                                required
                                className="input-host"
                            />
                            <span className="colon">:</span>
                            <input
                                id="port"
                                type="number"
                                value={port}
                                onChange={(e) => setPort(e.target.value)}
                                placeholder="7171"
                                required
                                className="input-port"
                            />
                        </div>
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
                            placeholder={mode === 'login' ? 'admin' : 'Choose username'}
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

                    {success && (
                        <div className="success-message">
                            <UserPlus size={16} />
                            {success}
                        </div>
                    )}

                    <button
                        type="submit"
                        className="btn-login"
                        disabled={loading}
                    >
                        {loading ? (
                            <>
                                <Loader2 size={16} className="animate-spin" />
                                {mode === 'login' ? 'Connecting...' : 'Registering...'}
                            </>
                        ) : (
                            mode === 'login' ? 'Connect' : 'Register'
                        )}
                    </button>
                </form>
            </div>

            <div className="login-footer">
                <p>RayhanDB Compass v1.0.0 &bull; Made with 💝 by Rayhan</p>
            </div>
        </div>
    );
};

export default LoginScreen;
