// src/pages/AdminPage.jsx
import { useState, useEffect } from 'react'
import { api } from '../api/client'
import { useAuth } from '../context/AuthContext'
import { useToast } from '../context/ToastContext'
import { Btn, Badge, Modal, Field, Alert, Empty, Spinner, SectionHeader, Tabs, Avatar } from '../components/ui'

const fm = "'JetBrains Mono', monospace"

export default function AdminPage() {
  const { admin } = useAuth()
  const toast = useToast()
  const [tab,      setTab]      = useState('users')
  const [users,    setUsers]    = useState([])
  const [merchants,setMerchants]= useState([])
  const [loading,  setLoading]  = useState(false)

  // Create user
  const [uOpen, setUOpen] = useState(false)
  const [uF,    setUF]    = useState({ email:'', username:'', password:'' })
  const [uErr,  setUErr]  = useState('')
  const [uLoad, setULoad] = useState(false)

  // Create merchant
  const [mOpen, setMOpen] = useState(false)
  const [mF,    setMF]    = useState({ name:'', email:'', password:'', phone:'', description:'' })
  const [mErr,  setMErr]  = useState('')
  const [mLoad, setMLoad] = useState(false)

  useEffect(() => { load() }, [tab])

  async function load() {
    setLoading(true)
    if (tab === 'users') {
      const { ok, data } = await api.adminGetUsers('page=1&limit=100')
      setLoading(false)
      if (ok) setUsers(data.data||[])
    } else {
      const { ok, data } = await api.adminGetMerchants('page=1&limit=100')
      setLoading(false)
      if (ok) setMerchants(data.data||[])
    }
  }

  async function toggleUser(id, active) {
    const { ok, data } = await (active ? api.adminDeactivateUser(id) : api.adminActivateUser(id))
    if (ok) { toast(active?'User deactivated':'User activated','success'); load() }
    else toast(data.error||'Failed','error')
  }

  async function toggleMerchant(id, active) {
    const { ok, data } = await (active ? api.adminDeactivateMerchant(id) : api.adminActivateMerchant(id))
    if (ok) { toast(active?'Merchant deactivated':'Merchant activated','success'); load() }
    else toast(data.error||'Failed','error')
  }

  async function createUser() {
    setUErr(''); setULoad(true)
    const { ok, data } = await api.adminCreateUser(uF)
    setULoad(false)
    if (ok) { toast('User created!','success'); setUOpen(false); setUF({ email:'',username:'',password:'' }); load() }
    else setUErr(data.error||(data.errors?Object.values(data.errors).join(' · '):'Failed'))
  }

  async function createMerchant() {
    setMErr(''); setMLoad(true)
    const { ok, data } = await api.adminCreateMerchant(mF)
    setMLoad(false)
    if (ok) { toast('Merchant created!','success'); setMOpen(false); setMF({ name:'',email:'',password:'',phone:'',description:'' }); load() }
    else setMErr(data.error||'Failed')
  }

  const su = k => e => setUF(p=>({...p,[k]:e.target.value}))
  const sm = k => e => setMF(p=>({...p,[k]:e.target.value}))

  return (
    <div>
      <SectionHeader title="ADMIN PANEL" sub={`Signed in as ${admin?.email}`} />

      <Tabs active={tab} onChange={setTab} tabs={[{ value:'users', label:'👤 Users' }, { value:'merchants', label:'🏪 Merchants' }]} />

      {tab === 'users' && (
        <>
          <div style={{ marginBottom:14 }}><Btn variant="primary" size="sm" onClick={()=>setUOpen(true)}>＋ Create User</Btn></div>
          {loading
            ? <div style={{ textAlign:'center', padding:48 }}><Spinner /></div>
            : users.length === 0 ? <Empty icon="👤" title="NO USERS" />
            : users.map((u,i) => (
              <div key={u.user_id} style={{ display:'flex', alignItems:'center', gap:11, padding:'11px 13px', background:'var(--surf)', border:'1px solid var(--border)', borderRadius:8, marginBottom:8, animation:`fadeUp .3s ease ${i*.03}s both` }}>
                <Avatar name={u.username} size={36} />
                <div style={{ flex:1, minWidth:0 }}>
                  <div style={{ display:'flex', alignItems:'center', gap:7 }}>
                    <span style={{ fontSize:14, fontWeight:700 }}>{u.username}</span>
                    <Badge color={u.active?'lime':'red'}>{u.active?'active':'inactive'}</Badge>
                  </div>
                  <div style={{ fontSize:11, fontFamily:fm, color:'var(--muted)' }}>{u.email}</div>
                </div>
                <Btn variant={u.active?'danger':'success'} size="sm" onClick={()=>toggleUser(u.user_id, u.active)}>
                  {u.active?'Deactivate':'Activate'}
                </Btn>
              </div>
            ))
          }
        </>
      )}

      {tab === 'merchants' && (
        <>
          <div style={{ marginBottom:14 }}><Btn variant="primary" size="sm" onClick={()=>setMOpen(true)}>＋ Create Merchant</Btn></div>
          {loading
            ? <div style={{ textAlign:'center', padding:48 }}><Spinner /></div>
            : merchants.length === 0 ? <Empty icon="🏪" title="NO MERCHANTS" />
            : merchants.map((m,i) => (
              <div key={m.merchant_id} style={{ display:'flex', alignItems:'center', gap:11, padding:'11px 13px', background:'var(--surf)', border:'1px solid var(--border)', borderRadius:8, marginBottom:8, animation:`fadeUp .3s ease ${i*.03}s both` }}>
                <Avatar name={m.name} size={36} />
                <div style={{ flex:1, minWidth:0 }}>
                  <div style={{ display:'flex', alignItems:'center', gap:7 }}>
                    <span style={{ fontSize:14, fontWeight:700 }}>{m.name}</span>
                    <Badge color={m.active?'amber':'red'}>{m.active?'active':'inactive'}</Badge>
                  </div>
                  <div style={{ fontSize:11, fontFamily:fm, color:'var(--muted)' }}>{m.email}</div>
                </div>
                <Btn variant={m.active?'danger':'success'} size="sm" onClick={()=>toggleMerchant(m.merchant_id, m.active)}>
                  {m.active?'Deactivate':'Activate'}
                </Btn>
              </div>
            ))
          }
        </>
      )}

      {/* Create User */}
      <Modal open={uOpen} onClose={()=>setUOpen(false)} title="CREATE USER">
        {uErr && <Alert type="error" style={{ marginBottom:12 }}>{uErr}</Alert>}
        <Field label="Email"><input type="email" value={uF.email} onChange={su('email')} placeholder="user@example.com" /></Field>
        <Field label="Username"><input value={uF.username} onChange={su('username')} placeholder="handle" /></Field>
        <Field label="Password"><input type="password" value={uF.password} onChange={su('password')} placeholder="PassW0rd!" /></Field>
        <Btn variant="primary" style={{ width:'100%', justifyContent:'center' }} onClick={createUser} loading={uLoad}>Create User</Btn>
      </Modal>

      {/* Create Merchant */}
      <Modal open={mOpen} onClose={()=>setMOpen(false)} title="CREATE MERCHANT">
        {mErr && <Alert type="error" style={{ marginBottom:12 }}>{mErr}</Alert>}
        <Field label="Name"><input value={mF.name} onChange={sm('name')} placeholder="Produtora XYZ" /></Field>
        <Field label="Email"><input type="email" value={mF.email} onChange={sm('email')} placeholder="contact@xyz.com" /></Field>
        <Field label="Password"><input type="password" value={mF.password} onChange={sm('password')} placeholder="PassW0rd!" /></Field>
        <Field label="Phone"><input value={mF.phone} onChange={sm('phone')} placeholder="44999000000" /></Field>
        <Field label="Description"><textarea value={mF.description} onChange={sm('description')} rows={2} placeholder="About…" /></Field>
        <Btn variant="primary" style={{ width:'100%', justifyContent:'center' }} onClick={createMerchant} loading={mLoad}>Create Merchant</Btn>
      </Modal>
    </div>
  )
}
