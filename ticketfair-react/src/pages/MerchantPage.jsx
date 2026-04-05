// src/pages/MerchantPage.jsx
import { useState, useEffect } from 'react'
import { api } from '../api/client'
import { useAuth } from '../context/AuthContext'
import { useToast } from '../context/ToastContext'
import { Btn, Badge, Modal, Field, Alert, Empty, Spinner, SectionHeader, PanelSection, Row, Divider } from '../components/ui'

const fd = "'Bebas Neue', Impact, sans-serif"
const fm = "'JetBrains Mono', monospace"
const fmtDate = iso => iso ? new Date(iso).toLocaleDateString('en-US',{month:'short',day:'numeric',year:'numeric'}) : '—'

const CATS = [
  ['general','General Admission'],['vip','VIP'],['early_bird','Early Bird'],
  ['reserved','Reserved Seating'],['group','Group Pack'],['day_pass','Day Pass'],
  ['tiered','Tiered / Release'],['complimentary','Complimentary / Free'],['demographic','Student / Senior'],
]
const CAT_COL = { vip:'amber', early_bird:'lime', complimentary:'green', reserved:'blue', demographic:'purple', group:'blue', day_pass:'amber', tiered:'purple' }

export default function MerchantPage() {
  const { merchant } = useAuth()
  const toast = useToast()
  const [events,   setEvents]   = useState([])
  const [loading,  setLoading]  = useState(true)

  // Create event modal
  const [ceOpen, setCeOpen] = useState(false)
  const [ceF, setCeF]       = useState({ name:'', description:'', location:'', start_time:'', end_time:'', capacity:'' })
  const [ceErr, setCeErr]   = useState('')
  const [ceLoad, setCeLoad] = useState(false)

  // Ticket types modal
  const [ttOpen,    setTtOpen]    = useState(false)
  const [ttEventID, setTtEventID] = useState(null)
  const [ttName,    setTtName]    = useState('')
  const [tts,       setTts]       = useState([])
  const [ttLoad,    setTtLoad]    = useState(false)
  const [ttF, setTtF]             = useState({ name:'', category:'general', price_cents:'', capacity:'', min_per_order:'1', max_per_order:'10' })
  const [ttErr,     setTtErr]     = useState('')
  const [ttAdd,     setTtAdd]     = useState(false)

  // Validate
  const [valID,  setValID]  = useState('')
  const [valRes, setValRes] = useState(null)
  const [valLoad,setValLoad]= useState(false)

  useEffect(() => { loadEvents() }, [merchant])

  async function loadEvents() {
    if (!merchant) return
    setLoading(true)
    const { ok, data } = await api.getEvents('page=1&limit=100')
    setLoading(false)
    if (ok) setEvents((data.data||[]).filter(e => e.merchant_id === merchant.merchant_id))
  }

  async function createEvent() {
    const { name, description, location, start_time, end_time, capacity } = ceF
    if (!name||!location||!start_time||!end_time||!capacity) { setCeErr('All fields required'); return }
    setCeLoad(true); setCeErr('')
    const { ok, data } = await api.createEvent({ name, description, location, capacity:parseInt(capacity), start_time: new Date(start_time).toISOString(), end_time: new Date(end_time).toISOString() })
    setCeLoad(false)
    if (ok) { toast('Event created!','success'); setCeOpen(false); setCeF({ name:'',description:'',location:'',start_time:'',end_time:'',capacity:'' }); loadEvents() }
    else setCeErr(data.error||'Failed')
  }

  async function openTT(id, name) {
    setTtEventID(id); setTtName(name); setTtErr(''); setTtOpen(true); await loadTTs(id)
  }

  async function loadTTs(id) {
    setTtLoad(true)
    const { ok, data } = await api.getMerchantTT(id)
    setTtLoad(false)
    if (ok) setTts(data.data||[])
  }

  async function addTT() {
    const { name, category, price_cents, capacity, min_per_order, max_per_order } = ttF
    if (!name||!capacity) { setTtErr('Name and capacity required'); return }
    setTtAdd(true); setTtErr('')
    const { ok, data } = await api.createTT(ttEventID, { name, category, price_cents:parseInt(price_cents)||0, capacity:parseInt(capacity), min_per_order:parseInt(min_per_order)||1, max_per_order:parseInt(max_per_order)||10 })
    setTtAdd(false)
    if (ok) { toast('Ticket type added!','success'); setTtF({ name:'',category:'general',price_cents:'',capacity:'',min_per_order:'1',max_per_order:'10' }); loadTTs(ttEventID) }
    else setTtErr(data.error||'Failed')
  }

  async function deleteTT(ttID) {
    if (!confirm('Delete this ticket type?')) return
    const { ok, data } = await api.deleteTT(ttEventID, ttID)
    if (ok) { toast('Deleted','success'); loadTTs(ttEventID) }
    else toast(data.error||'Cannot delete — tickets may have been sold','error')
  }

  async function validate() {
    if (!valID.trim()) return
    setValLoad(true); setValRes(null)
    const { ok, data } = await api.validateTicket(valID.trim())
    setValLoad(false)
    setValRes({ ok, msg: ok ? `✓ ${data.data?.ticket_type_name||'Ticket'} — ${data.data?.status}` : `✕ ${data.error||'Invalid ticket'}` })
    if (ok) { toast('Ticket validated! ✓','success'); setValID('') }
  }

  const ce = k => e => setCeF(p => ({...p,[k]:e.target.value}))
  const tt = k => e => setTtF(p => ({...p,[k]:e.target.value}))

  return (
    <div>
      <SectionHeader title="MY VENUE" sub={merchant?.name} />

      <PanelSection title="MY EVENTS" action={<Btn variant="primary" size="sm" onClick={()=>setCeOpen(true)}>＋ New Event</Btn>}>
        {loading
          ? <div style={{ textAlign:'center', padding:24 }}><Spinner /></div>
          : events.length === 0
            ? <Empty icon="🎪" title="NO EVENTS" sub="Create your first event" />
            : events.map((ev,i) => (
              <Row key={ev.event_id} icon="🎪"
                title={ev.name}
                sub={`${fmtDate(ev.start_time)} · ${ev.location} · ${ev.capacity} spots`}
                right={<>
                  <Btn variant="ghost" size="sm" onClick={()=>openTT(ev.event_id,ev.name)}>🎫 Tickets</Btn>
                  <Badge color={ev.active?'lime':'muted'}>{ev.active?'active':'inactive'}</Badge>
                </>}
                style={{ animation:`fadeUp .3s ease ${i*.04}s both` }}
              />
            ))
        }
      </PanelSection>

      <PanelSection title="VALIDATE TICKET">
        <p style={{ fontSize:13, color:'var(--muted)', marginBottom:14 }}>Paste a ticket UUID to mark it as used at the door.</p>
        <div style={{ display:'flex', gap:8 }}>
          <input value={valID} onChange={e=>setValID(e.target.value)} placeholder="Paste ticket UUID…" onKeyDown={e=>e.key==='Enter'&&validate()} style={{ flex:1 }} />
          <Btn variant="primary" onClick={validate} loading={valLoad}>Validate</Btn>
        </div>
        {valRes && <Alert type={valRes.ok?'success':'error'} style={{ marginTop:11 }}>{valRes.msg}</Alert>}
      </PanelSection>

      {/* Create Event Modal */}
      <Modal open={ceOpen} onClose={()=>setCeOpen(false)} title="NEW EVENT" width={580}>
        {ceErr && <Alert type="error" style={{ marginBottom:14 }}>{ceErr}</Alert>}
        <div style={{ display:'grid', gridTemplateColumns:'1fr 1fr', gap:11 }}>
          <Field label="Event Name" style={{ gridColumn:'1/-1' }}><input value={ceF.name} onChange={ce('name')} placeholder="Festival Night" /></Field>
          <Field label="Location"   style={{ gridColumn:'1/-1' }}><input value={ceF.location} onChange={ce('location')} placeholder="Arena, City" /></Field>
        </div>
        <Field label="Description"><textarea value={ceF.description} onChange={ce('description')} rows={3} placeholder="Describe the event…" /></Field>
        <div style={{ display:'grid', gridTemplateColumns:'1fr 1fr', gap:11 }}>
          <Field label="Start Time"><input type="datetime-local" value={ceF.start_time} onChange={ce('start_time')} /></Field>
          <Field label="End Time"><input type="datetime-local" value={ceF.end_time} onChange={ce('end_time')} /></Field>
        </div>
        <Field label="Total Capacity"><input type="number" value={ceF.capacity} onChange={ce('capacity')} placeholder="500" /></Field>
        <Btn variant="primary" style={{ width:'100%', justifyContent:'center' }} onClick={createEvent} loading={ceLoad}>Create Event →</Btn>
      </Modal>

      {/* Ticket Types Modal */}
      <Modal open={ttOpen} onClose={()=>setTtOpen(false)} title={`TICKETS — ${ttName}`} width={600}>
        {ttLoad
          ? <div style={{ textAlign:'center', padding:24 }}><Spinner /></div>
          : tts.length === 0
            ? <p style={{ fontSize:13, color:'var(--muted)', marginBottom:18 }}>No ticket types yet.</p>
            : <div style={{ marginBottom:18 }}>
                {tts.map(tt => (
                  <div key={tt.ticket_type_id} style={{ display:'flex', alignItems:'center', gap:10, padding:'10px 11px', background:'var(--panel)', border:'1px solid var(--border)', borderRadius:8, marginBottom:7 }}>
                    <span style={{ fontSize:17 }}>🎫</span>
                    <div style={{ flex:1 }}>
                      <div style={{ display:'flex', alignItems:'center', gap:7, marginBottom:2 }}>
                        <span style={{ fontSize:13, fontWeight:700 }}>{tt.name}</span>
                        <Badge color={CAT_COL[tt.category]||'muted'}>{tt.category}</Badge>
                        {!tt.active && <Badge color="red">inactive</Badge>}
                      </div>
                      <div style={{ fontSize:11, fontFamily:fm, color:'var(--muted)' }}>{tt.price_formatted} · {tt.available}/{tt.capacity} avail · Sold: {tt.sold}</div>
                    </div>
                    <Btn variant="danger" size="sm" onClick={()=>deleteTT(tt.ticket_type_id)}>Delete</Btn>
                  </div>
                ))}
              </div>
        }

        <Divider />
        <div style={{ fontFamily:fd, fontSize:15, letterSpacing:1, marginBottom:13, color:'var(--lime)' }}>ADD TICKET TYPE</div>
        {ttErr && <Alert type="error" style={{ marginBottom:11 }}>{ttErr}</Alert>}
        <div style={{ display:'grid', gridTemplateColumns:'1fr 1fr', gap:11 }}>
          <Field label="Name"><input value={ttF.name} onChange={tt('name')} placeholder="General Admission" /></Field>
          <Field label="Category">
            <select value={ttF.category} onChange={tt('category')}>
              {CATS.map(([v,l]) => <option key={v} value={v}>{l}</option>)}
            </select>
          </Field>
          <Field label="Price (cents)"><input type="number" value={ttF.price_cents} onChange={tt('price_cents')} placeholder="5000 = R$50" /></Field>
          <Field label="Capacity"><input type="number" value={ttF.capacity} onChange={tt('capacity')} placeholder="200" /></Field>
          <Field label="Min/order"><input type="number" value={ttF.min_per_order} onChange={tt('min_per_order')} /></Field>
          <Field label="Max/order"><input type="number" value={ttF.max_per_order} onChange={tt('max_per_order')} /></Field>
        </div>
        <Btn variant="primary" style={{ width:'100%', justifyContent:'center' }} onClick={addTT} loading={ttAdd}>Add Ticket Type</Btn>
      </Modal>
    </div>
  )
}
