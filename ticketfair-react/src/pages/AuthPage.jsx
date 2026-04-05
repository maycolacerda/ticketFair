// src/pages/AuthPage.jsx
import { useState } from 'react'
import { useAuth } from '../context/AuthContext'
import { useToast } from '../context/ToastContext'
import { api } from '../api/client'
import { Btn, Field, Alert, Spinner, Tabs } from '../components/ui'

const fd = "'Bebas Neue', Impact, sans-serif"
const fb = "'Outfit', system-ui, sans-serif"
const fm = "'JetBrains Mono', monospace"

export default function AuthPage() {
  const { login, register } = useAuth()
  const toast = useToast()
  const [tab, setTab] = useState('login')
  const [busy, setBusy] = useState(false)
  const [err,  setErr]  = useState('')
  const [ok,   setOk]   = useState('')

  // login fields
  const [lEmail, setLE] = useState('')
  const [lPass,  setLP] = useState('')
  // register fields
  const [rEmail, setRE] = useState('')
  const [rUser,  setRU] = useState('')
  const [rPass,  setRP] = useState('')
  // forgot fields
  const [fEmail, setFE] = useState('')
  const [fCode,  setFC] = useState('')
  const [fNew,   setFN] = useState('')
  const [fStep,  setFS] = useState(1)

  const clear = () => { setErr(''); setOk('') }

  async function doLogin(e) {
    e.preventDefault(); clear(); setBusy(true)
    try { await login(lEmail, lPass) }
    catch(e) { setErr(e.message) }
    finally { setBusy(false) }
  }

  async function doRegister(e) {
    e.preventDefault(); clear(); setBusy(true)
    try { await register(rEmail, rUser, rPass) }
    catch(e) { setErr(e.message) }
    finally { setBusy(false) }
  }

  async function doForgot(e) {
    e.preventDefault(); clear(); setBusy(true)
    await api.forgotPassword({ email: fEmail })
    setBusy(false)
    setOk('Code sent! In dev mode, check the API container logs.')
    setFS(2)
  }

  async function doReset(e) {
    e.preventDefault(); clear(); setBusy(true)
    const { ok: rOk, data } = await api.resetPassword({ email: fEmail, code: fCode, new_password: fNew })
    setBusy(false)
    if (rOk) { setOk('Password reset! Sign in with your new password.'); setTab('login'); setFS(1) }
    else setErr(data.error || 'Reset failed')
  }

  return (
    <div style={{ minHeight:'100vh', display:'flex', alignItems:'center', justifyContent:'center', background:'var(--black)', backgroundImage:'radial-gradient(ellipse 60% 50% at 50% 100%, rgba(212,255,0,.04) 0%, transparent 70%)' }}>
      <div style={{ width:440, maxWidth:'95vw', background:'var(--dark)', border:'1px solid var(--border)', borderRadius:14, padding:40 }}>

        <div style={{ fontFamily:fd, fontSize:36, letterSpacing:4, color:'var(--lime)', marginBottom:2 }}>TICKETFAIR</div>
        <div style={{ fontSize:11, fontFamily:fm, color:'var(--muted)', letterSpacing:1.5, marginBottom:32 }}>FAIR TICKETS · NO SCALPERS · REAL CONNECTIONS</div>

        <Tabs
          active={tab}
          onChange={t => { setTab(t); clear(); setFS(1) }}
          tabs={[{ value:'login', label:'Sign In' }, { value:'register', label:'Register' }, { value:'forgot', label:'Forgot' }]}
          style={{ marginBottom:20 }}
        />

        {err && <Alert type="error"   style={{ marginBottom:14 }}>{err}</Alert>}
        {ok  && <Alert type="success" style={{ marginBottom:14 }}>{ok}</Alert>}

        {/* ── Login ── */}
        {tab === 'login' && (
          <form onSubmit={doLogin}>
            <Field label="Email"><input type="email" value={lEmail} onChange={e=>setLE(e.target.value)} placeholder="you@example.com" required /></Field>
            <Field label="Password"><input type="password" value={lPass} onChange={e=>setLP(e.target.value)} placeholder="••••••••" required /></Field>
            <Submit busy={busy}>Sign In →</Submit>
          </form>
        )}

        {/* ── Register ── */}
        {tab === 'register' && (
          <form onSubmit={doRegister}>
            <Field label="Email"><input type="email" value={rEmail} onChange={e=>setRE(e.target.value)} placeholder="you@example.com" required /></Field>
            <Field label="Username"><input value={rUser} onChange={e=>setRU(e.target.value)} placeholder="yourhandle" required /></Field>
            <Field label="Password"><input type="password" value={rPass} onChange={e=>setRP(e.target.value)} placeholder="Min 8 chars, uppercase, number, symbol" required /></Field>
            <Submit busy={busy}>Create Account →</Submit>
          </form>
        )}

        {/* ── Forgot ── */}
        {tab === 'forgot' && (
          <>
            {fStep === 1 && (
              <form onSubmit={doForgot}>
                <Field label="Email"><input type="email" value={fEmail} onChange={e=>setFE(e.target.value)} placeholder="you@example.com" required /></Field>
                <Submit busy={busy}>Send Reset Code →</Submit>
              </form>
            )}
            {fStep === 2 && (
              <form onSubmit={doReset}>
                <Field label="6-Digit Code"><input value={fCode} onChange={e=>setFC(e.target.value)} placeholder="123456" maxLength={6} required /></Field>
                <Field label="New Password"><input type="password" value={fNew} onChange={e=>setFN(e.target.value)} placeholder="New strong password" required /></Field>
                <Submit busy={busy}>Reset Password →</Submit>
              </form>
            )}
          </>
        )}
      </div>
    </div>
  )
}

function Submit({ busy, children }) {
  return (
    <button type="submit" disabled={busy} style={{ width:'100%', padding:13, background:'var(--lime)', color:'var(--black)', border:'none', borderRadius:6, fontFamily:fb, fontSize:14, fontWeight:700, cursor:busy?'not-allowed':'pointer', opacity:busy?.8:1, marginTop:8, display:'flex', alignItems:'center', justifyContent:'center', gap:8, transition:'all .15s' }}>
      {busy ? <Spinner size={14} color="var(--black)" /> : children}
    </button>
  )
}
