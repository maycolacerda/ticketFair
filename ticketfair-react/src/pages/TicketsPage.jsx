// src/pages/TicketsPage.jsx
import { useState, useEffect } from 'react'
import { api } from '../api/client'
import { useToast } from '../context/ToastContext'
import { Btn, Badge, Modal, Field, Alert, Empty, Spinner, SectionHeader, Row } from '../components/ui'

const fm = "'JetBrains Mono', monospace"
const fd = "'Bebas Neue', Impact, sans-serif"
const fmtDate  = iso => iso ? new Date(iso).toLocaleDateString('en-US',{month:'short',day:'numeric',year:'numeric'}) : '—'
const fmtCents = c   => c === 0 ? 'Free' : `R$ ${(c/100).toFixed(2).replace('.',',')}`
const STATUS_COLOR = { active:'lime', used:'muted', refunded:'red', cancelled:'red' }

export default function TicketsPage() {
  const toast = useToast()
  const [tickets,  setTickets]  = useState([])
  const [loading,  setLoading]  = useState(true)
  const [giftOpen, setGiftOpen] = useState(false)
  const [giftID,   setGiftID]   = useState(null)
  const [conns,    setConns]    = useState([])
  const [recip,    setRecip]    = useState('')
  const [msg,      setMsg]      = useState('')
  const [giftErr,  setGiftErr]  = useState('')
  const [gifting,  setGifting]  = useState(false)

  useEffect(() => { load() }, [])

  async function load() {
    setLoading(true)
    const { ok, data } = await api.getTickets('page=1&limit=50')
    setLoading(false)
    if (ok) setTickets(data.data || [])
  }

  async function refund(txID) {
    if (!confirm('Refund this ticket?')) return
    const { ok, data } = await api.refundTicket({ transaction_id: txID })
    if (ok) { toast('Ticket refunded','success'); load() }
    else toast(data.error||'Refund failed','error')
  }

  async function openGift(ticketID) {
    setGiftID(ticketID); setGiftErr(''); setMsg(''); setRecip('')
    setGiftOpen(true)
    const { ok, data } = await api.getConnections('status=accepted&page=1&limit=50')
    if (ok) setConns(data.data || [])
  }

  async function submitGift() {
    if (!recip) { setGiftErr('Select a recipient'); return }
    setGifting(true)
    const { ok, data } = await api.giftTicket(giftID, { recipient_id: recip, message: msg })
    setGifting(false)
    if (ok) { toast(`Ticket gifted to ${data.data.recipient}! 🎁`,'success'); setGiftOpen(false); load() }
    else setGiftErr(data.error || 'Gift failed')
  }

  if (loading) return <div style={{ textAlign:'center', padding:48 }}><Spinner /></div>

  return (
    <div>
      <SectionHeader title="MY TICKETS" sub={`${tickets.filter(t=>t.status==='active').length} active`} />

      {tickets.length === 0
        ? <Empty icon="🎟" title="NO TICKETS" sub="Browse events to purchase your first ticket" />
        : tickets.map((t, i) => (
          <TicketCard key={t.ticket_id} ticket={t} idx={i} onRefund={refund} onGift={openGift} />
        ))
      }

      <Modal open={giftOpen} onClose={() => setGiftOpen(false)} title="GIFT TICKET">
        <p style={{ fontSize:13, color:'var(--muted)', marginBottom:14 }}>Only accepted connections can receive gifts. Ownership transfers immediately.</p>
        {giftErr && <Alert type="error" style={{ marginBottom:12 }}>{giftErr}</Alert>}
        <Field label="Recipient">
          <select value={recip} onChange={e=>setRecip(e.target.value)}>
            <option value="">Select a connection…</option>
            {conns.map(c => <option key={c.user.user_id} value={c.user.user_id}>{c.user.username} ({c.user.email})</option>)}
          </select>
        </Field>
        <Field label="Message (optional)">
          <textarea value={msg} onChange={e=>setMsg(e.target.value)} placeholder="Enjoy the show! 🎸" rows={3} />
        </Field>
        <Btn variant="primary" style={{ width:'100%', justifyContent:'center' }} onClick={submitGift} loading={gifting}>Send Gift 🎁</Btn>
      </Modal>
    </div>
  )
}

function TicketCard({ ticket: t, idx, onRefund, onGift }) {
  return (
    <div style={{ background:'var(--surf)', border:'1px solid var(--border)', borderRadius:10, overflow:'hidden', marginBottom:11, animation:`fadeUp .3s ease ${idx*.04}s both` }}>
      <div style={{ display:'flex', alignItems:'center', justifyContent:'space-between', padding:'13px 15px', borderBottom:'1px dashed var(--border)' }}>
        <div>
          <div style={{ fontFamily:fd, fontSize:15, letterSpacing:1 }}>{t.event?.name||'Event'}</div>
          <div style={{ fontSize:11, color:'var(--muted)', fontFamily:fm }}>{t.ticket_type_name||'General'} · {fmtDate(t.event?.start_time)}</div>
        </div>
        <Badge color={STATUS_COLOR[t.status]||'muted'}>{t.status}</Badge>
      </div>
      <div style={{ display:'flex', alignItems:'flex-start', gap:14, padding:15 }}>
        <div style={{ width:60, height:60, flexShrink:0, background:'var(--panel)', borderRadius:7, display:'flex', alignItems:'center', justifyContent:'center', fontSize:30 }}>
          {t.is_gift ? '🎁' : '🎟'}
        </div>
        <div style={{ flex:1 }}>
          <div style={{ fontSize:11, fontFamily:fm, color:'var(--muted)', marginBottom:3 }}>ID: {t.ticket_id.slice(0,20)}…</div>
          {t.is_gift && <div style={{ fontSize:11, color:'var(--blue)', marginBottom:3 }}>🎁 Gifted{t.gift_message?` — "${t.gift_message}"`:''}</div>}
          <div style={{ fontFamily:fm, fontSize:13, color:'var(--lime)' }}>{fmtCents(t.price_paid_cents)}</div>
        </div>
        <div style={{ display:'flex', gap:6, flexShrink:0 }}>
          {t.status==='active' && !t.is_gift && <Btn variant="blue" size="sm" onClick={()=>onGift(t.ticket_id)}>🎁 Gift</Btn>}
          {t.status==='active' && <Btn variant="danger" size="sm" onClick={()=>onRefund(t.transaction_id)}>Refund</Btn>}
        </div>
      </div>
    </div>
  )
}


