// src/pages/ConnectionsPage.jsx
import { useState, useEffect } from 'react'
import { api } from '../api/client'
import { useToast } from '../context/ToastContext'
import { Btn, Badge, Modal, Field, Alert, Empty, Spinner, SectionHeader, Tabs, Avatar } from '../components/ui'

const fm = "'JetBrains Mono', monospace"

export default function ConnectionsPage({ onPendingChange }) {
  const toast = useToast()
  const [tab,      setTab]      = useState('accepted')
  const [accepted, setAccepted] = useState([])
  const [pending,  setPending]  = useState([])
  const [loading,  setLoading]  = useState(true)
  const [addOpen,  setAddOpen]  = useState(false)
  const [addID,    setAddID]    = useState('')
  const [addErr,   setAddErr]   = useState('')
  const [addOk,    setAddOk]    = useState('')

  useEffect(() => { loadAll() }, [])

  async function loadAll() {
    setLoading(true)
    const [aRes, pRes] = await Promise.all([
      api.getConnections('status=accepted&page=1&limit=50'),
      api.getPendingRequests(),
    ])
    setLoading(false)
    if (aRes.ok) setAccepted(aRes.data.data || [])
    if (pRes.ok) {
      const p = pRes.data.data || []
      setPending(p)
      onPendingChange?.(p.length)
    }
  }

  async function respond(id, action) {
    const { ok, data } = await api.respondConnection(id, { action })
    if (ok) { toast(action==='accept'?'Connection accepted!':'Request declined','success'); loadAll() }
    else toast(data.error||'Failed','error')
  }

  async function remove(id) {
    if (!confirm('Remove this connection?')) return
    const { ok } = await api.removeConnection(id)
    if (ok) { toast('Connection removed','success'); loadAll() }
    else toast('Failed','error')
  }

  async function sendReq() {
    setAddErr(''); setAddOk('')
    if (!addID.trim()) { setAddErr('Enter a user ID'); return }
    const { ok, data } = await api.sendConnectionReq({ addressee_id: addID.trim() })
    if (ok) { setAddOk('Request sent!'); setAddID(''); setTimeout(()=>{ setAddOpen(false); setAddOk('') }, 1500) }
    else setAddErr(data.error||'Failed')
  }

  const list = tab === 'pending' ? pending : accepted

  return (
    <div>
      <SectionHeader title="CONNECTIONS" sub="Your network" action={<Btn variant="primary" size="sm" onClick={()=>setAddOpen(true)}>＋ Add</Btn>} />

      <Tabs active={tab} onChange={setTab} tabs={[{ value:'accepted', label:'Connected' }, { value:'pending', label:'Pending', count: pending.length }]} />

      {loading
        ? <div style={{ textAlign:'center', padding:48 }}><Spinner /></div>
        : list.length === 0
          ? <Empty icon="🔗" title={tab==='pending'?'NO REQUESTS':'NO CONNECTIONS'} sub={tab==='pending'?'No pending requests yet':'Add someone by their user ID'} />
          : list.map((c, i) => (
            <ConnCard key={c.connection_id} conn={c} isPending={tab==='pending'} idx={i}
              onAccept={()=>respond(c.connection_id,'accept')}
              onDecline={()=>respond(c.connection_id,'decline')}
              onRemove={()=>remove(c.connection_id)}
            />
          ))
      }

      <Modal open={addOpen} onClose={()=>setAddOpen(false)} title="ADD CONNECTION">
        <p style={{ fontSize:13, color:'var(--muted)', marginBottom:14 }}>Paste a user's UUID to send them a connection request.</p>
        {addErr && <Alert type="error"   style={{ marginBottom:11 }}>{addErr}</Alert>}
        {addOk  && <Alert type="success" style={{ marginBottom:11 }}>{addOk}</Alert>}
        <Field label="User ID"><input value={addID} onChange={e=>setAddID(e.target.value)} placeholder="Paste user UUID…" /></Field>
        <Btn variant="primary" style={{ width:'100%', justifyContent:'center' }} onClick={sendReq}>Send Request</Btn>
      </Modal>
    </div>
  )
}

function ConnCard({ conn: c, isPending, idx, onAccept, onDecline, onRemove }) {
  return (
    <div style={{ display:'flex', alignItems:'center', gap:11, padding:'11px 13px', background:'var(--surf)', border:'1px solid var(--border)', borderRadius:8, marginBottom:8, animation:`fadeUp .3s ease ${idx*.04}s both` }}>
      <Avatar name={c.user.username} size={38} />
      <div style={{ flex:1, minWidth:0 }}>
        <div style={{ fontSize:14, fontWeight:700 }}>{c.user.username}</div>
        <div style={{ fontSize:11, fontFamily:fm, color:'var(--muted)', overflow:'hidden', textOverflow:'ellipsis', whiteSpace:'nowrap' }}>{c.user.email}</div>
      </div>
      <div style={{ display:'flex', gap:6, flexShrink:0 }}>
        {isPending
          ? <><Btn variant="success" size="sm" onClick={onAccept}>Accept</Btn><Btn variant="danger" size="sm" onClick={onDecline}>Decline</Btn></>
          : <Btn variant="ghost" size="sm" onClick={onRemove}>Remove</Btn>
        }
      </div>
    </div>
  )
}
