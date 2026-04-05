// src/pages/EventsPage.jsx
import { useState, useEffect } from 'react'
import { api } from '../api/client'
import { useAuth } from '../context/AuthContext'
import { useToast } from '../context/ToastContext'
import { Btn, Badge, Modal, Field, Alert, Empty, Spinner, SectionHeader, Card, SearchInput, Skeleton } from '../components/ui'

const fd = "'Bebas Neue', Impact, sans-serif"
const fm = "'JetBrains Mono', monospace"

const CAT_COLOR = { vip:'amber', early_bird:'lime', complimentary:'green', reserved:'blue', demographic:'purple' }
const EVT_EMOJI = n => {
  const s = (n||'').toLowerCase()
  if (s.includes('rock')||s.includes('metal')) return '🎸'
  if (s.includes('jazz')||s.includes('blues')) return '🎷'
  if (s.includes('festival')||s.includes('summer')||s.includes('verão')) return '🎆'
  if (s.includes('theater')||s.includes('play')||s.includes('teatro')) return '🎭'
  if (s.includes('comedy')) return '🎤'
  return '🎟'
}

const fmtDate = iso => iso ? new Date(iso).toLocaleDateString('en-US', { month:'short', day:'numeric', year:'numeric' }) : '—'
const fmtCents = c => c === 0 ? 'Free' : `R$ ${(c/100).toFixed(2).replace('.',',')}`

