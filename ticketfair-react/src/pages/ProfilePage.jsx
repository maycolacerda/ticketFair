// src/pages/ProfilePage.jsx
import { useState, useEffect } from 'react'
import { api } from '../api/client'
import { useAuth } from '../context/AuthContext'
import { useToast } from '../context/ToastContext'
import { Btn, Badge, Modal, Field, Alert, Spinner, SectionHeader, Stat, Divider, PanelSection } from '../components/ui'

const fm = "'JetBrains Mono', monospace"

export default function ProfilePage() {
  const { user } = useAuth()
  const toast = useToast()
  const [profile,  setProfile]  = useState(null)
  const [loading,  setLoading]  = useState(true)
  const [editOpen, setEditOpen] = useState(false)
  const [cpOpen,   setCpOpen]   = useState(false)
  const [tktCount, setTktCount] = useState('—')
  const [connCount,setConnCount]= useState('—')

  useEffect(() => { load() }, [])

  async function load() {
    setLoading(true)
    const [pRes, tRes, cRes] = await Promise.all([
      api.getProfile(),
      api.getTickets('page=1&limit=1'),
      api.getConnections('status=accepted&page=1&limit=1'),
    ])
    setLoading(false)
    if (pRes.ok) setProfile(pRes.data.data)
    if (tRes.ok) setTktCount(tRes.data.total ?? '—')
    if (cRes.ok) setConnCount(cRes.data.total ?? '—')
  }

  if (loading) return <div style={{ textAlign:'center', padding:48 }}><Spinner /></div>

  return (
    <div>
      <SectionHeader title="MY PROFILE" />

      <div style={{ display:'grid', gridTemplateColumns:'1fr 1fr', gap:14, marginBottom:22 }}>
        <Stat value={tktCount}  label="Active Tickets" />
        <Stat value={connCount} label="Connections" color="var(--blue)" />
      </div>

      <PanelSection title="PERSONAL INFO" action={<Btn variant="ghost" size="sm" onClick={()=>setEditOpen(true)}>{profile ? 'Edit' : 'Create Profile'}</Btn>}>
        {profile ? (
          <div style={{ display:'grid', gridTemplateColumns:'1fr 1fr', gap:12 }}>
            {[['First Name',profile.first_name],['Last Name',profile.last_name],['Phone',profile.phone_number],['City',profile.address?.city||'—'],['State',profile.address?.state||'—'],['Country',profile.address?.country||'—']].map(([l,v]) => (
              <div key={l}>
                <div style={{ fontSize:10, fontFamily:fm, color:'var(--muted)', marginBottom:3 }}>{l}</div>
                <div style={{ fontSize:14, fontWeight:600 }}>{v}</div>
              </div>
            ))}
          </div>
        ) : (
          <p style={{ fontSize:13, color:'var(--muted)' }}>No profile yet. Create one to enable verification.</p>
        )}
      </PanelSection>

      <PanelSection title="VERIFICATION">
        <div style={{ display:'grid', gridTemplateColumns:'1fr 1fr', gap:12 }}>
          <VerifyCard type="email" verified={profile?.verified_email} label="📧 Email" onDone={load} disabled={!profile} />
          <VerifyCard type="phone" verified={profile?.verified_phone} label="📱 Phone" onDone={load} disabled={!profile} />
        </div>
      </PanelSection>

      <PanelSection title="SECURITY">
        <Btn variant="ghost" size="sm" onClick={()=>setCpOpen(true)}>Change Password</Btn>
        <div style={{ marginTop:10, fontSize:11, fontFamily:fm, color:'var(--muted)' }}>User ID: {user?.user_id}</div>
      </PanelSection>

      <ProfileEditModal open={editOpen} onClose={()=>setEditOpen(false)} profile={profile} onSaved={()=>{ setEditOpen(false); load() }} />
      <ChangePasswordModal open={cpOpen} onClose={()=>setCpOpen(false)} email={user?.email} />
    </div>
  )
}

function VerifyCard({ type, verified, label, onDone, disabled }) {
  const toast = useToast()
  const [sent,    setSent]    = useState(false)
  const [code,    setCode]    = useState('')
  const [loading, setLoading] = useState(false)

  async function send() {
    setLoading(true)
    const { ok, data } = await (type==='email' ? api.sendEmailVerify() : api.sendPhoneVerify())
    setLoading(false)
    if (ok) { toast(`Code sent! ${data.data?.message||''}`, 'info'); setSent(true) }
    else toast(data.error||'Failed','error')
  }

  async function verify() {
    setLoading(true)
    const { ok, data } = await (type==='email' ? api.verifyEmail({code}) : api.verifyPhone({code}))
    setLoading(false)
    if (ok && data.data?.verified) { toast(`${type} verified! ✓`,'success'); onDone() }
    else toast(data.error||'Invalid code','error')
  }

  return (
    <div style={{ background:'var(--panel)', border:'1px solid var(--border)', borderRadius:8, padding:13 }}>
      <div style={{ display:'flex', alignItems:'center', gap:7, marginBottom:11 }}>
        <span style={{ fontWeight:700, fontSize:13 }}>{label}</span>
        <Badge color={verified?'green':'red'}>{verified?'Verified':'Unverified'}</Badge>
      </div>
      {!verified && !disabled && (
        sent
          ? <div style={{ display:'flex', gap:6 }}>
              <input value={code} onChange={e=>setCode(e.target.value)} placeholder="6-digit code" maxLength={6} />
              <Btn variant="primary" size="sm" onClick={verify} loading={loading}>Verify</Btn>
            </div>
          : <Btn variant="ghost" size="sm" onClick={send} loading={loading}>Send Code</Btn>
      )}
      {disabled && !verified && <p style={{ fontSize:11, color:'var(--muted)', fontFamily:"'JetBrains Mono',monospace" }}>Create profile first</p>}
    </div>
  )
}

