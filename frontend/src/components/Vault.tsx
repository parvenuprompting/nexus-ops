import { useState, useEffect } from 'react';
import { GetSecrets, AddSecret, RemoveSecret } from '../../wailsjs/go/main/App';
import { Lock, Eye, EyeOff, Trash2, Plus, Key, Save } from 'lucide-react';

interface Secret {
    key: string;
    value: string;
}

export default function Vault() {
    const [secrets, setSecrets] = useState<Secret[]>([]);
    const [visible, setVisible] = useState<Record<string, boolean>>({});
    const [newKey, setNewKey] = useState('');
    const [newValue, setNewValue] = useState('');
    const [isAdding, setIsAdding] = useState(false);

    useEffect(() => {
        loadSecrets();
    }, []);

    const loadSecrets = async () => {
        const loaded = await GetSecrets();
        setSecrets(loaded || []);
    };

    const toggleVisibility = (key: string) => {
        setVisible(prev => ({ ...prev, [key]: !prev[key] }));
    };

    const handleAdd = async () => {
        if (!newKey || !newValue) return;

        await AddSecret(newKey, newValue);
        setNewKey('');
        setNewValue('');
        setIsAdding(false);
        await loadSecrets();
    };

    const handleDelete = async (key: string) => {
        if (confirm(`Delete secret '${key}'?`)) {
            await RemoveSecret(key);
            await loadSecrets();
        }
    };

    return (
        <div className="h-full flex flex-col p-6 space-y-6">
            <header className="flex items-center justify-between border-b border-white/10 pb-4">
                <div className="flex items-center gap-3">
                    <div className="p-2 bg-yellow-500/10 rounded-lg border border-yellow-500/20 text-yellow-500">
                        <Lock size={24} />
                    </div>
                    <div>
                        <h2 className="text-xl font-bold font-mono tracking-widest text-white">THE VAULT</h2>
                        <p className="text-xs text-gray-500 font-mono">SECURE KEY-VALUE STORAGE</p>
                    </div>
                </div>
                <button
                    onClick={() => setIsAdding(!isAdding)}
                    className="flex items-center gap-2 bg-cyber-primary/10 hover:bg-cyber-primary/20 text-cyber-primary border border-cyber-primary/50 px-4 py-2 rounded text-sm font-bold transition-all"
                >
                    <Plus size={16} /> ADD SECRET
                </button>
            </header>

            {isAdding && (
                <div className="bg-white/5 border border-white/10 rounded-lg p-4 animate-in fade-in slide-in-from-top-4">
                    <div className="flex gap-4 items-end">
                        <div className="flex-1 space-y-1">
                            <label className="text-[10px] uppercase font-bold text-gray-500">Key Name</label>
                            <div className="flex items-center gap-2 bg-black/40 border border-white/10 rounded px-3 py-2 focus-within:border-cyber-primary/50">
                                <Key size={14} className="text-gray-500" />
                                <input
                                    value={newKey}
                                    onChange={e => setNewKey(e.target.value.toUpperCase())}
                                    className="bg-transparent border-none outline-none text-sm font-mono w-full text-white placeholder-gray-600"
                                    placeholder="API_KEY_EXAMPLE"
                                    autoFocus
                                />
                            </div>
                        </div>
                        <div className="flex-1 space-y-1">
                            <label className="text-[10px] uppercase font-bold text-gray-500">Value</label>
                            <input
                                value={newValue}
                                onChange={e => setNewValue(e.target.value)}
                                className="bg-black/40 border border-white/10 rounded px-3 py-2 text-sm font-mono w-full text-white placeholder-gray-600 focus:border-cyber-primary/50 outline-none"
                                placeholder="Sensitive Value..."
                                type="password"
                            />
                        </div>
                        <button
                            onClick={handleAdd}
                            className="bg-green-500/10 hover:bg-green-500/20 text-green-500 border border-green-500/50 px-4 py-2 rounded h-[38px] flex items-center gap-2 font-bold text-sm transition-all"
                        >
                            <Save size={16} /> SAVE
                        </button>
                    </div>
                </div>
            )}

            <div className="flex-1 bg-black/40 border border-white/5 rounded-lg overflow-hidden flex flex-col">
                {secrets.length === 0 ? (
                    <div className="flex-1 flex flex-col items-center justify-center text-gray-600 opacity-50">
                        <Lock size={48} className="mb-2" />
                        <span className="text-xs font-mono">NO_SECRETS_STORED</span>
                    </div>
                ) : (
                    <div className="overflow-y-auto custom-scrollbar">
                        <table className="w-full text-left border-collapse">
                            <thead className="bg-white/5 text-xs uppercase font-bold text-gray-500 sticky top-0 backdrop-blur-sm">
                                <tr>
                                    <th className="p-4 bg-black/20">Key</th>
                                    <th className="p-4 bg-black/20">Value</th>
                                    <th className="p-4 bg-black/20 text-right">Actions</th>
                                </tr>
                            </thead>
                            <tbody className="divide-y divide-white/5 font-mono text-sm">
                                {secrets.map(secret => (
                                    <tr key={secret.key} className="hover:bg-white/5 transition-colors group">
                                        <td className="p-4 font-bold text-gray-300">{secret.key}</td>
                                        <td className="p-4 text-gray-500">
                                            {visible[secret.key] ? (
                                                <span className="text-cyber-primary">{secret.value}</span>
                                            ) : (
                                                <span className="tracking-widest">••••••••••••••</span>
                                            )}
                                        </td>
                                        <td className="p-4 text-right flex justify-end gap-2 opacity-100 lg:opacity-50 group-hover:opacity-100 transition-opacity">
                                            <button
                                                onClick={() => toggleVisibility(secret.key)}
                                                className="p-1.5 hover:bg-white/10 rounded text-gray-400 hover:text-white"
                                            >
                                                {visible[secret.key] ? <EyeOff size={14} /> : <Eye size={14} />}
                                            </button>
                                            <button
                                                onClick={() => handleDelete(secret.key)}
                                                className="p-1.5 hover:bg-red-500/10 rounded text-gray-400 hover:text-red-400"
                                            >
                                                <Trash2 size={14} />
                                            </button>
                                        </td>
                                    </tr>
                                ))}
                            </tbody>
                        </table>
                    </div>
                )}
            </div>
        </div>
    );
}
