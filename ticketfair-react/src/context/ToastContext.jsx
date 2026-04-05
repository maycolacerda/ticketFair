// src/context/ToastContext.jsx
import { createContext, useContext, useState, useCallback } from 'react'

const Ctx = createContext(null)
const COLORS = { success: '#00e676', error: '#ff4040', info: '#448aff', warning: '#ffb300' }
const ICONS  = { success: '✓', error: '✕', info: 'ℹ', warning: '⚠' }

export function ToastProvider({ children }) {
  const [toasts, setToasts] = useState([])

  const toast = useCallback((message, type = 'info', ms = 4500) => {
    const id = Date.now() + Math.random()
    setToasts(p => [...p, { id, message, type }])
    setTimeout(() => setToasts(p => p.filter(t => t.id !== id)), ms)
  }, [])

  return (
    <Ctx.Provider value={toast}>
      {children}
      <div style={{ position:'fixed', bottom:24, right:24, zIndex:9999, display:'flex', flexDirection:'column', gap:8 }}>
        {toasts.map(t => (
          <div key={t.id} style={{
            background: '#181818', border: '1px solid #2a2a2a',
            borderLeft: `3px solid ${COLORS[t.type]}`,
            borderRadius: 8, padding: '11px 16px',
            display: 'flex', alignItems: 'center', gap: 10,
            minWidth: 260, maxWidth: 380, fontSize: 13, fontWeight: 600,
            fontFamily: "'Outfit', sans-serif", color: '#f2ede4',
            animation: 'toastSlide .3s ease',
            boxShadow: '0 8px 24px rgba(0,0,0,.5)',
          }}>
            <span style={{ color: COLORS[t.type], flexShrink: 0 }}>{ICONS[t.type]}</span>
            <span style={{ lineHeight: 1.45 }}>{t.message}</span>
          </div>
        ))}
      </div>
    </Ctx.Provider>
  )
}

export const useToast = () => useContext(Ctx)
