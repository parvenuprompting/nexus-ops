import { useState, useEffect } from 'react';
import { GetActiveTasks } from '../../wailsjs/go/main/App';
import { EventsOn } from '../../wailsjs/runtime/runtime';

// Type definitions from Go
interface Task {
    id: string;
    type: string;
    startedAt: string; // ISO string
    status: string;
}

export default function Radar() {
    const [tasks, setTasks] = useState<Task[]>([]);

    useEffect(() => {
        // Hybrid Fetching: 1. Initial Load
        GetActiveTasks().then((initial) => {
            // Wails returns null if empty slice sometimes? Safe check.
            setTasks(initial || []);
        });

        // 2. Event Updates
        const unsub = EventsOn('activeTasks:update', (updatedTasks: Task[]) => {
            setTasks(updatedTasks || []);
        });

        // Optional: Local ticker to force re-render for relative times?
        // React handles re-renders on state change. "Duration" calculation needs re-render.
        const interval = setInterval(() => {
            setTasks(prev => [...prev]); // Trigger re-render to update timestamps
        }, 1000);

        return () => {
            unsub();
            clearInterval(interval);
        };
    }, []);

    const getBadgeColor = (type: string) => {
        if (type.includes('CPU')) return 'bg-orange-500/20 text-orange-400 border-orange-500/50';
        if (type.includes('I/O')) return 'bg-green-500/20 text-green-400 border-green-500/50';
        if (type.includes('Image')) return 'bg-purple-500/20 text-purple-400 border-purple-500/50';
        return 'bg-gray-500/20 text-gray-400 border-gray-500/50';
    };

    const formatDuration = (start: string) => {
        const diff = Date.now() - new Date(start).getTime();
        return (diff / 1000).toFixed(1) + 's';
    };

    return (
        <div className="h-full flex flex-col">
            <h2 className="text-xl font-bold mb-4 text-gray-300">Active Signals</h2>

            {tasks.length === 0 ? (
                <div className="flex-1 flex flex-col items-center justify-center text-gray-600 space-y-2">
                    <div className="text-4xl">📡</div>
                    <div className="font-mono text-sm italic">No active tasks. System idle.</div>
                </div>
            ) : (
                <div className="flex-1 overflow-y-auto space-y-2 pr-2">
                    {tasks.map(task => (
                        <div key={task.id} className="bg-black/40 border border-cyber-border p-3 rounded flex items-center justify-between font-mono text-sm hover:bg-white/5 transition-colors">
                            <div className="flex items-center gap-4">
                                <span className={`px-2 py-0.5 rounded text-xs border ${getBadgeColor(task.type)}`}>
                                    {task.type.includes('Image') ? 'IMG' : task.type.substring(0, 3).toUpperCase()}
                                </span>
                                <span className="text-gray-300">{task.id}</span>
                            </div>
                            <span className="text-cyber-primary">{formatDuration(task.startedAt)}</span>
                        </div>
                    ))}
                </div>
            )}
        </div>
    );
}
