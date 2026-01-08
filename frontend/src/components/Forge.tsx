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
    const [stats, setStats] = useState({ total: 0, processed: 0 }); // Note: processed isn't strictly tracked individually in app.go yet, but we can simulate

    useEffect(() => {
        // Listen for events
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

        // We can listen to metrics update to infer progress if we wanted, 
        // but for now let's keep it simple or implement better progress tracking later.

        return () => {
            // Cleanup if possible (Wails JS runtime doesn't expose easy unsub in v2 without storing the cancel function returned by EventsOn)
            // Actually eventsOn returns a function to cancel.
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
        // Event listeners handle the rest
    };

    return (
        <div className="glass-panel w-full max-w-2xl mx-auto p-8 flex flex-col items-center">
            <h2 className="text-2xl font-bold mb-8">Forge: Image Resizer</h2>

            <div className="w-full space-y-6">
                {/* Directory Selection */}
                <div className="flex flex-col gap-2">
                    <label className="text-sm text-gray-400">Input Directory</label>
                    <div className="flex gap-2">
                        <input
                            type="text"
                            value={dir}
                            readOnly
                            placeholder="No directory selected"
                            className="bg-black/30 border border-cyber-border rounded px-4 py-2 flex-1 text-gray-300 focus:outline-none focus:border-cyber-primary"
                        />
                        <button
                            onClick={handleSelect}
                            className="bg-cyber-card hover:bg-white/10 border border-cyber-border px-4 py-2 rounded transition-colors"
                        >
                            Browse
                        </button>
                    </div>
                </div>

                {/* Progress Bar */}
                <div className="h-4 bg-black/50 rounded-full overflow-hidden border border-cyber-border/50 relative">
                    <div
                        className={`h-full bg-cyber-primary shadow-[0_0_10px_#00f2ff] transition-all duration-500 ${status === 'scanning' ? 'w-full animate-pulse opacity-50' : ''
                            }`}
                        style={{
                            width: status === 'processing' || status === 'complete'
                                ? '100%'
                                : status === 'scanning' ? '100%' : '0%'
                        }}
                    ></div>
                    {/* If we had granular progress, we'd use width % */}
                </div>

                <div className="text-center text-sm font-mono text-cyber-primary">
                    {status === 'idle' && 'Ready'}
                    {status === 'scanning' && 'Scanning...'}
                    {status === 'processing' && `Processing... (${stats.total} found)`}
                    {status === 'complete' && 'Processing Complete! Check output folder.'}
                </div>

                {/* Action Button */}
                <button
                    onClick={handleStart}
                    disabled={!dir || status === 'scanning' || status === 'processing'}
                    className={`w-full py-3 rounded font-bold tracking-wider transition-all duration-300 ${!dir || status !== 'idle' && status !== 'complete'
                            ? 'bg-gray-800 text-gray-500 cursor-not-allowed'
                            : 'bg-cyber-primary/20 text-cyber-primary border border-cyber-primary hover:bg-cyber-primary hover:text-black shadow-[0_0_15px_rgba(0,242,255,0.3)]'
                        }`}
                >
                    INITIATE PROCESS
                </button>
            </div>
        </div>
    );
}
