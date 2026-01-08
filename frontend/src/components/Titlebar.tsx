import { useState } from 'react';
import { WindowMinimise, WindowToggleMaximise, Quit } from '../../wailsjs/runtime/runtime';

export default function Titlebar() {
    return (
        <div className="h-8 bg-cyber-bg w-full flex items-center justify-between select-none" style={{ widows: '1' }}>
            <div className="flex-1 h-full flex items-center pl-4 app-drag-region" style={{ '--wails-draggable': 'drag' } as any}>
                <span className="text-xs font-mono text-cyber-primary tracking-widest opacity-80">NEXUS OPS // HUB</span>
            </div>
            <div className="flex bg-cyber-bg">
                <button onClick={WindowMinimise} className="h-8 w-10 hover:bg-gray-800 text-gray-400 flex items-center justify-center transition-colors">
                    _
                </button>
                <button onClick={WindowToggleMaximise} className="h-8 w-10 hover:bg-gray-800 text-gray-400 flex items-center justify-center transition-colors">
                    []
                </button>
                <button onClick={Quit} className="h-8 w-10 hover:bg-red-500 hover:text-white text-gray-400 flex items-center justify-center transition-colors">
                    X
                </button>
            </div>
        </div>
    );
}
