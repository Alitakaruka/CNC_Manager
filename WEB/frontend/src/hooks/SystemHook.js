import { useState, useEffect, useRef } from 'react'
import wsClient from './WebSocketClient'

export const useSystemData = () => {
  const [systemData, setSystemData] = useState({
    connectionStatus: { online: 0, printing: 0, offline: 0, total: 0 },
    systemInfo: { uptime: '0d 0h 0m', activeConnections: 0, lastUpdate: new Date() },
    logs: []
  })
  const [isLoading, setIsLoading] = useState(false)
  const startedRef = useRef(false)

  // WS handlers
  useEffect(() => {
    if (!startedRef.current) {
      startedRef.current = true
      wsClient.connect()
    }
    const offStatus = wsClient.on('status', (data) => {
      setSystemData(prev => ({ ...prev, connectionStatus: data }))
    })
    const offInfo = wsClient.on('systemInfo', (data) => {
      setSystemData(prev => ({ 
  ...prev, 
  systemInfo: { ...data, lastUpdate: new Date() }
}))
    })
    const offLog = wsClient.on('log', (entry) => {
      setSystemData(prev => ({ ...prev, logs: [...prev.logs.slice(-9), entry] }))
    })
    return () => { offStatus(); offInfo(); offLog() }
  }, [])

  // HTTP fallback and manual updates
  const fetchSnapshot = async () => {
    try {
      const res = await fetch('/api/system/snapshot')
      if (res.ok) {
        const data = await res.json()
        setSystemData(prev => ({
          ...prev,
          connectionStatus: data.connectionStatus || prev.connectionStatus,
          systemInfo: data.systemInfo ? { ...data.systemInfo, lastUpdate: new Date() } : prev.systemInfo
        }))
      }
    } catch {}
  }

  const updateConnectionStatus = async () => {
    try {
      const res = await fetch('/api/system/status')
      if (res.ok) {
        const data = await res.json()
        setSystemData(prev => ({ ...prev, connectionStatus: data }))
        return data
      }
    } catch {}
    return systemData.connectionStatus
  }

  const updateSystemInfo = async () => {
    try {
      const res = await fetch('/api/system/info')
      if (res.ok) {
        const data = await res.json()
        setSystemData(prev => ({ ...prev, systemInfo: data }))
        return data
      }
    } catch {}
    return null
  }

  const updateLogs = async () => {
    try {
      const res = await fetch('/api/system/logs')
      if (res.ok) {
        const data = await res.json()
        setSystemData(prev => ({ ...prev, logs: data.slice(-10) }))
        return data
      }
    } catch {}
    return []
  }

  const updateAllData = async () => {
    setIsLoading(true)
    try {
      if (!wsClient.connected) await fetchSnapshot()
    } finally {
      setIsLoading(false)
    }
  }

  useEffect(() => {
    updateAllData()
    let interval
    const visibilityHandler = () => {
      clearInterval(interval)
      const fast = document.visibilityState === 'visible' && (systemData.connectionStatus.printing > 0)
      const ms = fast ? 1500 : 10000
      interval = setInterval(() => { if (!wsClient.connected) fetchSnapshot() }, ms)
    }
    document.addEventListener('visibilitychange', visibilityHandler)
    visibilityHandler()
    return () => { document.removeEventListener('visibilitychange', visibilityHandler); clearInterval(interval) }
  }, [systemData.connectionStatus.printing])

  return { systemData, isLoading, updateConnectionStatus, updateSystemInfo, updateLogs, updateAllData }
}

export default useSystemData
