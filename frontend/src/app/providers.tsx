'use client';

import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { WagmiProvider, createConfig, http } from 'wagmi';
import { mainnet, bsc, solana } from 'wagmi/chains';
import { ReactNode, useState } from 'react';

// Using mock config for MVP to prevent requirement of WalletConnect Project ID right now
export const config = createConfig({
  chains: [mainnet, bsc],
  transports: {
    [mainnet.id]: http(),
    [bsc.id]: http(),
  },
})

export default function Providers({ children }: { children: ReactNode }) {
  const [queryClient] = useState(() => new QueryClient());

  return (
    <WagmiProvider config={config}>
      <QueryClientProvider client={queryClient}>
        {children}
      </QueryClientProvider>
    </WagmiProvider>
  );
}
