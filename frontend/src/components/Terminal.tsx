import { useState, useRef, useEffect } from 'react';
import { RunCommand } from '../../wailsjs/go/main/App';
import { Terminal as TerminalIcon } from 'lucide-react';

export default function Terminal() {
    const [history, setHistory] = useState<string[]>(['> Nexus Ops Terminal v1.0', '> Type a command and press Enter...']);
    const [input, setInput] = useState('');
    const [isExecuting, setIsExecuting] = useState(false);
    const bottomRef = useRef<HTMLDivElement>(null);

    useEffect(() => {
        bottomRef.current?.scrollIntoView({ behavior: 'smooth' });
    }, [history]);

    const handleKeyDown = async (e: React.KeyboardEvent<HTMLInputElement>) => {
        if (e.key === 'Enter' && input.trim() && !isExecuting) {
            const cmd = input.trim();
            setHistory(prev => [...prev, `$ ${cmd}`]);
            setInput('');
            setIsExecuting(true);

            try {
                const output = await RunCommand(cmd);
                setHistory(prev => [...prev, output]);
            } catch (err: any) {
                setHistory(prev => [...prev, `Error: ${err}`]);
            } finally {
                setIsExecuting(false);
            }
        }
    };

    return (
        <div className="h-full flex flex-col p-6 space-y-4">
            <header className="flex items-center gap-3 border-b border-white/10 pb-4">
                <TerminalIcon className="text-cyber-primary" />
                <h2 className="text-xl font-bold font-mono tracking-widest text-white">TERMINAL</h2>
            </header>

            <div className="flex-1 bg-black/80 border border-white/10 rounded-lg p-4 font-mono text-sm overflow-hidden flex flex-col shadow-inner">
                <div className="flex-1 overflow-y-auto space-y-1 custom-scrollbar pb-2">
                    {history.map((line, i) => (
                        <div key={i} className={`${line.startsWith('$') ? 'text-cyber-primary font-bold mt-2' : 'text-gray-300 whitespace-pre-wrap'}`}>
                            {line}
                        </div>
                    ))}
                    <div ref={bottomRef} />
                </div>

                <div className="flex items-center gap-2 mt-2 border-t border-white/10 pt-2">
                    <span className="text-cyber-primary font-bold">{'>'}</span>
                    <input
                        type="text"
                        value={input}
                        onChange={(e) => setInput(e.target.value)}
                        onKeyDown={handleKeyDown}
                        disabled={isExecuting}
                        className="flex-1 bg-transparent border-none outline-none text-white placeholder-gray-600"
                        placeholder={isExecuting ? "Executing..." : "Enter command..."}
                        autoFocus
                    />
                </div>
            </div>
        </div>
    );
}