function ProfileEditModal({ open, onClose, profile, onSaved }) {
  const toast = useToast()
  const [err,  setErr]  = useState('')
  const [busy, setBusy] = useState(false)
  const [f, setF] = useState({ first_name:'', last_name:'', phone_number:'', street:'', city:'', state:'', country:'BR', zip_code:'' })

  useEffect(() => {
    if (profile && open) setF({
      first_name: profile.first_name||'', last_name: profile.last_name||'',
      phone_number: profile.phone_number||'',
      street: profile.address?.street||'', city: profile.address?.city||'',
      state: profile.address?.state||'', country: profile.address?.country||'BR',
      zip_code: profile.address?.zip_code||'',
    })
  }, [profile, open])

  const set = k => e => setF(p => ({...p, [k]: e.target.value}))

  async function save() {
    setErr(''); setBusy(true)
    const body = { first_name:f.first_name, last_name:f.last_name, phone_number:f.phone_number, address:{ street:f.street, city:f.city, state:f.state, country:f.country.toUpperCase(), zip_code:f.zip_code } }
    const fn = profile ? api.updateProfile : api.createProfile
    const { ok, data } = await fn(body)
    setBusy(false)
    if (ok) { toast('Profile saved!','success'); onSaved() }
    else setErr(data.error||'Save failed')
  }

  return (
    <Modal open={open} onClose={onClose} title={profile?'EDIT PROFILE':'CREATE PROFILE'}>
      {err && <Alert type="error" style={{ marginBottom:14 }}>{err}</Alert>}
      <div style={{ display:'grid', gridTemplateColumns:'1fr 1fr', gap:11 }}>
        <Field label="First Name"><input value={f.first_name} onChange={set('first_name')} /></Field>
        <Field label="Last Name"><input value={f.last_name} onChange={set('last_name')} /></Field>
      </div>
      <Field label="Phone (numbers only)"><input value={f.phone_number} onChange={set('phone_number')} placeholder="44999999999" /></Field>
      <Field label="Street"><input value={f.street} onChange={set('street')} placeholder="Rua das Flores, 123" /></Field>
      <div style={{ display:'grid', gridTemplateColumns:'1fr 1fr', gap:11 }}>
        <Field label="City"><input value={f.city} onChange={set('city')} /></Field>
        <Field label="State"><input value={f.state} onChange={set('state')} /></Field>
      </div>
      <div style={{ display:'grid', gridTemplateColumns:'1fr 1fr', gap:11 }}>
        <Field label="Country"><input value={f.country} onChange={set('country')} maxLength={2} placeholder="BR" /></Field>
        <Field label="ZIP Code"><input value={f.zip_code} onChange={set('zip_code')} /></Field>
      </div>
      <Btn variant="primary" style={{ width:'100%', justifyContent:'center' }} onClick={save} loading={busy}>Save Changes</Btn>
    </Modal>
  )
}

function ChangePasswordModal({ open, onClose, email }) {
  const toast = useToast()
  const [step, setStep] = useState(1)
  const [code, setCode] = useState('')
  const [np,   setNp]   = useState('')
  const [err,  setErr]  = useState('')
  const [ok_,  setOk]   = useState('')
  const [busy, setBusy] = useState(false)

  async function sendCode() {
    setBusy(true); await api.forgotPassword({ email }); setBusy(false)
    setOk('Code sent! Check API logs in dev mode.'); setStep(2)
  }

  async function reset() {
    setErr(''); setBusy(true)
    const { ok, data } = await api.resetPassword({ email, code, new_password: np })
    setBusy(false)
    if (ok) { toast('Password changed! Please sign in again.','success'); onClose() }
    else setErr(data.error||'Failed')
  }

  return (
    <Modal open={open} onClose={onClose} title="CHANGE PASSWORD">
      <p style={{ fontSize:13, color:'var(--muted)', marginBottom:14 }}>A reset code will be sent to <strong>{email}</strong></p>
      {err && <Alert type="error"   style={{ marginBottom:11 }}>{err}</Alert>}
      {ok_ && <Alert type="success" style={{ marginBottom:11 }}>{ok_}</Alert>}
      {step === 1 && <Btn variant="primary" style={{ width:'100%', justifyContent:'center' }} onClick={sendCode} loading={busy}>Send Reset Code</Btn>}
      {step === 2 && <>
        <Field label="Code"><input value={code} onChange={e=>setCode(e.target.value)} placeholder="6-digit code" maxLength={6} /></Field>
        <Field label="New Password"><input type="password" value={np} onChange={e=>setNp(e.target.value)} /></Field>
        <Btn variant="primary" style={{ width:'100%', justifyContent:'center' }} onClick={reset} loading={busy}>Reset Password</Btn>
      </>}
    </Modal>
  )
}
