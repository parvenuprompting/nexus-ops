type Props = {
    activeTab: string;
    onTabChange: (tab: string) => void;
};

export default function Sidebar({ activeTab, onTabChange }: Props) {
    const tabs = [
        { id: 'forge', label: 'Forge' },
        { id: 'radar', label: 'Radar' },
        { id: 'vault', label: 'The Vault' },
        { id: 'siphon', label: 'Siphon' },
        { id: 'terminal', label: 'Terminal' },
    ];

    return (
        <div className="w-64 h-full bg-black/40 backdrop-blur-xl border-r border-white/10 flex flex-col p-4 shadow-2xl z-20">
            <div className="mb-6 mt-2 text-center">
                <h1 className="text-xl font-black font-mono tracking-[0.2em] text-gray-200">
                    NEXUS HUB
                </h1>
                <div className="w-16 h-0.5 bg-cyber-primary/50 mx-auto mt-2 blur-[1px]"></div>
            </div>

            <nav className="flex-1 flex flex-col gap-3">
                {tabs.map((tab) => (
                    <button
                        key={tab.id}
                        onClick={() => onTabChange(tab.id)}
                        className={`
                            relative w-full py-3 px-4 rounded-sm font-bold tracking-wide transition-all duration-200 uppercase text-sm
                            border backdrop-blur-sm
                            ${activeTab === tab.id
                                ? 'bg-cyber-primary/10 border-cyber-primary/60 text-cyber-primary shadow-[0_0_15px_rgba(0,242,255,0.15)]'
                                : 'bg-white/5 border-white/5 text-gray-400 hover:bg-white/10 hover:text-gray-200 hover:border-white/20'
                            }
                        `}
                    >
                        {/* Active Indicator Line */}
                        {activeTab === tab.id && (
                            <div className="absolute left-0 top-0 bottom-0 w-1 bg-cyber-primary shadow-[0_0_8px_#00f2ff]"></div>
                        )}
                        {tab.label}
                    </button>
                ))}
            </nav>

            {/* Footer */}
            <div className="mt-auto pt-6 border-t border-white/10 text-center">
                <div className="text-[10px] text-gray-500 font-mono mb-2">
                    © 2026 Tiëndo Welles
                </div>
                <div className="flex items-center justify-center gap-2 text-[10px] uppercase font-bold tracking-widest text-green-500">
                    <div className="w-1.5 h-1.5 rounded-full bg-green-500 shadow-[0_0_8px_#22c55e]"></div>
                    ONLINE
                </div>
            </div>
        </div>
    );
}
