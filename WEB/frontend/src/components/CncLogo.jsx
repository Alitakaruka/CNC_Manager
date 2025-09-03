import React from 'react'

export default function CncLogo({ className = '' }) {
  return (
    <svg
      className={className}
      viewBox="0 0 24 24"
      fill="none"
      xmlns="http://www.w3.org/2000/svg"
      stroke="currentColor"
      strokeWidth="1.8"
      strokeLinecap="round"
      strokeLinejoin="round"
    >
      {/* Spindle head */}
      <rect x="6" y="3" width="12" height="7" rx="2" />
      {/* Spindle neck */}
      <path d="M12 10v3" />
      {/* Tool holder */}
      <path d="M9 13h6v2H9z" />
      {/* Milling bit */}
      <path d="M12 15l2.5 5h-5L12 15z" />
    </svg>
  )
}
