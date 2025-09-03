import { useState, useRef } from 'react'
import { motion } from 'framer-motion'
import { 
  X, 
  Printer, 
  Thermometer, 
  Zap, 
  Play, 
  Pause, 
  Square, 
  Upload, 
  FileText,
  Settings,
  Gauge,
  Move,
  ArrowUp,
  ArrowDown,
  ArrowLeft,
  ArrowRight,
  RotateCcw,
  Home,
  Target
} from 'lucide-react'
import MainPrinterData from '../MainPrinterData/MainPrinterData'
import PrinterGabarites from '../Gabarites/gabarites'
import PrinterCommands from '../Commands/PrinterCommands'
import toast from 'react-hot-toast'
import { useLocalization } from '../../../hooks/useLocalization.jsx'

export default function Details({ PrinterData, SetDetailsIsOpen }) {
  const fileRef = useRef(null)
  const [isUploading, setIsUploading] = useState(false)
  const { t } = useLocalization()

  if (!PrinterData) {
    return null
  }

  const handleStartPrint = async () => {
    if (!fileRef.current?.files?.[0]) {
      toast.error(t('notifications.fileRequired'))
      return
    }

    const file = fileRef.current.files[0]
    const formData = new FormData()
    formData.append('PrintFile', file)

    if (!PrinterData.uniqueKey) {
      toast.error(t('notifications.printerNotSelected'))
      return
    }

    setIsUploading(true)

    try {
      const query = `uniqueKey=${encodeURIComponent(PrinterData.uniqueKey)}`
      const response = await fetch(`/api/StartPrint?${query}`, {
        method: 'POST',
        body: formData
      })

      if (!response.ok) {
        throw new Error(t('notifications.printError'))
      }

      const result = await response.text()
      toast.success(t('notifications.printStarted'))
      console.log('Ответ сервера:', result)
    } catch (error) {
      toast.error(`${t('common.error')}: ${error.message}`)
    } finally {
      setIsUploading(false)
    }
  }

  const getStatusColor = () => {
    if (!PrinterData.isWorking) return 'text-danger-500'
    if (PrinterData.isPrinting) return 'text-success-500'
    return 'text-primary-500'
  }

  const getStatusText = () => {
    if (!PrinterData.isWorking) return t('printers.status.offline')
    if (PrinterData.isPrinting) return t('printers.status.printing')
    return t('printers.status.ready')
  }

  return (
    <motion.div
      className="card h-fit"
      initial={{ opacity: 0, x: 100, scale: 0.95 }}
      animate={{ opacity: 1, x: 0, scale: 1 }}
      exit={{ opacity: 0, x: 100, scale: 0.95 }}
      transition={{ type: "spring", damping: 25, stiffness: 300 }}
    >
      {/* Header */}
      <div className="flex items-center justify-between mb-6">
        <div className="flex items-center space-x-3">
          <div className="p-2 bg-primary-100 rounded-lg dark:bg-primary-900/40">
            <Printer className="h-6 w-6 text-primary-600 dark:text-primary-400" />
          </div>
          <div>
            <h3 className="text-lg font-semibold text-gray-900 dark:text-gray-100">
              {PrinterData.printerName || 'Unnamed Printer'}
            </h3>
            <div className="flex items-center space-x-2">
              <span className={`text-sm font-medium ${getStatusColor()}`}>
                {getStatusText()}
              </span>
              <span className="text-xs text-gray-500 dark:text-gray-400">
                {PrinterData.uniqueKey}
              </span>
            </div>
          </div>
        </div>
        <motion.button
          onClick={SetDetailsIsOpen}
          className="p-2 text-gray-400 hover:text-gray-600 hover:bg-gray-100 rounded-lg transition-colors duration-200 dark:hover:bg-gray-800"
          whileHover={{ scale: 1.1 }}
          whileTap={{ scale: 0.9 }}
        >
          <X className="h-5 w-5" />
        </motion.button>
      </div>

      {/* Temperature Status */}
      <div className="grid grid-cols-2 gap-4 mb-6">
        <motion.div
          className="p-4 bg-orange-50 border border-orange-200 rounded-lg dark:bg-orange-900/20 dark:border-orange-800"
          whileHover={{ scale: 1.02 }}
        >
          <div className="flex items-center space-x-2 mb-2">
            <Thermometer className="h-5 w-5 text-orange-600" />
            <span className="text-sm font-medium text-orange-800 dark:text-orange-300">{t('printers.details.nozzle')}</span>
          </div>
          <div className="text-2xl font-bold text-orange-900 dark:text-orange-200">
            {PrinterData.nozzleTemp || 0}°C
          </div>
        </motion.div>

        <motion.div
          className="p-4 bg-blue-50 border border-blue-200 rounded-lg dark:bg-blue-900/20 dark:border-blue-800"
          whileHover={{ scale: 1.02 }}
        >
          <div className="flex items-center space-x-2 mb-2">
            <Zap className="h-5 w-5 text-blue-600" />
            <span className="text-sm font-medium text-blue-800 dark:text-blue-300">{t('printers.details.bed')}</span>
          </div>
          <div className="text-2xl font-bold text-blue-900 dark:text-blue-200">
            {PrinterData.bedTemp || 0}°C
          </div>
        </motion.div>
      </div>

      {/* Print Controls */}
      <div className="space-y-4 mb-6">
        <h4 className="text-sm font-semibold text-gray-700 flex items-center space-x-2 dark:text-gray-200">
          <Play className="h-4 w-4" />
          <span>{t('printers.controls.printControl')}</span>
        </h4>

        {/* File Upload */}
        <div className="space-y-3">
          <label className="block text-sm font-medium text-gray-700 dark:text-gray-300">
            {t('printers.controls.gcodeFile')}
          </label>
          <div className="flex items-center space-x-3">
            <input
              type="file"
              ref={fileRef}
              accept=".gcode"
              className="hidden"
              id="fileInput"
            />
            <label
              htmlFor="fileInput"
              className="flex-1 flex items-center justify-center px-4 py-2 border-2 border-dashed border-gray-300 rounded-lg hover:border-primary-400 hover:bg-primary-50 transition-colors duration-200 cursor-pointer dark:border-gray-700 dark:hover:bg-gray-800"
            >
              <Upload className="h-5 w-5 text-gray-400 mr-2" />
              <span className="text-sm text-gray-600 dark:text-gray-300">
                {fileRef.current?.files?.[0]?.name || t('printers.controls.selectFile')}
              </span>
            </label>
          </div>
        </div>

        {/* Print Button */}
        <motion.button
          onClick={handleStartPrint}
          disabled={!fileRef.current?.files?.[0] || isUploading}
          className={`w-full py-3 px-4 rounded-lg font-medium transition-all duration-200 flex items-center justify-center space-x-2 ${
            !fileRef.current?.files?.[0] || isUploading
              ? 'bg-gray-300 text-gray-500 cursor-not-allowed dark:bg-gray-700 dark:text-gray-400'
              : 'btn-success'
          }`}
          whileHover={fileRef.current?.files?.[0] && !isUploading ? { scale: 1.02 } : {}}
          whileTap={fileRef.current?.files?.[0] && !isUploading ? { scale: 0.98 } : {}}
        >
          {isUploading ? (
            <>
              <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-white"></div>
              <span>{t('printers.controls.upload')}</span>
            </>
          ) : (
            <>
              <Play className="h-4 w-4" />
              <span>{t('printers.controls.startPrint')}</span>
            </>
          )}
        </motion.button>
      </div>

      {/* Printer Information */}
      <div className="space-y-4 mb-6">
        <h4 className="text-sm font-semibold text-gray-700 flex items-center space-x-2 dark:text-gray-200">
          <Settings className="h-4 w-4" />
          <span>{t('printers.details.title')}</span>
        </h4>
        
        <div className="grid grid-cols-2 gap-3 text-sm">
          <div>
            <span className="text-gray-500 dark:text-gray-400">{t('printers.details.type')}:</span>
            <div className="font-medium dark:text-gray-100">{PrinterData.printerType || 'Unknown'}</div>
          </div>
          <div>
            <span className="text-gray-500 dark:text-gray-400">{t('printers.details.version')}:</span>
            <div className="font-medium dark:text-gray-100">{PrinterData.version || 'Unknown'}</div>
          </div>
          <div>
            <span className="text-gray-500 dark:text-gray-400">{t('printers.details.connection')}:</span>
            <div className="font-medium dark:text-gray-100">{PrinterData.typeOfConnection || 'Unknown'}</div>
          </div>
          <div>
            <span className="text-gray-500 dark:text-gray-400">{t('printers.details.status')}:</span>
            <div className={`font-medium ${getStatusColor()}`}>
              {getStatusText()}
            </div>
          </div>
        </div>
      </div>

      {/* Printer Commands */}
      <div className="space-y-4">
        <h4 className="text-sm font-semibold text-gray-700 flex items-center space-x-2 dark:text-gray-200">
          <Move className="h-4 w-4" />
          <span>{t('printers.controls.printControl')}</span>
        </h4>
        
        <PrinterCommands uniqueKey={PrinterData.uniqueKey} />
      </div>
    </motion.div>
  )
}