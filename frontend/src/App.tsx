import { useState } from 'react';
import Titlebar from './components/Titlebar';
import Sidebar from './components/Sidebar';
import Forge from './components/Forge';
import Radar from './components/Radar';

function App() {
    const [activeTab, setActiveTab] = useState('forge');

    return (
        <div className="flex flex-col h-screen bg-cyber-bg text-cyber-text overflow-hidden">
            <Titlebar />

            <div className="flex-1 flex overflow-hidden">
                <Sidebar activeTab={activeTab} onTabChange={setActiveTab} />

                <main className="flex-1 p-8 overflow-y-auto relative">
                    {/* Background decorations */}
                    <div className="absolute top-0 left-0 w-full h-full pointer-events-none opacity-20 z-0">
                        <div className="absolute top-[-10%] right-[-5%] w-96 h-96 bg-cyber-secondary rounded-full blur-[120px]"></div>
                        <div className="absolute bottom-[-10%] left-[-5%] w-96 h-96 bg-cyber-primary rounded-full blur-[120px]"></div>
                    </div>

                    <div className="relative z-10 h-full">
                        {activeTab === 'forge' && <Forge />}
                        {activeTab === 'radar' && <Radar />}

                        {(activeTab !== 'forge' && activeTab !== 'radar') && (
                            <div className="glass-panel h-full flex items-center justify-center flex-col text-gray-500">
                                <div className="text-6xl mb-4 opacity-50">🚧</div>
                                <h3 className="text-xl font-bold font-mono uppercase tracking-widest mb-2">{activeTab} Module</h3>
                                <p className="text-sm">Under Construction</p>
                            </div>
                        )}
                    </div>
                </main>
            </div>
        </div>
    );
}

export default App;
