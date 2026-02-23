import React, { useEffect, useState } from 'react';
import './PaintSplash.css';
import logoFQL from '../assets/logo.png';

const PaintSplash = ({ onComplete }) => {
    const [stage, setStage] = useState(0);

    useEffect(() => {
        // Stage 1: Splash starts (0ms)
        const t1 = setTimeout(() => setStage(1), 100);

        // Stage 2: Fill screen (800ms)
        const t2 = setTimeout(() => setStage(2), 1000);

        // Stage 3: Reveal logo & text (2000ms)
        const t3 = setTimeout(() => setStage(3), 2000);

        // Stage 4: Finish after 10 seconds (10000ms) 
        const t4 = setTimeout(() => {
            setStage(4);
            if (onComplete) setTimeout(onComplete, 500);
        }, 10000);

        return () => {
            clearTimeout(t1);
            clearTimeout(t2);
            clearTimeout(t3);
            clearTimeout(t4);
        };
    }, [onComplete]);

    if (stage === 4) return null;

    return (
        <div className={`splash-container ${stage >= 4 ? 'fade-out' : ''}`}>
            <div className="paint-source">
                <div className={`paint-drop ${stage >= 1 ? 'drop' : ''}`} />
            </div>

            <div className={`paint-spread ${stage >= 1 ? 'spread' : ''}`} />

            <div className={`splash-logo ${stage >= 3 ? 'visible' : ''}`}>
                <img
                    src={logoFQL}
                    alt="Introduction Logo"
                    className="splash-logo-img-original"
                />
                <h1 className="logo-text-splash">Rayhan Compass</h1>
                <p className="loading-text">Loading...</p>
            </div>
        </div>
    );
};

export default PaintSplash;
