import { useState } from 'react';
import { SiphonCheckHTTP } from '../../wailsjs/go/main/App';
import { Globe, Activity, Clock, Play, AlertCircle, CheckCircle } from 'lucide-react';

interface ScanResult {
    url: string;
    statusCode: number;
    latencyMs: number;
    error: string;
    timestamp: string;
}

export default function Siphon() {
    const [url, setUrl] = useState('https://google.com');
    const [results, setResults] = useState<ScanResult[]>([]);
    const [scanning, setScanning] = useState(false);

    const handleScan = async () => {
        if (!url || scanning) return;

        setScanning(true);
        try {
            // Ensure URL has protocol
            const target = url.startsWith('http') ? url : `https://${url}`;
            const res = await SiphonCheckHTTP(target);

            setResults(prev => [{
                url: target,
                statusCode: res.statusCode,
                latencyMs: res.latencyMs,
                error: res.error,
                timestamp: new Date().toLocaleTimeString()
            }, ...prev]);
        } catch (e) {
            console.error(e);
        } finally {
            setScanning(false);
        }
    };

    return (
        <div className="h-full flex flex-col p-6 space-y-6">
            <header className="flex items-center gap-3 border-b border-white/10 pb-4">
                <div className="p-2 bg-blue-500/10 rounded-lg border border-blue-500/20 text-blue-500">
                    <Globe size={24} />
                </div>
                <div>
                    <h2 className="text-xl font-bold font-mono tracking-widest text-white">SIPHON</h2>
                    <p className="text-xs text-gray-500 font-mono">NETWORK SCANNER // LATENCY CHECK</p>
                </div>
            </header>

            {/* Input Area */}
            <div className="bg-white/5 p-4 rounded-lg flex gap-4 items-center">
                <input
                    type="text"
                    value={url}
                    onChange={e => setUrl(e.target.value)}
                    onKeyDown={e => e.key === 'Enter' && handleScan()}
                    className="flex-1 bg-black/40 border border-white/10 rounded px-4 py-3 text-white font-mono outline-none focus:border-cyber-primary"
                    placeholder="Enter target URL (e.g. google.com)"
                />
                <button
                    onClick={handleScan}
                    disabled={scanning}
                    className="bg-cyber-primary text-black px-6 py-3 rounded font-bold uppercase tracking-wider hover:bg-cyan-400 transition-colors disabled:opacity-50 flex items-center gap-2"
                >
                    {scanning ? <Activity className="animate-spin" size={18} /> : <Play size={18} />}
                    SCAN
                </button>
            </div>

            {/* Results Grid */}
            <div className="flex-1 overflow-y-auto space-y-2 pr-2 custom-scrollbar">
                {results.map((res, i) => (
                    <div key={i} className="bg-black/40 border border-white/5 rounded p-4 flex items-center justify-between group hover:border-white/10 animate-in fade-in slide-in-from-top-2">
                        <div className="flex items-center gap-4">
                            {res.error || res.statusCode >= 400 ? (
                                <AlertCircle className="text-red-500" size={20} />
                            ) : (
                                <CheckCircle className="text-green-500" size={20} />
                            )}
                            <div>
                                <div className="font-bold text-gray-200">{res.url}</div>
                                <div className="text-[10px] text-gray-600 font-mono">{res.timestamp}</div>
                            </div>
                        </div>

                        <div className="flex items-center gap-6">
                            <div className="text-right">
                                <div className="text-[10px] text-gray-500 uppercase font-bold">Status</div>
                                <div className={`font-mono font-bold ${res.error || res.statusCode >= 400 ? 'text-red-500' : 'text-green-400'}`}>
                                    {res.error ? 'ERR' : res.statusCode}
                                </div>
                            </div>
                            <div className="text-right w-20">
                                <div className="text-[10px] text-gray-500 uppercase font-bold">Latency</div>
                                <div className="font-mono text-cyber-primary flex items-center justify-end gap-1">
                                    {res.latencyMs > 0 ? `${res.latencyMs}ms` : '-'}
                                    <Clock size={12} className="opacity-50" />
                                </div>
                            </div>
                        </div>
                    </div>
                ))}

                {results.length === 0 && (
                    <div className="h-full flex flex-col items-center justify-center opacity-30 min-h-[200px]">
                        <Activity size={48} className="mb-2" />
                        <span className="text-xs font-mono">READY_TO_SCAN</span>
                    </div>
                )}
            </div>
        </div>
    );
}
