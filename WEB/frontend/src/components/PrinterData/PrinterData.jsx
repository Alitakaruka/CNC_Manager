import { useRef, useState } from 'react'
import { Printer, Monitor, FileText } from 'lucide-react'
import PrintersTable from './Tables/PrinterTable'
import Details from './PrinterDetails/PrinterDetails'
import PrinterLogs from './PrinterLogs/PrinterLog'
import { useLocalization } from '../../hooks/useLocalization.jsx'

export default function PrinterData() {
  const [detailsIsOpen, setDetailsIsOpen] = useState(false)
  const [selectedCnc, setSelectedCnc] = useState(null)
  const { t } = useLocalization()

  const handleCncSelect = (cnc) => {
    setSelectedCnc(cnc)
    setDetailsIsOpen(true)
  }

  const handleCloseDetails = () => {
    setDetailsIsOpen(false)
    setSelectedCnc(null)
  }

  return (
    <div className="space-y-4">
      {/* Main Content */}
      <div className="grid grid-cols-1 lg:grid-cols-12 gap-6">
        {/* Reserved Details Panel (no layout shift) */}
        <div className="lg:col-span-3">
          {detailsIsOpen && selectedCnc && (
            <Details 
              PrinterData={selectedCnc} 
              SetDetailsIsOpen={handleCloseDetails}
            />
          )}
        </div>

        {/* CNC Stations Table - constant width */}
        <div className="lg:col-span-6 transition-all duration-300">
          <PrintersTable 
            SetNowPrinter={setSelectedCnc}
            SetDetailsIsOpen={handleCncSelect}
          />
        </div>

        {/* CNC Station Logs - constant width (same as Details for symmetry) */}
        <div className="lg:col-span-3 transition-all duration-300">
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