export default function EventsPage() {
  const { isAuthenticated } = useAuth()
  const toast = useToast()

  const [events, setEvents] = useState([])
  const [total,  setTotal]  = useState(0)
  const [loading, setLoading] = useState(true)
  const [search, setSearch] = useState('')

  // Event detail
  const [ev,      setEv]      = useState(null)
  const [tts,     setTts]     = useState([])
  const [ttBusy,  setTtBusy]  = useState(false)
  const [selTT,   setSelTT]   = useState(null)
  const [qty,     setQty]     = useState(1)
  const [buyErr,  setBuyErr]  = useState('')
  const [buying,  setBuying]  = useState(false)

  // Payment modal
  const [payOpen, setPayOpen] = useState(false)
  const [intent,  setIntent]  = useState(null)

  useEffect(() => { load() }, [])

  async function load() {
    setLoading(true)
    const { ok, data } = await api.getEvents('page=1&limit=100')
    setLoading(false)
    if (ok) { setEvents(data.data || []); setTotal(data.total || 0) }
  }

  async function openEv(e) {
    setEv(e); setSelTT(null); setQty(1); setBuyErr('')
    setTtBusy(true); setTts([])
    const { ok, data } = await api.getTicketTypes(e.event_id)
    setTtBusy(false)
    if (ok) setTts(data.data || [])
  }

  async function purchase() {
    if (!selTT) { setBuyErr('Select a ticket type first'); return }
    if (!isAuthenticated) { toast('Sign in to buy tickets', 'error'); return }
    setBuying(true); setBuyErr('')
    const { ok, data } = await api.createIntent({ event_id: ev.event_id, ticket_type_id: selTT.ticket_type_id, quantity: qty })
    setBuying(false)
    if (!ok) { setBuyErr(data.error || data.data?.error || (data.status === 503 ? 'Stripe not configured' : 'Failed')); return }
    setIntent(data.data)
    setEv(null)
    setPayOpen(true)
  }

  const filtered = events.filter(e => !search || e.name.toLowerCase().includes(search.toLowerCase()) || e.location.toLowerCase().includes(search.toLowerCase()))
  const total$ = selTT ? fmtCents(selTT.price_cents * qty) : '—'

  return (
    <div>
      <SectionHeader title="UPCOMING EVENTS" sub={`${total} events available`} />
      <SearchInput value={search} onChange={setSearch} placeholder="Search events, locations, artists…" />

      {loading ? (
        <div style={{ display:'grid', gridTemplateColumns:'repeat(auto-fill,minmax(240px,1fr))', gap:20 }}>
          {[1,2,3,4,5,6].map(i => <CardSkeleton key={i} />)}
        </div>
      ) : filtered.length === 0 ? (
        <Empty icon="🎪" title="NO EVENTS" sub="Check back soon" />
      ) : (
        <div style={{ display:'grid', gridTemplateColumns:'repeat(auto-fill,minmax(240px,1fr))', gap:20 }}>
          {filtered.map((e, i) => <EvCard key={e.event_id} ev={e} idx={i} onClick={() => openEv(e)} />)}
        </div>
      )}

      {/* Event detail */}
      <Modal open={!!ev} onClose={() => setEv(null)} title={ev?.name||''} width={600}>
        {ev && <>
          <div style={{ height:190, borderRadius:8, marginBottom:18, background:'var(--panel)', display:'flex', alignItems:'center', justifyContent:'center', fontSize:60, position:'relative', overflow:'hidden' }}>
            {EVT_EMOJI(ev.name)}
            {ev.image_url && <img src={ev.image_url} alt="" style={{ position:'absolute', inset:0, width:'100%', height:'100%', objectFit:'cover' }} onError={e=>e.target.style.display='none'} />}
          </div>

          <div style={{ display:'grid', gridTemplateColumns:'1fr 1fr', gap:10, marginBottom:16 }}>
            {[['📅 Date', fmtDate(ev.start_time)], ['📍 Location', ev.location], ['⏰ Ends', fmtDate(ev.end_time)], ['🎟 Capacity', `${ev.capacity} spots`]].map(([l,v]) => (
              <div key={l} style={{ background:'var(--panel)', borderRadius:7, padding:'9px 11px' }}>
                <div style={{ fontSize:10, fontFamily:fm, color:'var(--muted)', marginBottom:3 }}>{l}</div>
                <div style={{ fontSize:13, fontWeight:600 }}>{v}</div>
              </div>
            ))}
          </div>

          {ev.description && <p style={{ fontSize:13, color:'var(--muted)', lineHeight:1.6, marginBottom:16 }}>{ev.description}</p>}

          <div style={{ fontSize:10, fontFamily:fm, color:'var(--muted)', letterSpacing:1, textTransform:'uppercase', marginBottom:8 }}>SELECT TICKET TYPE</div>
          {ttBusy ? <div style={{ textAlign:'center', padding:20 }}><Spinner /></div>
          : tts.length === 0 ? <p style={{ fontSize:13, color:'var(--muted)', marginBottom:14 }}>No ticket types available yet.</p>
          : <div style={{ marginBottom:14 }}>{tts.map(tt => <TTCard key={tt.ticket_type_id} tt={tt} selected={selTT?.ticket_type_id===tt.ticket_type_id} onClick={() => { setSelTT(tt); setBuyErr('') }} />)}</div>
          }

          <div style={{ display:'flex', alignItems:'flex-end', gap:12 }}>
            <Field label="Qty" style={{ margin:0, flex:'0 0 76px' }}>
              <input type="number" min={1} max={selTT?.max_per_order||10} value={qty} onChange={e=>setQty(Math.max(1, parseInt(e.target.value)||1))} style={{ textAlign:'center' }} />
            </Field>
            <div style={{ flex:1 }}>
              <div style={{ fontSize:10, fontFamily:fm, color:'var(--muted)', marginBottom:5 }}>TOTAL</div>
              <div style={{ fontFamily:fd, fontSize:26, color:'var(--lime)', letterSpacing:1 }}>{total$}</div>
            </div>
            <Btn variant="primary" size="lg" onClick={purchase} loading={buying}>Buy Ticket</Btn>
          </div>
          {buyErr && <Alert type="error" style={{ marginTop:10 }}>{buyErr}</Alert>}
        </>}
      </Modal>

      {/* Payment confirmation */}
      <Modal open={payOpen} onClose={() => setPayOpen(false)} title="CONFIRM PAYMENT">
        {intent && <>
          <div style={{ background:'var(--panel)', borderRadius:8, padding:14, marginBottom:18, fontSize:13, lineHeight:1.9 }}>
            <div><strong>Event:</strong> {intent.event_name}</div>
            <div><strong>Quantity:</strong> {intent.quantity}</div>
            <div><strong>Amount:</strong> <span style={{ fontFamily:fd, fontSize:22, color:'var(--lime)' }}>{fmtCents(intent.amount)}</span></div>
            <div style={{ marginTop:8, fontSize:11, fontFamily:fm, color:'var(--muted)' }}>Intent: {intent.payment_intent_id}</div>
          </div>
          <Alert type="info" style={{ marginBottom:16 }}>
            Test mode. After clicking, run:<br/>
            <code style={{ background:'rgba(0,0,0,.3)', padding:'2px 6px', borderRadius:4, display:'inline-block', marginTop:6, fontSize:11 }}>
              docker compose exec stripe-cli stripe trigger payment_intent.succeeded
            </code>
          </Alert>
          <Btn variant="primary" style={{ width:'100%', justifyContent:'center', padding:13 }} onClick={() => { setPayOpen(false); toast('Intent created! Trigger the webhook to complete.','info') }}>
            Confirm (Test Mode)
          </Btn>
        </>}
      </Modal>
    </div>
  )
}

