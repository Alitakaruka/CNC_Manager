import { useRef, useState } from 'react'
import { Printer, Monitor, FileText } from 'lucide-react'
import PrintersTable from './Tables/PrinterTable'
import Details from './PrinterDetails/PrinterDetails'
import PrinterLogs from './PrinterLogs/PrinterLog'
import { useLocalization } from '../../hooks/useLocalization.jsx'

export default function PrinterData() {
  const [detailsIsOpen, setDetailsIsOpen] = useState(false)
  const [selectedPrinter, setSelectedPrinter] = useState(null)
  const { t } = useLocalization()

  const handlePrinterSelect = (printer) => {
    setSelectedPrinter(printer)
    setDetailsIsOpen(true)
  }

  const handleCloseDetails = () => {
    setDetailsIsOpen(false)
    setSelectedPrinter(null)
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div className="flex items-center space-x-3">
          <div className="p-3 bg-primary-100 rounded-xl dark:bg-primary-900/40">
            <Printer className="h-8 w-8 text-primary-600 dark:text-primary-400" />
          </div>
          <div>
            <h1 className="text-2xl font-bold text-gray-900 dark:text-gray-100">{t('printers.title')}</h1>
            <p className="text-gray-600 dark:text-gray-400">{t('printers.subtitle')}</p>
          </div>
        </div>
      </div>

      {/* Main Content */}
      <div className="grid grid-cols-1 lg:grid-cols-12 gap-6">
        {/* Printer Details Panel */}
        {detailsIsOpen && selectedPrinter && (
          <div className="lg:col-span-4">
            <Details 
              PrinterData={selectedPrinter} 
              SetDetailsIsOpen={handleCloseDetails}
            />
          </div>
        )}

        {/* Printers Table */}
        <div className={`${detailsIsOpen ? 'lg:col-span-5' : 'lg:col-span-8'} transition-all duration-300`}>
          <PrintersTable 
            SetNowPrinter={setSelectedPrinter}
            SetDetailsIsOpen={handlePrinterSelect}
          />
        </div>

        {/* Printer Logs */}
        <div className={`${detailsIsOpen ? 'lg:col-span-3' : 'lg:col-span-4'} transition-all duration-300`}>
          <PrinterLogs />
        </div>
      </div>

      {/* Empty State */}
      {!detailsIsOpen && (
        <div className="text-center py-12">
          <div className="max-w-md mx-auto">
            <div className="p-4 bg-gray-100 rounded-full w-16 h-16 mx-auto mb-4 flex items-center justify-center dark:bg-gray-800">
              <Monitor className="h-8 w-8 text-gray-400" />
            </div>
            <h3 className="text-lg font-medium text-gray-900 mb-2 dark:text-gray-100">
              {t('printers.empty.selectPrinter')}
            </h3>
            <p className="text-gray-500 dark:text-gray-400">
              {t('printers.empty.selectPrinterSubtitle')}
            </p>
          </div>
        </div>
      )}
    </div>
  )
}  