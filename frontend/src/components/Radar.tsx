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
        GetActiveTasks().then((initial) => {
            setTasks(initial || []);
        });

        const unsub = EventsOn('activeTasks:update', (updatedTasks: Task[]) => {
            setTasks(updatedTasks || []);
        });

        // 1s ticker for relative time updates
        const interval = setInterval(() => {
            setTasks(prev => [...prev]);
        }, 1000);

        return () => {
            unsub();
            clearInterval(interval);
        };
    }, []);

    const getBadgeStyle = (type: string) => {
        if (type.includes('CPU')) return 'bg-orange-500/10 text-orange-400 border-orange-500/30 shadow-[0_0_10px_rgba(249,115,22,0.1)]';
        if (type.includes('I/O')) return 'bg-green-500/10 text-green-400 border-green-500/30 shadow-[0_0_10px_rgba(34,197,94,0.1)]';
        if (type.includes('Image')) return 'bg-purple-500/10 text-purple-400 border-purple-500/30 shadow-[0_0_10px_rgba(168,85,247,0.1)]';
        return 'bg-gray-500/10 text-gray-400 border-gray-500/30';
    };

    const formatDuration = (start: string) => {
        const diff = Date.now() - new Date(start).getTime();
        return (diff / 1000).toFixed(1) + 's';
    };

    return (
        <div className="h-full flex flex-col p-4">
            <div className="flex items-center justify-between mb-6 border-b border-white/5 pb-2">
                <h2 className="text-xl font-bold text-gray-200 tracking-wider font-mono">
                    RADAR // <span className="text-cyber-primary text-sm">LIVE FEED</span>
                </h2>
                <div className="flex gap-2 text-xs font-mono text-gray-500">
                    <span>ACTIVE_TASKS: <span className="text-white">{tasks.length}</span></span>
                </div>
            </div>

            {tasks.length === 0 ? (
                <div className="flex-1 flex flex-col items-center justify-center text-gray-600 space-y-4">
                    <div className="w-16 h-16 rounded-full border border-gray-700 flex items-center justify-center animate-pulse">
                        <div className="w-12 h-12 rounded-full bg-gray-800/50"></div>
                    </div>
                    <div className="font-mono text-sm italic tracking-widest opacity-50">NO_SIGNAL_FOUND</div>
                </div>
            ) : (
                <div className="flex-1 overflow-y-auto space-y-2 pr-2 custom-scrollbar">
                    {tasks.map(task => (
                        <div key={task.id} className="group bg-black/40 border border-white/5 p-3 rounded flex items-center justify-between font-mono text-sm hover:bg-white/5 hover:border-white/10 transition-all">
                            <div className="flex items-center gap-4">
                                <span className={`px-2 py-1 rounded text-[10px] font-bold tracking-wider border ${getBadgeStyle(task.type)} uppercase`}>
                                    {task.type.includes('Image') ? 'IMG_PROCESS' : task.type.replace('-Bound', '')}
                                </span>
                                <span className="text-gray-300 group-hover:text-white transition-colors">{task.id}</span>
                            </div>
                            <span className="text-cyber-primary font-bold">{formatDuration(task.startedAt)}</span>
                        </div>
                    ))}
                </div>
            )}
        </div>
    );
}
