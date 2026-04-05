// src/context/AuthContext.jsx
import { createContext, useContext, useState } from 'react'
import { api } from '../api/client'

const Ctx = createContext(null)

export function AuthProvider({ children }) {
  const [user,     setUser]     = useState(() => { try { return JSON.parse(localStorage.getItem('tf_user'))     } catch { return null } })
  const [merchant, setMerchant] = useState(() => { try { return JSON.parse(localStorage.getItem('tf_merchant')) } catch { return null } })
  const [admin,    setAdmin]    = useState(() => { try { return JSON.parse(localStorage.getItem('tf_admin'))    } catch { return null } })

  async function login(email, password) {
    const { ok, data } = await api.login({ email, password })
    if (!ok) throw new Error(data.error || 'Invalid credentials')
    localStorage.setItem('tf_token', data.data.token)
    localStorage.setItem('tf_user',  JSON.stringify(data.data.user))
    setUser(data.data.user)

    // also try merchant login (same creds may work)
    const { ok: mOk, data: mData } = await api.loginMerchant({ email, password })
    if (mOk) {
      localStorage.setItem('tf_merchant_token', mData.data.token)
      localStorage.setItem('tf_merchant', JSON.stringify(mData.data.merchant))
      setMerchant(mData.data.merchant)
    }
    return data.data.user
  }

  async function register(email, username, password) {
    const { ok, data } = await api.register({ email, username, password })
    if (!ok) throw new Error(data.error || (data.errors ? Object.values(data.errors).join(' · ') : 'Registration failed'))
    return login(email, password)
  }

  async function loginAdmin(email, password) {
    const { ok, data } = await api.loginAdmin({ email, password })
    if (!ok) throw new Error(data.error || 'Invalid credentials')
    localStorage.setItem('tf_admin_token', data.data.token)
    localStorage.setItem('tf_admin', JSON.stringify(data.data.admin))
    setAdmin(data.data.admin)
    return data.data.admin
  }

  function logout() {
    ['tf_token','tf_user','tf_merchant_token','tf_merchant','tf_admin_token','tf_admin'].forEach(k => localStorage.removeItem(k))
    setUser(null); setMerchant(null); setAdmin(null)
    api.logout().catch(() => {})
  }

  return (
    <Ctx.Provider value={{
      user, merchant, admin,
      isAuthenticated: !!user,
      isMerchant: !!merchant,
      isAdmin: !!admin,
      login, register, loginAdmin, logout,
    }}>
      {children}
    </Ctx.Provider>
  )
}

export const useAuth = () => useContext(Ctx)
