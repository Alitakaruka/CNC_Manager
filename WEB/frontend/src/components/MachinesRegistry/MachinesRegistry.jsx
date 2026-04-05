import { useEffect, useState, useRef, useCallback } from 'react'
import { Server, RefreshCw, Plug } from 'lucide-react'
import toast from 'react-hot-toast'
import { useLocalization } from '../../hooks/useLocalization.jsx'
import wsClient from '../../hooks/WebSocketClient'
import { ReconnectCNC } from '../../hooks/ConnectionHook'
import './MachinesRegistry.css'

const WS_REGISTRY_REQUEST = 'GetRegistry'
/** Событие, если сервер шлёт обёртку { type: "registry", data: [...] } */
const WS_REGISTRY_EVENT = 'registry'

function getRegistryFlat(row) {
  if (!row || typeof row !== 'object') return {}
  const dto = row.DTO
  if (dto && typeof dto === 'object') {
    return { ...dto, ...row }
  }
  return row
}

function getRegistryRowKey(row) {
  const r = getRegistryFlat(row)
  return r.uniqueKey || r.UniqueKey || ''
}

function getMachineName(m) {
  const r = getRegistryFlat(m)
  return r.CNCName || r.TARGET_MACHINE_NAME || r.printerName || r.customName || ''
}

function getConnectionMethod(m) {
  const r = getRegistryFlat(m)
  const raw =
    r.connectionType ??
    r.ConnectionType ??
    r.typeOfConnection ??
    r.TypeOfConnection ??
    ''
  return String(raw).trim()
}

function getConnectionDetails(m) {
  const r = getRegistryFlat(m)
  const details = r.connectionData ?? r.ConnectionData ?? r.connectionString ?? r.ConnectionString ?? ''
  return String(details).trim()
}

/** Ответ GetRegistry: массив, или { data: [] }, или событие registry уже с массивом */
function normalizeRegistryPayload(raw) {
  if (raw == null) return []
  if (Array.isArray(raw)) return raw
  if (typeof raw === 'object') {
    if (Array.isArray(raw.data)) return raw.data
    if (Array.isArray(raw.machines)) return raw.machines
  }
  return []
}

function looksLikeRegistryMessage(msg) {
  if (!Array.isArray(msg)) return false
  if (msg.length === 0) return true
  return msg.every((el) => el && typeof el === 'object' && !Array.isArray(el))
}

