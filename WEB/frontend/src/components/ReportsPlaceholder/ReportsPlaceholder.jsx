import { BarChart3 } from 'lucide-react'
import { useLocalization } from '../../hooks/useLocalization.jsx'

export default function ReportsPlaceholder() {
  const { t } = useLocalization()
  return (
    <div className="flex flex-col items-center justify-center py-24 text-center">
      <BarChart3 className="h-14 w-14 text-gray-300 dark:text-gray-600 mb-4" />
      <h2 className="text-lg font-semibold text-gray-900 dark:text-gray-100">{t('reportsPage.title')}</h2>
      <p className="mt-2 text-gray-500 dark:text-gray-400 max-w-md">{t('reportsPage.subtitle')}</p>
    </div>
  )
}
