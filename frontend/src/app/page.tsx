'use client';

import { useAccount } from 'wagmi';
import { useStore } from '../store/useStore';
import { Activity, Wallet, BarChart2, Zap, Settings, Search, LayoutDashboard, Bell } from 'lucide-react';
import { useWebSocket } from '../hooks/useWebSocket';
import { useState } from 'react';
import { ConnectButton } from '@rainbow-me/rainbowkit';

export default function Dashboard() {
  const { address, isConnected } = useAccount();
  const portfolioBalance = useStore((state) => state.portfolioBalance);
  const notifications = useStore((state) => state.notifications);
  const addNotification = useStore((state) => state.addNotification);
  
  const [tradeSymbol, setTradeSymbol] = useState('SOL');
  const [tradeAmount, setTradeAmount] = useState('100');
  const [isTrading, setIsTrading] = useState(false);
  
  // Initialize WebSocket connection
  useWebSocket();

  const handleTrade = async () => {
    setIsTrading(true);
    try {
      const res = await fetch('http://localhost:8080/api/trade/execute', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          user_id: 1, // Mock user ID
          symbol: tradeSymbol,
          amount: parseFloat(tradeAmount),
          is_buy: true,
        }),
      });
      if (!res.ok) throw new Error('Trade API failed');
      
      addNotification({
        title: 'Trade Queued',
        message: 'Your order was sent to the execution engine.',
        type: 'info'
      });
    } catch (err) {
      console.error(err);
      addNotification({
        title: 'Connection Error',
        message: 'Failed to reach the API server. Check your connection.',
        type: 'error'
      });
    } finally {
      setIsTrading(false);
    }
  };

  return (
    <div className="flex h-screen bg-zinc-950 text-white">
      {/* Sidebar */}
      <div className="w-64 border-r border-zinc-800 bg-zinc-900/50 p-4 flex flex-col gap-6">
        <div className="flex items-center gap-2 px-2 text-xl font-bold tracking-tighter text-emerald-400">
          <Zap className="w-6 h-6 fill-emerald-400" />
          CRYPTOSUM
        </div>

        <nav className="flex flex-col gap-2">
          <NavItem icon={<LayoutDashboard />} label="Dashboard" active />
          <NavItem icon={<Search />} label="Scanner" />
          <NavItem icon={<BarChart2 />} label="Markets" />
          <NavItem icon={<Activity />} label="Positions" />
          <NavItem icon={<Settings />} label="Settings" />
        </nav>

        <div className="mt-auto">
          <ConnectButton 
            chainStatus="icon" 
            showBalance={false} 
            accountStatus={{
              smallScreen: 'avatar',
              largeScreen: 'full',
            }}
          />
        </div>
      </div>

      {/* Main Content */}
      <div className="flex-1 overflow-auto p-8">
        <header className="flex justify-between items-center mb-8">
          <h1 className="text-3xl font-bold">Terminal Overview</h1>
          <div className="flex items-center gap-4">
             <div className="relative group">
               <button className="p-2 hover:bg-zinc-900 rounded-full transition">
                 <Bell className="w-5 h-5 text-zinc-400" />
                 {notifications.length > 0 && (
                   <span className="absolute top-1 right-1 w-2 h-2 bg-emerald-500 rounded-full" />
                 )}
               </button>
               {notifications.length > 0 && (
                 <div className="absolute right-0 mt-2 w-64 bg-zinc-900 border border-zinc-800 rounded-lg shadow-xl opacity-0 group-hover:opacity-100 transition-opacity p-2 z-50 pointer-events-none group-hover:pointer-events-auto">
                   {notifications.slice(-3).map((notif, i) => (
                     <div key={i} className={`p-2 text-sm border-b border-zinc-800 last:border-0 ${notif.type === 'error' ? 'text-red-400' : 'text-emerald-400'}`}>
                       <div className="font-bold">{notif.title}</div>
                       <div className="text-zinc-400">{notif.message}</div>
                     </div>
                   ))}
                 </div>
               )}
             </div>
             <div className="px-4 py-2 bg-zinc-900 rounded-lg border border-zinc-800 flex items-center gap-3">
               <div className="w-2 h-2 rounded-full bg-emerald-500 animate-pulse" />
               <span className="text-sm text-zinc-400">API Connected</span>
             </div>
          </div>
        </header>

        <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
          <StatCard title="Total Portfolio" value={`$${portfolioBalance.toFixed(2)}`} />
          <StatCard title="24h P&L" value="+$482.50" trend="up" />
          <StatCard title="Active Positions" value="3" />
        </div>

        <div className="mt-8 grid grid-cols-1 lg:grid-cols-3 gap-6">
          <div className="col-span-2 bg-zinc-900 border border-zinc-800 rounded-xl p-6 h-96">
            <h2 className="text-lg font-semibold mb-4">Market Scanner</h2>
            <div className="flex items-center justify-center h-full text-zinc-500">
              Select a token to view AI analysis
            </div>
          </div>
          <div className="col-span-1 bg-zinc-900 border border-zinc-800 rounded-xl p-6">
            <h2 className="text-lg font-semibold mb-4">Quick Trade</h2>
            {/* Quick Trade Stub */}
            <div className="space-y-4">
              <div>
                <label className="text-xs text-zinc-400 uppercase">Token</label>
                <input 
                  type="text" 
                  value={tradeSymbol}
                  onChange={(e) => setTradeSymbol(e.target.value)}
                  className="w-full bg-zinc-950 border border-zinc-800 rounded p-2 mt-1" 
                />
              </div>
              <div>
                <label className="text-xs text-zinc-400 uppercase">Amount (USDC)</label>
                <input 
                  type="number" 
                  value={tradeAmount}
                  onChange={(e) => setTradeAmount(e.target.value)}
                  className="w-full bg-zinc-950 border border-zinc-800 rounded p-2 mt-1" 
                />
              </div>
              <button 
                onClick={handleTrade}
                disabled={isTrading}
                className={`w-full py-3 font-bold rounded mt-4 flex items-center justify-center gap-2 transition ${isTrading ? 'bg-emerald-500/50 cursor-not-allowed text-zinc-800' : 'bg-emerald-500 hover:bg-emerald-400 active:bg-emerald-600 text-zinc-950 shadow-[0_0_15px_rgba(16,185,129,0.3)]'}`}
              >
                {isTrading ? (
                  <>
                    <div className="w-4 h-4 border-2 border-zinc-800 border-t-transparent rounded-full animate-spin" />
                    Executing...
                  </>
                ) : (
                  'Execute Buy'
                )}
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}

function NavItem({ icon, label, active = false }: { icon: React.ReactNode, label: string, active?: boolean }) {
  return (
    <button className={`flex items-center gap-3 px-4 py-3 rounded-lg transition-colors ${active ? 'bg-zinc-800 text-white' : 'text-zinc-400 hover:bg-zinc-800/50 hover:text-white'}`}>
      <div className="w-5 h-5">{icon}</div>
      <span className="font-medium">{label}</span>
    </button>
  );
}

function StatCard({ title, value, trend }: { title: string, value: string, trend?: 'up' | 'down' }) {
  return (
    <div className="bg-zinc-900 border border-zinc-800 rounded-xl p-6">
      <h3 className="text-zinc-400 text-sm font-medium">{title}</h3>
      <div className="mt-2 flex items-baseline gap-2">
        <span className="text-3xl font-bold">{value}</span>
        {trend && (
          <span className={`text-sm font-medium ${trend === 'up' ? 'text-emerald-400' : 'text-red-400'}`}>
            {trend === 'up' ? '↗' : '↘'}
          </span>
        )}
      </div>
    </div>
  );
}
