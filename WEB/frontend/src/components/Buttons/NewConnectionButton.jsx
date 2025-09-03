import { useState } from 'react'
import { Plus, Wifi } from 'lucide-react'
import Modal from '../ModalWindows/ModalForm'
import ConnectionSelector from '../Connection/ConnectionSelector'
import { useLocalization } from '../../hooks/useLocalization.jsx'

export default function NewConnectionButton() {
  const [isOpen, setIsOpen] = useState(false)
  const { t } = useLocalization()

  return (
    <>
      <button
        onClick={() => setIsOpen(true)}
        className="btn-primary flex items-center space-x-2 shadow-glow hover:scale-105 transition-transform duration-200"
      >
        <Plus className="h-4 w-4" />
        <Wifi className="h-4 w-4" />
        <span className="hidden sm:inline">{t('connections.title')}</span>
      </button>

      {isOpen && (
        <Modal onClose={() => setIsOpen(false)} title={t('connections.title')}>
          <ConnectionSelector />
        </Modal>
      )}
    </>
  )
}