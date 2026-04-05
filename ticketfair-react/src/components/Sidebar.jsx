// src/components/Sidebar.jsx
import { useState } from 'react'
import { useAuth } from '../context/AuthContext'
import { Avatar } from './ui'

const NAV = [
  { id:'events',      icon:'🎪', label:'Discover' },
  { id:'feed',        icon:'👥', label:'Friends Feed' },
  { id:'tickets',     icon:'🎟', label:'My Tickets' },
  { id:'payments',    icon:'💳', label:'Payments' },
  { id:'connections', icon:'🔗', label:'Connections' },
  { id:'profile',     icon:'👤', label:'Profile' },
]

const fd = "'Bebas Neue', Impact, sans-serif"
const fb = "'Outfit', system-ui, sans-serif"
const fm = "'JetBrains Mono', monospace"

export default function Sidebar({ active, onNavigate, pendingCount }) {
  const { user, merchant, admin, logout } = useAuth()
  const [open, setOpen] = useState(false)

  const w = open ? 220 : 60

  return (
    <>
      {/* Toggle */}
      <button
        onClick={() => setOpen(o => !o)}
        style={{ position:'fixed', top:13, left:13, zIndex:60, width:32, height:32, background:'var(--panel)', border:'1px solid var(--border)', borderRadius:6, cursor:'pointer', display:'flex', alignItems:'center', justifyContent:'center', fontSize:14, color:'var(--muted)', transition:'border-color .15s' }}
        onMouseEnter={e => e.currentTarget.style.borderColor = 'var(--lime)'}
        onMouseLeave={e => e.currentTarget.style.borderColor = 'var(--border)'}
      >{open ? '←' : '☰'}</button>

      <aside style={{ position:'fixed', top:0, left:0, height:'100vh', width:w, background:'var(--dark)', borderRight:'1px solid var(--border)', display:'flex', flexDirection:'column', padding:'16px 0', zIndex:50, transition:'width .25s ease', overflow:'hidden' }}>
        {/* Logo */}
        <div style={{ padding:'8px 16px 18px', borderBottom:'1px solid var(--border)', marginBottom:10, marginTop:6 }}>
          <div style={{ fontFamily:fd, fontSize:open?18:14, letterSpacing:3, color:'var(--lime)', whiteSpace:'nowrap', transition:'font-size .2s' }}>
            {open ? 'TICKETFAIR' : 'TF'}
          </div>
        </div>

        {/* Nav */}
        <div style={{ flex:1, padding:'0 6px', display:'flex', flexDirection:'column', gap:2 }}>
          {NAV.map(item => (
            <NavBtn key={item.id} item={item} active={active===item.id} open={open} onClick={() => onNavigate(item.id)} badge={item.id==='connections' ? pendingCount : 0} />
          ))}
          {(merchant || admin) && (
            <NavBtn item={{ id:'merchant', icon:'🏪', label:'My Venue' }} active={active==='merchant'} open={open} onClick={() => onNavigate('merchant')} />
          )}
          {admin && (
            <NavBtn item={{ id:'admin', icon:'⚙️', label:'Admin' }} active={active==='admin'} open={open} onClick={() => onNavigate('admin')} />
          )}
        </div>

        {/* User / logout */}
        <div style={{ padding:'10px 6px 0', borderTop:'1px solid var(--border)' }}>
          <div onClick={logout} style={{ display:'flex', alignItems:'center', gap:10, padding:'8px', borderRadius:8, background:'var(--panel)', border:'1px solid var(--border)', cursor:'pointer', overflow:'hidden', transition:'border-color .15s' }}
            onMouseEnter={e => e.currentTarget.style.borderColor = 'var(--red)'}
            onMouseLeave={e => e.currentTarget.style.borderColor = 'var(--border)'}
          >
            <Avatar name={user?.username||'?'} size={28} />
            {open && (
              <div style={{ overflow:'hidden' }}>
                <div style={{ fontSize:12, fontWeight:700, whiteSpace:'nowrap', overflow:'hidden', textOverflow:'ellipsis' }}>{user?.username}</div>
                <div style={{ fontSize:10, color:'var(--red)', fontFamily:fm }}>Sign out</div>
              </div>
            )}
          </div>
        </div>
      </aside>
    </>
  )
}

function NavBtn({ item, active, open, onClick, badge=0 }) {
  const [h, setH] = useState(false)
  const fd_ = "'Bebas Neue', Impact, sans-serif"
  const fb_ = "'Outfit', system-ui, sans-serif"

  return (
    <div onClick={onClick} onMouseEnter={()=>setH(true)} onMouseLeave={()=>setH(false)} style={{ display:'flex', alignItems:'center', gap:10, padding:'9px 10px', borderRadius:8, cursor:'pointer', background:active?'rgba(212,255,0,.07)':h?'var(--panel)':'transparent', border:`1px solid ${active?'rgba(212,255,0,.2)':'transparent'}`, color:active?'var(--lime)':h?'var(--text)':'var(--muted)', transition:'all .15s', position:'relative', overflow:'hidden', minHeight:38 }}>
      <span style={{ fontSize:17, flexShrink:0, width:20, textAlign:'center' }}>{item.icon}</span>
      {open && <span style={{ fontSize:13, fontFamily:fb_, fontWeight:600, whiteSpace:'nowrap' }}>{item.label}</span>}
      {badge > 0 && (
        <div style={{ position:open?'static':'absolute', top:open?'auto':5, right:open?'auto':5, marginLeft:open?'auto':0, background:'var(--red)', color:'#fff', borderRadius:99, padding:open?'1px 5px':'1px 4px', fontSize:9, fontFamily:"'JetBrains Mono',monospace", minWidth:14, textAlign:'center' }}>
          {badge}
        </div>
      )}
    </div>
  )
}