export default function MachinesRegistry() {
  const [rows, setRows] = useState([])
  const [loading, setLoading] = useState(true)
  const [connectingKey, setConnectingKey] = useState('')
  const { t } = useLocalization()
  const startedRef = useRef(false)

  const requestRegistryViaWs = useCallback(async () => {
    try {
      if (wsClient.isConnected) {
        await wsClient.request(WS_REGISTRY_REQUEST, {})
        return true
      }
    } catch (err) {
      console.log('MachinesRegistry', WS_REGISTRY_REQUEST, err.message)
    }
    return false
  }, [])

  const applyRegistryPayload = useCallback((raw) => {
    setRows(normalizeRegistryPayload(raw))
    setLoading(false)
  }, [])

  const refresh = useCallback(async () => {
    setLoading(true)
    await requestRegistryViaWs()
    setLoading(false)
  }, [requestRegistryViaWs])

  useEffect(() => {
    if (!startedRef.current) {
      startedRef.current = true
      wsClient.connect()
    }

    const offRegistry = wsClient.on(WS_REGISTRY_EVENT, (data) => {
      applyRegistryPayload(data)
    })

    const offMessage = wsClient.on('message', (msg) => {
      if (looksLikeRegistryMessage(msg)) {
        applyRegistryPayload(msg)
      }
    })

    const load = async () => {
      setLoading(true)
      if (wsClient.isConnected) {
        await requestRegistryViaWs()
        setLoading(false)
      } else {
        setLoading(false)
      }
    }
    load()

    const offOpen = wsClient.on('open', async () => {
      setLoading(true)
      await requestRegistryViaWs()
      setLoading(false)
    })

    return () => {
      offRegistry()
      offMessage()
      offOpen()
    }
  }, [applyRegistryPayload, requestRegistryViaWs])

  const handleConnect = async (row) => {
    const key = getRegistryRowKey(row)
    if (!key) {
      toast.error(t('machinesPage.missingKey'))
      return
    }
    setConnectingKey(key)
    try {
      await ReconnectCNC(key)
      toast.success(t('machinesPage.connectSuccess'))
      await refresh()
    } catch (e) {
      toast.error(e?.message || t('machinesPage.connectError'))
    } finally {
      setConnectingKey('')
    }
  }

  const methodLabel = (raw) => {
    const u = String(raw).toUpperCase()
    if (u === 'COM') return t('machinesPage.methodCOM')
    if (u === 'IP') return t('machinesPage.methodIP')
    if (u === 'WIFI' || u === 'WI-FI') return t('machinesPage.methodWifi')
    if (u === 'USB') return t('machinesPage.methodUSB')
    return raw || t('machinesPage.methodUnknown')
  }

  return (
    <div className="machines-registry space-y-4">
      <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-3">
        <div className="flex items-center gap-3">
          <Server className="h-7 w-7 text-primary-600 dark:text-primary-400" />
          <div>
            <h2 className="text-xl font-semibold text-gray-900 dark:text-gray-100">
              {t('machinesPage.title')}
            </h2>
            <p className="text-sm text-gray-500 dark:text-gray-400">{t('machinesPage.subtitle')}</p>
          </div>
        </div>
        <button
          type="button"
          onClick={() => refresh()}
          disabled={loading}
          className="inline-flex items-center gap-2 px-4 py-2 rounded-lg bg-primary-100 text-primary-800 hover:bg-primary-200 dark:bg-gray-800 dark:text-primary-300 dark:hover:bg-gray-700 transition-colors disabled:opacity-50"
        >
          <RefreshCw className={`h-4 w-4 ${loading ? 'animate-spin' : ''}`} />
          {t('machinesPage.refresh')}
        </button>
      </div>

      <div className="machines-registry__table-wrap glass-effect rounded-xl border border-white/20 dark:border-gray-700 overflow-x-auto">
        <table className="w-full">
          <thead>
            <tr className="table-header">
              <th className="px-4 py-3 text-left rounded-tl-lg">{t('machinesPage.columns.name')}</th>
              <th className="px-4 py-3 text-left">{t('machinesPage.columns.uniqueKey')}</th>
              <th className="px-4 py-3 text-left">{t('machinesPage.columns.connectionMethod')}</th>
              <th className="px-4 py-3 text-left">{t('machinesPage.columns.connectionData')}</th>
              <th className="px-4 py-3 text-left rounded-tr-lg">{t('machinesPage.columns.actions')}</th>
            </tr>
          </thead>
          <tbody>
            {loading && rows.length === 0 ? (
              <tr>
                <td colSpan={5} className="px-4 py-8 text-center text-gray-500 dark:text-gray-400">
                  {t('status.loading')}
                </td>
              </tr>
            ) : rows.length === 0 ? (
              <tr>
                <td colSpan={5} className="px-4 py-8 text-center text-gray-500 dark:text-gray-400">
                  {t('machinesPage.empty')}
                </td>
              </tr>
            ) : (
              rows.map((row) => {
                const key = getRegistryRowKey(row)
                const name = getMachineName(row) || t('machinesPage.unnamed')
                const methodRaw = getConnectionMethod(row)
                const details = getConnectionDetails(row)
                const busy = connectingKey === key
                return (
                  <tr
                    key={key || `${name}-${methodRaw}-${details}`}
                    className="table-row border-b border-gray-100 dark:border-gray-700 hover:bg-primary-50/80 dark:hover:bg-gray-800/80 transition-colors"
                  >
                    <td className="px-4 py-3 font-medium text-gray-900 dark:text-gray-100">{name}</td>
                    <td className="px-4 py-3 text-sm font-mono text-gray-700 dark:text-gray-300">{key || '—'}</td>
                    <td className="px-4 py-3 text-sm text-gray-700 dark:text-gray-300">
                      {methodRaw ? methodLabel(methodRaw) : t('machinesPage.connectionDataUnavailable')}
                    </td>
                    <td className="px-4 py-3 text-sm text-gray-600 dark:text-gray-400 break-all max-w-[220px]">
                      {details || t('machinesPage.connectionDataUnavailable')}
                    </td>
                    <td className="px-4 py-3">
                      <button
                        type="button"
                        disabled={!key || busy}
                        onClick={() => handleConnect(row)}
                        className="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-lg text-sm font-medium bg-primary-600 text-white hover:bg-primary-700 disabled:opacity-40 disabled:pointer-events-none dark:bg-primary-500 dark:hover:bg-primary-600"
                      >
                        <Plug className="h-3.5 w-3.5" />
                        {busy ? t('machinesPage.connecting') : t('machinesPage.connect')}
                      </button>
                    </td>
                  </tr>
                )
              })
            )}
          </tbody>
        </table>
      </div>
    </div>
  )
}
