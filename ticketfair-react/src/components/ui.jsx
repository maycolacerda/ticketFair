// src/components/ui.jsx
import { useState } from 'react'

// ── helpers ──────────────────────────────────────────
const fd = "'Bebas Neue', Impact, sans-serif"
const fb = "'Outfit', system-ui, sans-serif"
const fm = "'JetBrains Mono', monospace"

// ── Spinner ──────────────────────────────────────────
export function Spinner({ size = 16, color = 'var(--lime)' }) {
  return <div style={{ width:size, height:size, flexShrink:0, border:'2px solid rgba(255,255,255,.1)', borderTopColor:color, borderRadius:'50%', animation:'spin .6s linear infinite' }} />
}

// ── Btn ───────────────────────────────────────────────
const BV = {
  primary: { bg:'var(--lime)',                color:'var(--black)', border:'none' },
  ghost:   { bg:'transparent',               color:'var(--muted)', border:'1px solid var(--border)' },
  danger:  { bg:'rgba(255,64,64,.1)',         color:'var(--red)',   border:'1px solid rgba(255,64,64,.25)' },
  success: { bg:'rgba(0,230,118,.1)',         color:'var(--green)', border:'1px solid rgba(0,230,118,.25)' },
  blue:    { bg:'rgba(68,138,255,.1)',        color:'var(--blue)',  border:'1px solid rgba(68,138,255,.25)' },
  amber:   { bg:'rgba(255,179,0,.1)',         color:'var(--amber)', border:'1px solid rgba(255,179,0,.25)' },
}
const BP = { sm:'5px 10px', md:'9px 18px', lg:'13px 28px' }
const BF = { sm:11, md:13, lg:14 }

export function Btn({ children, variant='ghost', size='md', onClick, disabled, loading, style={}, type='button' }) {
  const v = BV[variant] || BV.ghost
  return (
    <button type={type} onClick={onClick} disabled={disabled||loading} style={{ display:'inline-flex', alignItems:'center', gap:6, padding:BP[size]||BP.md, borderRadius:6, fontFamily:fb, fontSize:BF[size]||13, fontWeight:600, cursor:(disabled||loading)?'not-allowed':'pointer', opacity:disabled?.45:1, border:v.border, background:v.bg, color:v.color, transition:'all .15s', whiteSpace:'nowrap', ...style }}>
      {loading ? <Spinner size={BF[size]||13} color={v.color} /> : children}
    </button>
  )
}

// ── Badge ─────────────────────────────────────────────
const BC = {
  lime:   { bg:'rgba(212,255,0,.1)',    c:'var(--lime)',   b:'rgba(212,255,0,.2)' },
  amber:  { bg:'rgba(255,179,0,.1)',    c:'var(--amber)',  b:'rgba(255,179,0,.2)' },
  red:    { bg:'rgba(255,64,64,.1)',    c:'var(--red)',    b:'rgba(255,64,64,.2)' },
  green:  { bg:'rgba(0,230,118,.1)',    c:'var(--green)',  b:'rgba(0,230,118,.2)' },
  blue:   { bg:'rgba(68,138,255,.1)',   c:'var(--blue)',   b:'rgba(68,138,255,.2)' },
  purple: { bg:'rgba(179,136,255,.1)',  c:'var(--purple)', b:'rgba(179,136,255,.2)' },
  muted:  { bg:'rgba(102,102,102,.08)', c:'var(--muted)', b:'var(--faint)' },
}
export function Badge({ children, color='muted', style={} }) {
  const c = BC[color] || BC.muted
  return <span style={{ display:'inline-flex', alignItems:'center', gap:3, padding:'2px 7px', borderRadius:99, fontSize:10, fontFamily:fm, background:c.bg, color:c.c, border:`1px solid ${c.b}`, ...style }}>{children}</span>
}

// ── Modal ─────────────────────────────────────────────
export function Modal({ open, onClose, title, children, width=480 }) {
  if (!open) return null
  return (
    <div onClick={e => e.target===e.currentTarget&&onClose()} style={{ position:'fixed', inset:0, background:'rgba(0,0,0,.8)', backdropFilter:'blur(8px)', zIndex:200, display:'flex', alignItems:'center', justifyContent:'center', padding:20 }}>
      <div style={{ background:'var(--dark)', border:'1px solid var(--border)', borderRadius:12, width:'100%', maxWidth:width, maxHeight:'90vh', overflowY:'auto', animation:'modalIn .2s ease' }}>
        <div style={{ display:'flex', alignItems:'center', justifyContent:'space-between', padding:'20px 24px 0' }}>
          <h3 style={{ fontFamily:fd, fontSize:22, letterSpacing:2 }}>{title}</h3>
          <button onClick={onClose} style={{ background:'none', border:'none', color:'var(--muted)', fontSize:20, cursor:'pointer', lineHeight:1 }}>✕</button>
        </div>
        <div style={{ padding:'16px 24px 24px' }}>{children}</div>
      </div>
    </div>
  )
}

