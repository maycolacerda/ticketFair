// src/api/client.js
const BASE = import.meta.env.VITE_API_URL || 'http://localhost:8000/api/v1'

function tok()  { return localStorage.getItem('tf_token')          || '' }
function mTok() { return localStorage.getItem('tf_merchant_token') || '' }
function aTok() { return localStorage.getItem('tf_admin_token')    || '' }

async function req(path, method = 'GET', body = null, auth = 'client') {
  const token = auth === 'merchant' ? mTok() : auth === 'admin' ? aTok() : tok()
  const headers = { 'Content-Type': 'application/json' }
  if (token) headers['Authorization'] = `Bearer ${token}`
  const opts = { method, headers }
  if (body) opts.body = JSON.stringify(body)
  try {
    const res = await fetch(BASE + path, opts)
    const data = await res.json().catch(() => ({}))
    return { ok: res.ok, status: res.status, data }
  } catch {
    return { ok: false, status: 0, data: { error: 'Connection failed — is the API running?' } }
  }
}

// helpers
const G  = (p)    => req(p)
const P  = (p, b) => req(p, 'POST',   b)
const Pu = (p, b) => req(p, 'PUT',    b)
const D  = (p)    => req(p, 'DELETE')
const MG  = (p)    => req(p, 'GET',    null, 'merchant')
const MP  = (p, b) => req(p, 'POST',   b,    'merchant')
const MPu = (p, b) => req(p, 'PUT',    b,    'merchant')
const MD  = (p)    => req(p, 'DELETE', null, 'merchant')
const AG  = (p)    => req(p, 'GET',    null, 'admin')
const AP  = (p, b) => req(p, 'POST',   b,    'admin')

export const api = {
  // ── Auth ──────────────────────────────────────────
  login:          (b) => P('/public/auth/client/login',    b),
  loginMerchant:  (b) => P('/public/auth/merchant/login',  b),
  loginAdmin:     (b) => P('/admin/auth/login',            b),
  register:       (b) => P('/public/auth/register',        b),
  logout:         ()  => P('/public/auth/logout'),
  forgotPassword: (b) => P('/public/auth/password/forgot', b),
  resetPassword:  (b) => P('/public/auth/password/reset',  b),

  // ── Public events ─────────────────────────────────
  getEvents:      (qs = '') => G(`/public/events?${qs}`),
  getEvent:       (id)      => G(`/public/events/${id}`),
  getTicketTypes: (id)      => G(`/public/events/${id}/ticket-types`),

  // ── Private — user ────────────────────────────────
  getProfile:    ()    => G('/private/profile/'),
  createProfile: (b)   => P('/private/profile/', b),
  updateProfile: (b)   => Pu('/private/profile/', b),

  // ── Verification ──────────────────────────────────
  sendEmailVerify: () => P('/private/verify/email/send'),
  verifyEmail:  (b)   => P('/private/verify/email', b),
  sendPhoneVerify: () => P('/private/verify/phone/send'),
  verifyPhone:  (b)   => P('/private/verify/phone', b),

  // ── Tickets ───────────────────────────────────────
  getTickets:     (qs = '') => G(`/private/tickets?${qs}`),
  getTicket:      (id)      => G(`/private/tickets/${id}`),
  purchaseTicket: (b)       => P('/private/tickets/purchase', b),
  refundTicket:   (b)       => P('/private/tickets/refund', b),
  giftTicket:     (id, b)   => P(`/private/tickets/${id}/gift`, b),

  // ── Transactions ──────────────────────────────────
  getTransactions: (qs = '') => G(`/private/transactions?${qs}`),

  // ── Payments (Stripe) ─────────────────────────────
  getPayments:   (qs = '') => G(`/private/payments?${qs}`),
  createIntent:  (b)       => P('/private/payments/intent', b),
  refundPayment: (id)      => P(`/private/payments/${id}/refund`),

  // ── Connections ───────────────────────────────────
  getConnections:      (qs = '') => G(`/private/connections?${qs}`),
  getPendingRequests:  ()        => G('/private/connections/requests'),
  getConnectionEvents: (qs = '') => G(`/private/connections/events?${qs}`),
  sendConnectionReq:   (b)       => P('/private/connections', b),
  respondConnection:   (id, b)   => P(`/private/connections/${id}/respond`, b),
  removeConnection:    (id)      => D(`/private/connections/${id}`),

  // ── Merchant ──────────────────────────────────────
  createEvent:      (b)          => MP('/merchant/events/new', b),
  updateEvent:      (id, b)      => MPu(`/merchant/events/${id}`, b),
  getMerchantTT:    (eid)        => MG(`/merchant/events/${eid}/ticket-types`),
  createTT:         (eid, b)     => MP(`/merchant/events/${eid}/ticket-types`, b),
  updateTT:         (eid, id, b) => MPu(`/merchant/events/${eid}/ticket-types/${id}`, b),
  deleteTT:         (eid, id)    => MD(`/merchant/events/${eid}/ticket-types/${id}`),
  validateTicket:   (id)         => MP(`/merchant/tickets/${id}/validate`),
  uploadEventImage: (id, form)   => {
    return fetch(`${BASE}/merchant/events/${id}/image`, {
      method: 'POST',
      headers: mTok() ? { Authorization: `Bearer ${mTok()}` } : {},
      body: form,
    }).then(r => r.json().then(d => ({ ok: r.ok, data: d })))
  },

  // ── Admin ─────────────────────────────────────────
  adminGetUsers:           (qs = '') => AG(`/admin/users?${qs}`),
  adminCreateUser:         (b)       => AP('/admin/users', b),
  adminDeactivateUser:     (id)      => AP(`/admin/users/${id}/deactivate`),
  adminActivateUser:       (id)      => AP(`/admin/users/${id}/activate`),
  adminGetMerchants:       (qs = '') => AG(`/admin/merchants?${qs}`),
  adminCreateMerchant:     (b)       => AP('/admin/merchants', b),
  adminDeactivateMerchant: (id)      => AP(`/admin/merchants/${id}/deactivate`),
  adminActivateMerchant:   (id)      => AP(`/admin/merchants/${id}/activate`),
}
