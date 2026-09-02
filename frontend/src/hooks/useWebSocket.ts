import { useEffect, useRef } from 'react';
import { useStore } from '../store/useStore';

const WS_URL = 'ws://localhost:8080/ws/stream';

export function useWebSocket() {
  const ws = useRef<WebSocket | null>(null);
  const addNotification = useStore((state) => state.addNotification);

  useEffect(() => {
    ws.current = new WebSocket(WS_URL);

    ws.current.onopen = () => {
      console.log('WebSocket connected');
    };

    ws.current.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data);
        console.log('WS Message:', data);
        
        // Handle trade execution results
        if (data.type === 'TRADE_RESULT') {
          addNotification({
            title: data.payload.success ? 'Trade Executed' : 'Trade Failed',
            message: data.payload.message,
            type: data.payload.success ? 'success' : 'error'
          });
        }
        
      } catch (err) {
        console.error('Failed to parse WS message', err);
      }
    };

    ws.current.onclose = () => {
      console.log('WebSocket disconnected');
    };

    return () => {
      if (ws.current) {
        ws.current.close();
      }
    };
  }, [addNotification]);

  const sendMessage = (message: any) => {
    if (ws.current && ws.current.readyState === WebSocket.OPEN) {
      ws.current.send(JSON.stringify(message));
    }
  };

  return { sendMessage };
}
