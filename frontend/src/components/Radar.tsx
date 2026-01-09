import { useState, useEffect } from 'react';
import { GetActiveTasks, GetMetrics } from '../../wailsjs/go/main/App';
import { EventsOn } from '../../wailsjs/runtime/runtime';
import { AreaChart, Area, XAxis, YAxis, Tooltip, ResponsiveContainer, BarChart, Bar } from 'recharts';
import { Activity, Cpu, Server, AlertCircle, CheckCircle2 } from 'lucide-react';
import { MetricsDTO, Task } from '../types';

export default function Radar() {
    const [tasks, setTasks] = useState<Task[]>([]);
    const [metrics, setMetrics] = useState<MetricsDTO>({ activeGoroutines: 0, tasksCompleted: 0, errors: 0 });
    const [history, setHistory] = useState<{ time: string; goroutines: number }[]>([]);

    useEffect(() => {
        // Initial Fetch
        GetActiveTasks().then(initial => setTasks(initial || []));
        GetMetrics().then(initial => setMetrics(initial));

        // Subscriptions
        const unsubTasks = EventsOn('activeTasks:update', (updated: Task[]) => setTasks(updated || []));
        const unsubMetrics = EventsOn('metrics:update', (updated: MetricsDTO) => {
            setMetrics(updated);
            setHistory(prev => {
                const now = new Date();
                const timeStr = `${now.getHours()}:${now.getMinutes()}:${now.getSeconds()}`;
                const newHistory = [...prev, { time: timeStr, goroutines: updated.activeGoroutines }];
                if (newHistory.length > 20) newHistory.shift(); // Keep last 20 points
                return newHistory;
            });
        });

        // Ticker for relative time
        const interval = setInterval(() => setTasks(prev => [...prev]), 1000);

        return () => {
            unsubTasks();
            unsubMetrics();
            clearInterval(interval);
        };
    }, []);

    const formatDuration = (start: string) => {
        const diff = Date.now() - new Date(start).getTime();
        return (diff / 1000).toFixed(1) + 's';
    };

    return (
        <div className="h-full flex flex-col p-6 space-y-6">
            <header className="flex items-center justify-between border-b border-white/10 pb-4">
                <div className="space-y-1">
                    <h2 className="text-2xl font-black font-mono tracking-widest text-white flex items-center gap-3">
                        RADAR <span className="text-cyber-primary text-sm bg-cyber-primary/10 px-2 py-0.5 rounded animate-pulse">LIVE FEED</span>
                    </h2>
                    <p className="text-xs text-gray-500 font-mono">SYSTEM TELEMETRY & TASK ORCHESTRATION</p>
                </div>
                <div className="flex gap-4">
                    <StatCard icon={<Cpu size={16} />} label="GOROUTINES" value={metrics.activeGoroutines} color="text-cyber-primary" />
                    <StatCard icon={<CheckCircle2 size={16} />} label="COMPLETED" value={metrics.tasksCompleted} color="text-green-400" />
                    <StatCard icon={<AlertCircle size={16} />} label="ERRORS" value={metrics.errors} color="text-red-400" />
                </div>
            </header>

            <div className="flex-1 grid grid-cols-1 lg:grid-cols-3 gap-6 overflow-hidden">
                {/* Visualizer Column */}
                <div className="lg:col-span-2 flex flex-col gap-6">
                    {/* Activity Chart */}
                    <div className="flex-1 bg-black/40 border border-white/5 rounded-lg p-4 flex flex-col">
                        <h3 className="text-xs font-bold text-gray-400 mb-4 flex items-center gap-2">
                            <Activity size={14} /> THREAD ACTIVITY
                        </h3>
                        <div className="flex-1 min-h-[200px]">
                            <ResponsiveContainer width="100%" height="100%">
                                <AreaChart data={history}>
                                    <defs>
                                        <linearGradient id="colorGro" x1="0" y1="0" x2="0" y2="1">
                                            <stop offset="5%" stopColor="#00f2ff" stopOpacity={0.3} />
                                            <stop offset="95%" stopColor="#00f2ff" stopOpacity={0} />
                                        </linearGradient>
                                    </defs>
                                    <XAxis dataKey="time" stroke="#444" fontSize={10} tickLine={false} />
                                    <YAxis stroke="#444" fontSize={10} tickLine={false} />
                                    <Tooltip
                                        contentStyle={{ backgroundColor: '#000', borderColor: '#333', color: '#fff' }}
                                        itemStyle={{ color: '#00f2ff' }}
                                    />
                                    <Area type="monotone" dataKey="goroutines" stroke="#00f2ff" fillOpacity={1} fill="url(#colorGro)" />
                                </AreaChart>
                            </ResponsiveContainer>
                        </div>
                    </div>
                </div>

                {/* Task Stream */}
                <div className="bg-black/40 border border-white/5 rounded-lg p-4 flex flex-col overflow-hidden">
                    <h3 className="text-xs font-bold text-gray-400 mb-4 flex items-center gap-2">
                        <Server size={14} /> ACTIVE TASKS ({tasks.length})
                    </h3>
                    <div className="flex-1 overflow-y-auto space-y-2 pr-2 custom-scrollbar">
                        {tasks.length === 0 ? (
                            <div className="h-full flex flex-col items-center justify-center opacity-30">
                                <Activity size={48} className="mb-2" />
                                <span className="text-xs font-mono">AWAITING_TASKS</span>
                            </div>
                        ) : (
                            tasks.map(task => (
                                <div key={task.id} className="bg-white/5 p-3 rounded border-l-2 border-cyber-primary flex justify-between items-center group hover:bg-white/10 transition-all">
                                    <div className="flex flex-col">
                                        <span className="text-xs font-bold text-gray-200">{task.name}</span>
                                        <span className="text-[10px] text-gray-500 font-mono">{task.id}</span>
                                    </div>
                                    <span className="text-cyber-primary font-mono text-xs font-bold">
                                        {formatDuration(task.startedAt)}
                                    </span>
                                </div>
                            ))
                        )}
                    </div>
                </div>
            </div>
        </div>
    );
}

const StatCard = ({ icon, label, value, color }: any) => (
    <div className="flex items-center gap-3 bg-white/5 px-4 py-2 rounded border border-white/5">
        <div className={`p-2 rounded bg-black/50 ${color}`}>{icon}</div>
        <div>
            <div className="text-[10px] text-gray-500 font-bold tracking-wider">{label}</div>
            <div className={`text-lg font-mono font-bold ${color}`}>{value}</div>
        </div>
    </div>
);
