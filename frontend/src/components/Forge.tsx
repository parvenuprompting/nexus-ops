import { useState, useEffect } from 'react';
import { SelectDirectory, StartProcessing } from '../../wailsjs/go/main/App';
import { EventsOn } from '../../wailsjs/runtime/runtime';

// Define explicit type for window.runtime
declare global {
    interface Window {
        runtime: any;
    }
}

export default function Forge() {
    const [dir, setDir] = useState<string>('');
    const [status, setStatus] = useState<'idle' | 'scanning' | 'processing' | 'complete'>('idle');
    const [stats, setStats] = useState({ total: 0, processed: 0 });

    useEffect(() => {
        const unsubError = EventsOn('processing:error', (err: string) => {
            alert('Error: ' + err);
            setStatus('idle');
        });

        const unsubStart = EventsOn('processing:started', (count: number) => {
            setStats({ total: count, processed: 0 });
            setStatus('processing');
        });

        const unsubComplete = EventsOn('processing:complete', () => {
            setStatus('complete');
            setStats(s => ({ ...s, processed: s.total }));
        });

        return () => {
            unsubError();
            unsubStart();
            unsubComplete();
        };
    }, []);

    const handleSelect = async () => {
        const path = await SelectDirectory();
        if (path) setDir(path);
    };

    const handleStart = async () => {
        if (!dir) return;
        setStatus('scanning');
        await StartProcessing(dir);
    };

    return (
        <div className="h-full w-full flex items-center justify-center p-8">
            {/* Main Glass Card */}
            <div className="w-full max-w-lg bg-black/40 backdrop-blur-xl border border-white/10 rounded-lg p-1 shadow-2xl">
                {/* Header */}
                <div className="border-b border-white/5 p-4 bg-white/5 rounded-t-lg">
                    <h2 className="text-center font-bold text-gray-200 tracking-wider">Forge: Image Resizer</h2>
                </div>

                {/* Content */}
                <div className="p-8 flex flex-col gap-6">

                    {/* Directory Selection */}
                    <button
                        onClick={handleSelect}
                        className="w-full py-4 bg-white/5 hover:bg-white/10 border border-white/10 hover:border-white/20 text-gray-300 font-medium rounded transition-all active:scale-[0.98] group"
                    >
                        {dir ? (
                            <span className="text-cyber-primary font-mono text-sm break-all">{dir}</span>
                        ) : (
                            <span className="group-hover:text-white transition-colors">Select Input Folder</span>
                        )}
                    </button>

                    {/* Status Text */}
                    <div className="text-center space-y-1">
                        {!dir && <p className="text-sm text-gray-500 italic">No directory selected</p>}
                    </div>

                    {/* Action Button */}
                    <button
                        onClick={handleStart}
                        disabled={!dir || status === 'scanning' || status === 'processing'}
                        className={`
                            w-full py-3 rounded text-sm font-bold tracking-widest uppercase transition-all duration-300 border
                            ${!dir || status !== 'idle' && status !== 'complete'
                                ? 'bg-black/20 text-gray-700 border-transparent cursor-not-allowed'
                                : 'bg-cyber-primary/10 text-cyber-primary border-cyber-primary/50 hover:bg-cyber-primary hover:text-black shadow-[0_0_20px_rgba(0,242,255,0.1)] hover:shadow-[0_0_20px_rgba(0,242,255,0.4)]'
                            }
                        `}
                    >
                        {status === 'scanning' ? 'Scanning...' : status === 'processing' ? 'Processing...' : 'Start Processing'}
                    </button>

                    {/* Progress Monitor */}
                    <div className="space-y-2 pt-4 border-t border-white/5">
                        <div className="flex justify-between text-xs text-gray-400 font-mono">
                            <span>STATUS</span>
                            <span className={status === 'processing' ? 'text-cyber-primary animate-pulse' : 'text-gray-500'}>
                                {status === 'idle' && 'READY'}
                                {status === 'scanning' && 'SCANNING_FILES...'}
                                {status === 'processing' && `PROCESSING [${stats.total} found]`}
                                {status === 'complete' && 'TASK_COMPLETE'}
                            </span>
                        </div>

                        {/* Progress Bar Container */}
                        <div className="h-6 w-full bg-black/50 rounded border border-white/5 relative overflow-hidden">
                            {/* Bar */}
                            <div
                                className={`h-full bg-gradient-to-r from-cyber-primary/50 to-cyber-primary transition-all duration-300 ${status === 'scanning' && 'animate-pulse w-full'}`}
                                style={{
                                    width: status === 'complete' ? '100%' : status === 'processing' ? '100%' : status === 'scanning' ? '100%' : '0%'
                                }}
                            ></div>

                            {/* Percentage Text Overlay */}
                            <div className="absolute inset-0 flex items-center justify-center text-[10px] font-bold text-white/80 drop-shadow-md">
                                {status === 'complete' ? '100%' : status === 'processing' || status === 'scanning' ? 'WORKING...' : '0%'}
                            </div>
                        </div>

                        <p className="text-center text-[10px] text-gray-600 italic mt-2">
                            Check Radar for live telemetry.
                        </p>
                    </div>

                </div>
            </div>
        </div>
    );
}
