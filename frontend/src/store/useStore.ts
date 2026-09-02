import { create } from 'zustand';

export interface Notification {
  title: string;
  message: string;
  type: 'success' | 'error' | 'info';
}

interface AppState {
  portfolioBalance: number;
  setPortfolioBalance: (balance: number) => void;
  activeToken: string | null;
  setActiveToken: (token: string | null) => void;
  notifications: Notification[];
  addNotification: (notification: Notification) => void;
  clearNotifications: () => void;
}

export const useStore = create<AppState>((set) => ({
  portfolioBalance: 0,
  setPortfolioBalance: (balance) => set({ portfolioBalance: balance }),
  activeToken: null,
  setActiveToken: (token) => set({ activeToken: token }),
  notifications: [],
  addNotification: (notification) => set((state) => ({ notifications: [...state.notifications, notification] })),
  clearNotifications: () => set({ notifications: [] }),
}));
