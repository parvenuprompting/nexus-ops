type Props = {
    activeTab: string;
    onTabChange: (tab: string) => void;
};

export default function Sidebar({ activeTab, onTabChange }: Props) {
    const tabs = [
        { id: 'forge', label: 'Forge', icon: '⚡' },
        { id: 'radar', label: 'Radar', icon: '📡' },
        { id: 'vault', label: 'The Vault', icon: '🔒' },
        { id: 'siphon', label: 'Siphon', icon: '🌊' },
        { id: 'terminal', label: 'Terminal', icon: '💻' },
    ];

    return (
        <div className="w-64 h-full bg-cyber-bg/50 backdrop-blur-md border-r border-cyber-border flex flex-col p-4">
            <div className="mb-8">
                <h1 className="text-2xl font-bold bg-gradient-to-r from-cyber-primary to-cyber-secondary bg-clip-text text-transparent">
                    NEXUS HUB
                </h1>
                <p className="text-xs text-gray-500 tracking-wider mt-1">v1.2.0-wails</p>
            </div>

            <nav className="flex-1 space-y-2">
                {tabs.map((tab) => (
                    <button
                        key={tab.id}
                        onClick={() => onTabChange(tab.id)}
                        className={`w-full text-left px-4 py-3 rounded-lg flex items-center gap-3 transition-all duration-200 ${activeTab === tab.id
                                ? 'bg-cyber-primary/10 text-cyber-primary border border-cyber-primary/20 shadow-[0_0_15px_rgba(0,242,255,0.2)]'
                                : 'text-gray-400 hover:text-white hover:bg-white/5'
                            }`}
                    >
                        <span>{tab.icon}</span>
                        <span className="font-medium tracking-wide">{tab.label}</span>
                    </button>
                ))}
            </nav>

            <div className="mt-auto pt-4 border-t border-cyber-border">
                <div className="flex items-center gap-2 text-xs text-gray-500">
                    <div className="w-2 h-2 rounded-full bg-green-500 animate-pulse"></div>
                    SYSTEM ONLINE
                </div>
            </div>
        </div>
    );
}
