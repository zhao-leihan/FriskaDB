import React from 'react';
import { Minus, Square, X } from 'lucide-react';
import './TitleBar.css';

const TitleBar = ({ connected }) => {
    const logoStyle = {
        width: '24px',
        height: '24px',
        backgroundColor: '#991B1B',
        WebkitMask: `url(logo.png) no-repeat center / contain`,
        mask: `url(logo.png) no-repeat center / contain`,
    };

    return (
        <div className="titlebar">
            <div className="titlebar-logo">
                <div style={logoStyle} className="titlebar-logo-icon"></div>
                <div className="logo-text">Friska Compass</div>
                {connected && (
                    <div className="connection-indicator">
                        <div className="status-dot animate-pulse" />
                        <span>Connected</span>
                    </div>
                )}
            </div>
            <div className="titlebar-controls">
                <button
                    className="control-btn"
                    onClick={() => window.electron.minimizeWindow()}
                    aria-label="Minimize"
                >
                    <Minus size={16} />
                </button>
                <button
                    className="control-btn"
                    onClick={() => window.electron.maximizeWindow()}
                    aria-label="Maximize"
                >
                    <Square size={14} />
                </button>
                <button
                    className="control-btn close-btn"
                    onClick={() => window.electron.closeWindow()}
                    aria-label="Close"
                >
                    <X size={16} />
                </button>
            </div>
        </div>
    );
};

export default TitleBar;
