import React, { useEffect, useState } from 'react'
import { 
  Printer, 
  Wifi, 
  WifiOff, 
  Play, 
  Pause, 
  AlertCircle, 
  CheckCircle,
  Clock,
  Thermometer,
  Zap
} from 'lucide-react'
import { useLocalization } from '../../../hooks/useLocalization.jsx'

function PrintersTable({ SetNowPrinter, SetDetailsIsOpen }) {
  const [printers, setPrinters] = useState([])
  const [error, setError] = useState(null)
  const [loading, setLoading] = useState(true)
  const { t } = useLocalization()

  // Mock data for demonstration
  const mockPrinters = [
    {
      uniqueKey: 'printer-001',
      printerName: 'Ender 3 Pro',
      printerType: 'FDM',
      version: 'v2.0.1',
      isWorking: true,
      typeOfConnection: 'USB',
      isPrinting: false,
      nozzleTemp: 200,
      bedTemp: 60,
      progress: 0,
      timeRemaining: 0
    },
    {
      uniqueKey: 'printer-002', 
      printerName: 'Prusa i3 MK3S+',
      printerType: 'FDM',
      version: 'v3.0.0',
      isWorking: true,
      typeOfConnection: 'WiFi',
      isPrinting: true,
      nozzleTemp: 215,
      bedTemp: 65,
      progress: 45,
      timeRemaining: 3600
    },
    {
      uniqueKey: 'printer-003',
      printerName: 'Anycubic Photon',
      printerType: 'SLA',
      version: 'v1.5.2',
      isWorking: false,
      typeOfConnection: 'COM',
      isPrinting: false,
      nozzleTemp: 0,
      bedTemp: 0,
      progress: 0,
      timeRemaining: 0
    }
  ]

  function updateTableData(newPrinters) {
    let changed = false
    const updated = [...printers]

    for (const newPrinter of newPrinters) {
      const index = updated.findIndex(p => p.uniqueKey === newPrinter.uniqueKey)

      if (index !== -1) {
        const old = updated[index]
        if (JSON.stringify(old) !== JSON.stringify(newPrinter)) {
          updated[index] = newPrinter
          changed = true
        }
      } else {
        updated.push(newPrinter)
        changed = true
      }
    }

    if (changed) {
      setPrinters(updated)
    }
  }

  useEffect(() => {
    const getDataFromServer = async () => {
      try {
        setLoading(true)
        const response = await fetch('/api/Printers', {
          method: 'GET',
          headers: { 'Content-Type': 'application/json' }
        })
        
        if (!response.ok) throw new Error('Failed to fetch printer list')
        
        const json = await response.json()
        updateTableData(json)
        setError(null)
      } catch (err) {
        console.log("Error loading printers: " + err.message)
        // Убираем показ ошибки - просто оставляем пустой список
        setError(null)
        
        // Не загружаем mock данные - оставляем пустой список
        setTimeout(() => {
          updateTableData([])
          setError(null)
          console.log('Нет подключенных принтеров')
        }, 500)
      } finally {
        setLoading(false)
      }
    }

    getDataFromServer()
    const interval = setInterval(getDataFromServer, 10000) // Увеличили интервал до 10 секунд
    return () => clearInterval(interval)
  }, [])

  const showPrinterDetails = (printer) => {
    SetNowPrinter(printer)
    SetDetailsIsOpen(printer)
  }

  const getStatusIcon = (printer) => {
    if (!printer.isWorking) {
      return <WifiOff className="h-4 w-4 text-danger-500" />
    }
    if (printer.isPrinting) {
      return <Play className="h-4 w-4 text-success-500 animate-pulse" />
    }
    return <Wifi className="h-4 w-4 text-success-500" />
  }

  const getStatusText = (printer) => {
    if (!printer.isWorking) return t('printers.status.offline')
    if (printer.isPrinting) return t('printers.status.printing')
    return t('printers.status.ready')
  }

  const getStatusClass = (printer) => {
    if (!printer.isWorking) return 'status-offline'
    if (printer.isPrinting) return 'status-printing'
    return 'status-online'
  }

  const formatTime = (seconds) => {
    if (seconds === 0) return '--'
    const hours = Math.floor(seconds / 3600)
    const minutes = Math.floor((seconds % 3600) / 60)
    return `${hours}${t('time.hours')} ${minutes}${t('time.minutes')}`
  }

  if (loading) {
    return (
      <div className="card">
        <div className="flex items-center justify-center py-12">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary-600"></div>
          <span className="ml-3 text-gray-600 dark:text-gray-300">{t('status.loading')} {t('navigation.printers').toLowerCase()}...</span>
        </div>
      </div>
    )
  }

  return (
    <div className="card">
      {/* Header */}
      <div className="flex items-center justify-between mb-6">
        <div className="flex items-center space-x-3">
          <Printer className="h-6 w-6 text-primary-600 dark:text-primary-400" />
          <h2 className="text-xl font-semibold text-gray-900 dark:text-gray-100">{t('navigation.printers')}</h2>
          <span className="bg-primary-100 text-primary-800 text-xs font-medium px-2.5 py-0.5 rounded-full dark:bg-primary-900/40 dark:text-primary-300">
            {printers.length}
          </span>
        </div>
        <div className="flex items-center space-x-2 text-sm text-gray-500 dark:text-gray-400">
          <div className="flex items-center space-x-1">
            <div className="w-2 h-2 bg-success-500 rounded-full"></div>
            <span>{t('header.online')}: {printers.filter(p => p.isWorking).length}</span>
          </div>
          <div className="flex items-center space-x-1">
            <div className="w-2 h-2 bg-warning-500 rounded-full animate-pulse"></div>
            <span>{t('header.printing')}: {printers.filter(p => p.isPrinting).length}</span>
          </div>
        </div>
      </div>

      {/* Убрали показ ошибки */}

      {/* Table */}
      <div className="overflow-x-auto">
        <table className="w-full">
          <thead>
            <tr className="table-header">
              <th className="px-4 py-3 text-left rounded-l-lg">{t('navigation.printers')}</th>
              <th className="px-4 py-3 text-left">{t('printers.details.type')}</th>
              <th className="px-4 py-3 text-left">{t('printers.details.status')}</th>
              <th className="px-4 py-3 text-left">{t('printers.details.nozzle')} / {t('printers.details.bed')}</th>
              <th className="px-4 py-3 text-left">{t('printers.details.progress')}</th>
              <th className="px-4 py-3 text-left rounded-r-lg">{t('printers.details.connection')}</th>
            </tr>
          </thead>
          <tbody>
            {printers.map((printer, index) => (
              <tr
                key={printer.uniqueKey}
                className="table-row border-b border-gray-100 hover:bg-primary-50 transition-colors duration-200 cursor-pointer dark:border-gray-700 dark:hover:bg-gray-800"
                onClick={() => showPrinterDetails(printer)}
              >
                <td className="px-4 py-4">
                  <div className="flex items-center space-x-3">
                    <div className="p-2 bg-primary-100 rounded-lg dark:bg-primary-900/40">
                      <Printer className="h-5 w-5 text-primary-600 dark:text-primary-400" />
                    </div>
                    <div>
                      <div className="font-medium text-gray-900 dark:text-gray-100">{printer.printerName || 'Unnamed'}</div>
                      <div className="text-sm text-gray-500 dark:text-gray-400">v{printer.version || 'Unknown'}</div>
                    </div>
                  </div>
                </td>
                <td className="px-4 py-4">
                  <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-gray-100 text-gray-800 dark:bg-gray-800 dark:text-gray-200">
                    {printer.printerType || 'Unknown'}
                  </span>
                </td>
                <td className="px-4 py-4">
                  <div className="flex items-center space-x-2">
                    {getStatusIcon(printer)}
                    <span className={getStatusClass(printer)}>
                      {getStatusText(printer)}
                    </span>
                  </div>
                </td>
                <td className="px-4 py-4">
                  <div className="flex items-center space-x-3 text-sm">
                    <div className="flex items-center space-x-1">
                      <Thermometer className="h-4 w-4 text-orange-500" />
                      <span>{printer.nozzleTemp}°C</span>
                    </div>
                    <div className="flex items-center space-x-1">
                      <Zap className="h-4 w-4 text-blue-500" />
                      <span>{printer.bedTemp}°C</span>
                    </div>
                  </div>
                </td>
                <td className="px-4 py-4">
                  {printer.isPrinting ? (
                    <div className="space-y-1">
                      <div className="flex items-center justify-between text-sm">
                        <span>{printer.progress}%</span>
                        <Clock className="h-4 w-4 text-gray-400" />
                      </div>
                      <div className="w-full bg-gray-200 rounded-full h-2 dark:bg-gray-700">
                        <div
                          className="bg-success-500 h-2 rounded-full transition-all duration-500"
                          style={{ width: `${printer.progress}%` }}
                        />
                      </div>
                      <div className="text-xs text-gray-500 dark:text-gray-400">
                        {formatTime(printer.timeRemaining)}
                      </div>
                    </div>
                  ) : (
                    <span className="text-gray-400">--</span>
                  )}
                </td>
                <td className="px-4 py-4">
                  <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-blue-100 text-blue-800 dark:bg-blue-900/40 dark:text-blue-300">
                    {printer.typeOfConnection || 'Unknown'}
                  </span>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      {/* Empty State */}
      {printers.length === 0 && !loading && (
        <div className="text-center py-12">
          <Printer className="h-12 w-12 text-gray-400 mx-auto mb-4" />
          <h3 className="text-lg font-medium text-gray-900 mb-2 dark:text-gray-100">
            {t('printers.empty.title')}
          </h3>
          <p className="text-gray-500 dark:text-gray-400">
            {t('printers.empty.subtitle')}
          </p>
        </div>
      )}
    </div>
  )
}
export default PrintersTable

