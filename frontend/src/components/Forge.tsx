import { useState, useEffect } from 'react';
import { SelectDirectory, StartProcessing } from '../../wailsjs/go/main/App';
import { EventsOn } from '../../wailsjs/runtime/runtime';
import { AppError } from '../types';

// Define explicit type for window.runtime
declare global {
    interface Window {
        runtime: any;
    }
}

type ProcessingStatus = 'idle' | 'scanning' | 'processing' | 'complete';

export default function Forge() {
    const [selectedDir, setSelectedDir] = useState<string>('');
    const [status, setStatus] = useState<ProcessingStatus>('idle');
    const [stats, setStats] = useState({ total: 0, processed: 0 });
    const [settings, setSettings] = useState({ format: 'jpg', aspectRatio: 'original' });

    useEffect(() => {
        const unsubError = EventsOn('processing:error', (err: AppError) => {
            alert(`[${err.code}] ${err.message}`);
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
        if (path) setSelectedDir(path);
    };

    const handleStart = async () => {
        if (!selectedDir) return;
        setStatus('scanning'); // Keep scanning status for initial file count
        await StartProcessing(selectedDir, settings);
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
                        {selectedDir ? (
                            <span className="text-cyber-primary font-mono text-sm break-all">{selectedDir}</span>
                        ) : (
                            <span className="group-hover:text-white transition-colors">Select Input Folder</span>
                        )}
                    </button>

                    {/* Active Task (Simple) */}
                    <div className="text-center opacity-50 mb-8 font-mono text-sm min-h-[20px]">
                        {selectedDir || <i>No directory selected</i>}
                    </div>

                    {/* Settings UI */}
                    <div className="mb-8 space-y-4">
                        {/* Format Selection */}
                        <div className="flex justify-center gap-4">
                            {['jpg', 'png'].map(fmt => (
                                <button
                                    key={fmt}
                                    onClick={() => setSettings(s => ({ ...s, format: fmt }))}
                                    className={`px-4 py-2 rounded font-bold uppercase tracking-wider text-xs border transition-all ${settings.format === fmt
                                            ? 'bg-cyber-primary text-black border-cyber-primary shadow-[0_0_10px_rgba(0,242,255,0.3)]'
                                            : 'bg-transparent text-gray-500 border-white/10 hover:border-white/30 hover:text-gray-300'
                                        }`}
                                >
                                    .{fmt}
                                </button>
                            ))}
                        </div>

                        {/* Aspect Ratio Grid */}
                        <div className="grid grid-cols-3 gap-2 max-w-xs mx-auto">
                            {['original', '16:9', '9:16', '1:1', '3:2', '2:3'].map(ratio => (
                                <button
                                    key={ratio}
                                    onClick={() => setSettings(s => ({ ...s, aspectRatio: ratio }))}
                                    className={`px-2 py-2 rounded font-mono text-xs border transition-all ${settings.aspectRatio === ratio
                                            ? 'bg-white/10 text-white border-white/40 shadow-inner'
                                            : 'bg-transparent text-gray-600 border-white/5 hover:bg-white/5 hover:text-gray-400'
                                        }`}
                                >
                                    {ratio === 'original' ? 'ORIG' : ratio}
                                </button>
                            ))}
                        </div>
                    </div>

                    {/* Action Button */}
                    <button
                        onClick={handleStart}
                        disabled={!selectedDir || status === 'scanning' || status === 'processing'}
                        className={`
                            w-full py-3 rounded text-sm font-bold tracking-widest uppercase transition-all duration-300 border
                            ${!selectedDir || status !== 'idle' && status !== 'complete'
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