function EvCard({ ev, idx, onClick }) {
  const [h, setH] = useState(false)
  return (
    <div onClick={onClick} onMouseEnter={()=>setH(true)} onMouseLeave={()=>setH(false)}
      style={{ background:'var(--surf)', border:`1px solid ${h?'var(--lime)':'var(--border)'}`, borderRadius:10, overflow:'hidden', cursor:'pointer', transition:'all .2s', transform:h?'translateY(-4px)':'none', animation:`fadeUp .35s ease ${idx*.04}s both` }}>
      <div style={{ height:180, background:'var(--panel)', display:'flex', alignItems:'center', justifyContent:'center', fontSize:52, position:'relative', overflow:'hidden' }}>
        {EVT_EMOJI(ev.name)}
        {ev.image_url && <img src={ev.image_url} alt="" style={{ position:'absolute', inset:0, width:'100%', height:'100%', objectFit:'cover' }} onError={e=>e.target.style.display='none'} />}
      </div>
      <div style={{ padding:14 }}>
        <div style={{ fontSize:10, fontFamily:fm, color:'var(--lime)', letterSpacing:1, textTransform:'uppercase', marginBottom:5 }}>{fmtDate(ev.start_time)}</div>
        <div style={{ fontFamily:fd, fontSize:19, letterSpacing:1, marginBottom:4, lineHeight:1.1 }}>{ev.name}</div>
        <div style={{ fontSize:12, color:'var(--muted)', marginBottom:12 }}>📍 {ev.location}</div>
        <div style={{ display:'flex', alignItems:'center', justifyContent:'space-between', paddingTop:10, borderTop:'1px solid var(--border)' }}>
          <span style={{ fontSize:12, fontFamily:fm, color:'var(--amber)' }}>From {fmtCents(ev.min_price_cents||5000)}</span>
          <Badge color={ev.capacity>100?'lime':ev.capacity>10?'amber':'red'}>{ev.capacity} left</Badge>
        </div>
      </div>
    </div>
  )
}

function TTCard({ tt, selected, onClick }) {
  return (
    <div onClick={onClick} style={{ display:'flex', alignItems:'center', gap:10, background:selected?'rgba(212,255,0,.04)':'var(--panel)', border:`1px solid ${selected?'var(--lime)':'var(--border)'}`, borderRadius:8, padding:11, marginBottom:7, cursor:'pointer', transition:'all .15s' }}>
      <div style={{ width:17, height:17, borderRadius:'50%', border:`2px solid ${selected?'var(--lime)':'var(--border)'}`, display:'flex', alignItems:'center', justifyContent:'center', flexShrink:0, transition:'border-color .15s' }}>
        {selected && <div style={{ width:7, height:7, borderRadius:'50%', background:'var(--lime)' }} />}
      </div>
      <div style={{ flex:1 }}>
        <div style={{ display:'flex', alignItems:'center', gap:7, marginBottom:2 }}>
          <span style={{ fontSize:13, fontWeight:700 }}>{tt.name}</span>
          <Badge color={CAT_COLOR[tt.category]||'muted'}>{tt.category}</Badge>
        </div>
        <div style={{ fontSize:11, color:'var(--muted)', fontFamily:fm }}>
          {tt.available} available · max {tt.max_per_order}/order
          {!tt.on_sale && <span style={{ color:'var(--red)', marginLeft:8 }}>· Not on sale</span>}
        </div>
      </div>
      <div style={{ fontFamily:fd, fontSize:19, color:'var(--lime)', letterSpacing:1 }}>{tt.price_formatted}</div>
    </div>
  )
}

function CardSkeleton() {
  return (
    <div style={{ background:'var(--surf)', border:'1px solid var(--border)', borderRadius:10, overflow:'hidden' }}>
      <div style={{ height:180, background:'var(--panel)', animation:'pulse 1.5s ease infinite' }} />
      <div style={{ padding:14, display:'flex', flexDirection:'column', gap:8 }}>
        <Skeleton w="40%" h={10} /><Skeleton w="80%" h={16} /><Skeleton w="55%" h={11} />
      </div>
    </div>
  )
}
