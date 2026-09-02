'use client';

import '@rainbow-me/rainbowkit/styles.css';
import { getDefaultConfig, RainbowKitProvider, darkTheme } from '@rainbow-me/rainbowkit';
import { WagmiProvider } from 'wagmi';
import { mainnet, bsc, polygon, optimism, arbitrum, base } from 'wagmi/chains';
import { QueryClientProvider, QueryClient } from '@tanstack/react-query';
import { ReactNode, useState } from 'react';

// Using a public default for the MVP, but ideally loaded from env
const projectId = process.env.NEXT_PUBLIC_WALLETCONNECT_PROJECT_ID || '1deccdc4e8c16dcb4dafb4b50c042971';

export const config = getDefaultConfig({
  appName: 'Cryptosum',
  projectId: projectId,
  chains: [mainnet, bsc, polygon, optimism, arbitrum, base],
  ssr: true,
});

export default function Providers({ children }: { children: ReactNode }) {
  const [queryClient] = useState(() => new QueryClient());

  return (
    <WagmiProvider config={config}>
      <QueryClientProvider client={queryClient}>
        <RainbowKitProvider theme={darkTheme()}>
          {children}
        </RainbowKitProvider>
      </QueryClientProvider>
    </WagmiProvider>
  );
}