// ─────────────────────────────────────────────────────
// PaymentsPage
// ─────────────────────────────────────────────────────
export function PaymentsPage() {
  const toast = useToast()
  const [payments, setPayments] = useState([])
  const [loading,  setLoading]  = useState(true)

  useEffect(() => {
    api.getPayments('page=1&limit=50').then(({ ok, data }) => {
      setLoading(false)
      if (ok) setPayments(data.data || [])
    })
  }, [])

  async function refund(id) {
    if (!confirm('Issue a Stripe refund?')) return
    const { ok, data } = await api.refundPayment(id)
    if (ok) { toast('Refund issued','success'); setPayments(p => p.map(x => x.payment_id===id ? {...x, status:'refunded'} : x)) }
    else toast(data.error||'Refund failed','error')
  }

  const S = { succeeded:'green', pending:'amber', failed:'red', refunded:'muted', canceled:'muted' }

  if (loading) return <div style={{ textAlign:'center', padding:48 }}><Spinner /></div>

  return (
    <div>
      <SectionHeader title="PAYMENTS" sub="Stripe payment history" />
      {payments.length === 0
        ? <Empty icon="💳" title="NO PAYMENTS" sub="Your Stripe history will appear here" />
        : payments.map((p, i) => (
          <Row key={p.payment_id} icon="💳"
            title={`${p.currency.toUpperCase()} ${(p.amount/100).toFixed(2)}`}
            sub={`${p.stripe_payment_id} · ${new Date(p.created_at).toLocaleDateString()}`}
            right={<>
              <Badge color={S[p.status]||'muted'}>{p.status}</Badge>
              {p.status==='succeeded' && <Btn variant="danger" size="sm" onClick={()=>refund(p.payment_id)}>Refund</Btn>}
            </>}
            style={{ animation:`fadeUp .3s ease ${i*.04}s both` }}
          />
        ))
      }
    </div>
  )
}


// ─────────────────────────────────────────────────────
// FeedPage
// ─────────────────────────────────────────────────────
export function FeedPage() {
  const [events,  setEvents]  = useState([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    api.getConnectionEvents('page=1&limit=20').then(({ ok, data }) => {
      setLoading(false)
      if (ok) setEvents(data.data || [])
    })
  }, [])

  if (loading) return <div style={{ textAlign:'center', padding:48 }}><Spinner /></div>

  return (
    <div>
      <SectionHeader title="FRIENDS FEED" sub="Events your connections are attending" />
      {events.length === 0
        ? <Empty icon="🎪" title="QUIET HERE" sub="Add connections and see which events they're going to" />
        : events.map((ev, i) => (
          <div key={ev.event_id} style={{ borderLeft:'3px solid var(--lime)', padding:'13px 15px', marginBottom:11, background:'var(--surf)', borderRadius:'0 8px 8px 0', animation:`fadeUp .3s ease ${i*.05}s both` }}>
            <div style={{ fontFamily:fd, fontSize:17, letterSpacing:1, marginBottom:3 }}>{ev.name}</div>
            <div style={{ fontSize:11, fontFamily:fm, color:'var(--muted)', marginBottom:11 }}>
              {new Date(ev.start_time).toLocaleDateString('en-US',{month:'short',day:'numeric',year:'numeric'})} · {ev.location}
            </div>
            <div style={{ display:'flex', gap:7, flexWrap:'wrap' }}>
              {ev.attendees.map(a => (
                <div key={a.user_id} style={{ display:'flex', alignItems:'center', gap:6, padding:'4px 10px', background:'var(--panel)', borderRadius:99, fontSize:12, border:'1px solid var(--border)' }}>
                  <div style={{ width:18, height:18, borderRadius:'50%', background:'var(--faint)', display:'flex', alignItems:'center', justifyContent:'center', fontSize:9, fontFamily:fd }}>
                    {a.username[0].toUpperCase()}
                  </div>
                  {a.username}
                </div>
              ))}
              {ev.attendees_count > ev.attendees.length && (
                <div style={{ padding:'4px 10px', background:'var(--panel)', borderRadius:99, fontSize:12, border:'1px solid var(--border)', color:'var(--muted)' }}>
                  +{ev.attendees_count - ev.attendees.length} more
                </div>
              )}
            </div>
          </div>
        ))
      }
    </div>
  )
}