// ── Field ─────────────────────────────────────────────
export function Field({ label, children, error, style={} }) {
  return (
    <div style={{ marginBottom:16, ...style }}>
      {label && <label style={{ display:'block', fontSize:10, fontFamily:fm, color:'var(--muted)', letterSpacing:1, textTransform:'uppercase', marginBottom:6 }}>{label}</label>}
      {children}
      {error && <p style={{ fontSize:11, color:'var(--red)', fontFamily:fm, marginTop:4 }}>{error}</p>}
    </div>
  )
}

// ── Alert ─────────────────────────────────────────────
const AC = {
  error:   { bg:'rgba(255,64,64,.06)',   b:'rgba(255,64,64,.2)',   c:'var(--red)' },
  success: { bg:'rgba(0,230,118,.06)',   b:'rgba(0,230,118,.2)',   c:'var(--green)' },
  info:    { bg:'rgba(68,138,255,.06)',  b:'rgba(68,138,255,.2)',  c:'var(--blue)' },
  warning: { bg:'rgba(255,179,0,.06)',   b:'rgba(255,179,0,.2)',   c:'var(--amber)' },
}
export function Alert({ children, type='error', style={} }) {
  const c = AC[type]||AC.error
  return <div style={{ padding:'10px 14px', borderRadius:6, fontSize:12, fontFamily:fm, background:c.bg, border:`1px solid ${c.b}`, color:c.c, lineHeight:1.5, ...style }}>{children}</div>
}

// ── Skeleton ──────────────────────────────────────────
export function Skeleton({ w='100%', h=14, style={} }) {
  return <div style={{ width:w, height:h, borderRadius:4, background:'linear-gradient(90deg,var(--surf) 25%,var(--panel) 50%,var(--surf) 75%)', backgroundSize:'200% 100%', animation:'shimmer 1.5s infinite', ...style }} />
}

// ── Empty ─────────────────────────────────────────────
export function Empty({ icon='🎟', title, sub, action }) {
  return (
    <div style={{ textAlign:'center', padding:'56px 32px' }}>
      <div style={{ fontSize:48, opacity:.2, marginBottom:16 }}>{icon}</div>
      <div style={{ fontFamily:fd, fontSize:22, letterSpacing:2, marginBottom:8 }}>{title}</div>
      {sub && <p style={{ fontSize:13, color:'var(--muted)' }}>{sub}</p>}
      {action && <div style={{ marginTop:20 }}>{action}</div>}
    </div>
  )
}

// ── Card ──────────────────────────────────────────────
export function Card({ children, onClick, style={} }) {
  const [h, setH] = useState(false)
  return (
    <div onClick={onClick} onMouseEnter={()=>setH(true)} onMouseLeave={()=>setH(false)} style={{ background:'var(--surf)', border:`1px solid ${h&&onClick?'var(--faint)':'var(--border)'}`, borderRadius:10, overflow:'hidden', transition:'all .2s', transform:h&&onClick?'translateY(-3px)':'none', cursor:onClick?'pointer':'default', ...style }}>
      {children}
    </div>
  )
}

// ── Tabs ──────────────────────────────────────────────
export function Tabs({ tabs, active, onChange, style={} }) {
  return (
    <div style={{ display:'flex', background:'var(--panel)', borderRadius:8, padding:4, gap:4, marginBottom:24, ...style }}>
      {tabs.map(t => (
        <button key={t.value} onClick={()=>onChange(t.value)} style={{ flex:1, padding:'8px 10px', borderRadius:6, border:'none', background:active===t.value?'var(--surf)':'transparent', color:active===t.value?'var(--text)':'var(--muted)', fontFamily:fb, fontSize:13, fontWeight:600, cursor:'pointer', transition:'all .15s', display:'flex', alignItems:'center', justifyContent:'center', gap:6 }}>
          {t.label}
          {t.count > 0 && <Badge color="red" style={{padding:'1px 5px',fontSize:9}}>{t.count}</Badge>}
        </button>
      ))}
    </div>
  )
}

