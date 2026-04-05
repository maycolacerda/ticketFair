// src/App.jsx
import { useState, useEffect } from 'react'
import { useAuth } from './context/AuthContext'
import { api } from './api/client'
import Sidebar from './components/Sidebar'
import AuthPage from './pages/AuthPage'
import EventsPage from './pages/EventsPage'
import TicketsPage, { PaymentsPage, FeedPage } from './pages/TicketsPage'
import ConnectionsPage from './pages/ConnectionsPage'
import ProfilePage from './pages/ProfilePage'
import MerchantPage from './pages/MerchantPage'
import AdminPage from './pages/AdminPage'

const TITLES = {
  events: 'DISCOVER',
  feed: 'FRIENDS FEED',
  tickets: 'MY TICKETS',
  payments: 'PAYMENTS',
  connections: 'CONNECTIONS',
  profile: 'MY PROFILE',
  merchant: 'MY VENUE',
  admin: 'ADMIN PANEL',
}

export default function App() {
  const { isAuthenticated, user, merchant, admin } = useAuth()
  const [page,    setPage]    = useState('events')
  const [pending, setPending] = useState(0)

  // Poll for pending connection requests every 30s
  useEffect(() => {
    if (!isAuthenticated) return
    const poll = () =>
      api.getPendingRequests().then(({ ok, data }) => {
        if (ok) setPending((data.data || []).length)
      })
    poll()
    const t = setInterval(poll, 30000)
    return () => clearInterval(t)
  }, [isAuthenticated])

  if (!isAuthenticated) return <AuthPage />

  const fd = "'Bebas Neue', Impact, sans-serif"
  const fm = "'JetBrains Mono', monospace"
  const fb = "'Outfit', system-ui, sans-serif"

  const pages = {
    events:      <EventsPage />,
    feed:        <FeedPage />,
    tickets:     <TicketsPage />,
    payments:    <PaymentsPage />,
    connections: <ConnectionsPage onPendingChange={setPending} />,
    profile:     <ProfilePage />,
    merchant:    (merchant || admin) ? <MerchantPage /> : <EventsPage />,
    admin:       admin ? <AdminPage /> : <EventsPage />,
  }

  return (
    <div style={{ display: 'flex', minHeight: '100vh' }}>
      <Sidebar active={page} onNavigate={setPage} pendingCount={pending} />

      <main style={{ marginLeft: 60, flex: 1, display: 'flex', flexDirection: 'column', minHeight: '100vh' }}>
        {/* Topbar */}
        <header style={{
          height: 56,
          borderBottom: '1px solid var(--border)',
          background: 'rgba(9,9,9,.92)',
          backdropFilter: 'blur(12px)',
          display: 'flex', alignItems: 'center', justifyContent: 'space-between',
          padding: '0 24px',
          position: 'sticky', top: 0, zIndex: 40,
        }}>
          <div style={{ fontFamily: fd, fontSize: 17, letterSpacing: 3 }}>
            {TITLES[page] || page.toUpperCase()}
          </div>

          <div style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
            {merchant && (
              <span style={{ fontSize: 10, fontFamily: fm, color: 'var(--amber)', background: 'rgba(255,179,0,.1)', border: '1px solid rgba(255,179,0,.2)', padding: '3px 8px', borderRadius: 99 }}>
                MERCHANT
              </span>
            )}
            {admin && (
              <span style={{ fontSize: 10, fontFamily: fm, color: 'var(--red)', background: 'rgba(255,64,64,.1)', border: '1px solid rgba(255,64,64,.2)', padding: '3px 8px', borderRadius: 99 }}>
                ADMIN
              </span>
            )}
            <div
              onClick={() => setPage('profile')}
              style={{
                display: 'flex', alignItems: 'center', gap: 7,
                padding: '5px 12px',
                background: 'var(--panel)', border: '1px solid var(--border)',
                borderRadius: 99, fontSize: 12, fontWeight: 600, cursor: 'pointer',
                fontFamily: fb, transition: 'border-color .15s',
              }}
              onMouseEnter={e => e.currentTarget.style.borderColor = 'var(--lime)'}
              onMouseLeave={e => e.currentTarget.style.borderColor = 'var(--border)'}
            >
              <div style={{
                width: 24, height: 24, borderRadius: '50%',
                background: 'var(--lime)', color: 'var(--black)',
                display: 'flex', alignItems: 'center', justifyContent: 'center',
                fontSize: 11, fontWeight: 700,
              }}>
                {(user?.username || '?')[0].toUpperCase()}
              </div>
              <span>{user?.username}</span>
            </div>
          </div>
        </header>

        {/* Page content */}
        <div
          key={page}
          style={{ flex: 1, padding: '30px 24px', animation: 'fadeUp .3s ease' }}
        >
          {pages[page] || <EventsPage />}
        </div>
      </main>
    </div>
  )
}