// ── Section Header ────────────────────────────────────
export function SectionHeader({ title, sub, action }) {
  return (
    <div style={{ display:'flex', alignItems:'flex-start', justifyContent:'space-between', marginBottom:24 }}>
      <div>
        <h2 style={{ fontFamily:fd, fontSize:'clamp(22px,4vw,34px)', letterSpacing:2 }}>{title}</h2>
        {sub && <p style={{ fontSize:11, fontFamily:fm, color:'var(--muted)', marginTop:4 }}>{sub}</p>}
      </div>
      {action}
    </div>
  )
}

// ── Stat ──────────────────────────────────────────────
export function Stat({ value, label, color='var(--lime)' }) {
  return (
    <div style={{ background:'var(--surf)', border:'1px solid var(--border)', borderRadius:10, padding:20 }}>
      <div style={{ fontFamily:fd, fontSize:44, color, letterSpacing:1, lineHeight:1 }}>{value}</div>
      <div style={{ fontSize:10, fontFamily:fm, color:'var(--muted)', marginTop:6, letterSpacing:1, textTransform:'uppercase' }}>{label}</div>
    </div>
  )
}

// ── Avatar ────────────────────────────────────────────
export function Avatar({ name='?', size=32 }) {
  return (
    <div style={{ width:size, height:size, borderRadius:'50%', background:'var(--faint)', color:'var(--text)', display:'flex', alignItems:'center', justifyContent:'center', fontFamily:fd, fontSize:size*.42, flexShrink:0 }}>
      {(name[0]||'?').toUpperCase()}
    </div>
  )
}

// ── Divider ───────────────────────────────────────────
export function Divider({ style={} }) {
  return <div style={{ height:1, background:'var(--border)', margin:'18px 0', ...style }} />
}

// ── SearchInput ───────────────────────────────────────
export function SearchInput({ value, onChange, placeholder='' }) {
  return (
    <div style={{ display:'flex', alignItems:'center', background:'var(--panel)', border:'1px solid var(--border)', borderRadius:8, overflow:'hidden', marginBottom:24 }}>
      <span style={{ padding:'0 14px', color:'var(--muted)', fontSize:16, flexShrink:0 }}>⌕</span>
      <input value={value} onChange={e=>onChange(e.target.value)} placeholder={placeholder} style={{ flex:1, border:'none', padding:'12px 0', background:'transparent', outline:'none' }} />
    </div>
  )
}

// ── Panel Section ─────────────────────────────────────
export function PanelSection({ title, action, children, style={} }) {
  return (
    <div style={{ background:'var(--surf)', border:'1px solid var(--border)', borderRadius:10, padding:20, marginBottom:20, ...style }}>
      <div style={{ display:'flex', alignItems:'center', justifyContent:'space-between', marginBottom:16, paddingBottom:12, borderBottom:'1px solid var(--border)' }}>
        <div style={{ fontFamily:fd, fontSize:17, letterSpacing:1 }}>{title}</div>
        {action}
      </div>
      {children}
    </div>
  )
}

// ── Row (list item) ───────────────────────────────────
export function Row({ icon, title, sub, right, onClick, style={} }) {
  const [h, setH] = useState(false)
  return (
    <div onClick={onClick} onMouseEnter={()=>setH(true)} onMouseLeave={()=>setH(false)} style={{ display:'flex', alignItems:'center', gap:12, padding:'11px 13px', background:'var(--surf)', border:`1px solid ${h?'var(--faint)':'var(--border)'}`, borderRadius:8, marginBottom:8, transition:'border-color .15s', cursor:onClick?'pointer':'default', ...style }}>
      {icon && <span style={{ fontSize:20, flexShrink:0 }}>{icon}</span>}
      <div style={{ flex:1, minWidth:0 }}>
        <div style={{ fontSize:14, fontWeight:600, overflow:'hidden', textOverflow:'ellipsis', whiteSpace:'nowrap' }}>{title}</div>
        {sub && <div style={{ fontSize:11, fontFamily:fm, color:'var(--muted)', marginTop:1 }}>{sub}</div>}
      </div>
      {right && <div style={{ display:'flex', gap:6, flexShrink:0, alignItems:'center' }}>{right}</div>}
    </div>
  )
}
